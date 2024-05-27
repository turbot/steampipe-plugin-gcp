package gcp

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"unicode/utf8"

	"github.com/mitchellh/go-homedir"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/memoize"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/impersonate"
	"google.golang.org/api/option"
)

func getLastPathElement(path string) string {
	if path == "" {
		return ""
	}

	pathItems := strings.Split(path, "/")
	return pathItems[len(pathItems)-1]
}

type projectInfo struct {
	Project string `json:"project,omitempty"`
}

// Constants for Standard Column Descriptions
const (
	ColumnDescriptionAkas     = "Array of globally unique identifier strings (also known as) for the resource."
	ColumnDescriptionTags     = "A map of tags for the resource."
	ColumnDescriptionTitle    = "Title of the resource."
	ColumnDescriptionProject  = "The GCP Project in which the resource is located."
	ColumnDescriptionLocation = "The GCP multi-region, region, or zone in which the resource is located."
)

//// TRANSFORM FUNCTIONS

func lastPathElement(_ context.Context, d *transform.TransformData) (interface{}, error) {
	return getLastPathElement(types.SafeString(d.Value)), nil
}

// if the caching is required other than per connection, build a cache key for the call and use it in Memoize
// since getProject is a call, caching should be per connection
var getProjectMemoized = plugin.HydrateFunc(getProjectUncached).Memoize(memoize.WithCacheKeyFunction(getProjectCacheKey))

// Build a cache key for the call to getProject.
func getProjectCacheKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	key := fmt.Sprintf("getGCPProjectInfo%s", "")
	return key, nil
}

func getProject(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (any, error) {
	projectId, err := getProjectMemoized(ctx, d, h)
	if err != nil {
		return nil, err
	}

	return projectId, nil
}

func getProjectUncached(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cacheKey := "getGCPProjectInfo"
	var err error
	var projectData *projectInfo
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		projectData = cachedData.(*projectInfo)
	} else {
		projectData, err = activeProject(ctx, d)
		if err != nil {
			return nil, err
		}
		// save to extension cache
		d.ConnectionManager.Cache.Set(cacheKey, projectData)
	}
	return projectData.Project, nil
}

func activeProject(ctx context.Context, d *plugin.QueryData) (*projectInfo, error) {
	// have we already created and cached the session?
	serviceCacheKey := "gcp_project_id"

	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*projectInfo), nil
	}

	var err error
	var projectData *projectInfo
	gcpProject := os.Getenv("GCP_PROJECT")
	sdkCoreProject := os.Getenv("CLOUDSDK_CORE_PROJECT")
	projectFromConfig := getProjectFromConfig(d.Connection)

	plugin.Logger(ctx).Debug("activeProject", "gcp_project_env_var", gcpProject)
	plugin.Logger(ctx).Debug("activeProject", "cloudsdk_core_project_env_var", sdkCoreProject)
	plugin.Logger(ctx).Debug("activeProject", "config_project", projectFromConfig)

	if projectFromConfig != "" {
		projectData = &projectInfo{
			Project: projectFromConfig,
		}
	} else if sdkCoreProject != "" {
		projectData = &projectInfo{
			Project: sdkCoreProject,
		}
	} else if gcpProject != "" {
		projectData = &projectInfo{
			Project: gcpProject,
		}
	} else {
		projectData, err = getProjectFromCLI(ctx)
		if err != nil {
			return nil, err
		}
	}

	plugin.Logger(ctx).Debug("activeProject", "project_data", projectData)

	// No active project is set
	if projectData == nil {
		return nil, fmt.Errorf("an active project must be set")
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, projectData)

	return projectData, nil
}

func getProjectFromCLI(ctx context.Context) (*projectInfo, error) {
	// The default install paths are used to find Google cloud CLI.
	// This is for security, so that any path in the calling program's Path environment is not used to execute Google Cloud CLI.
	// https://stackoverflow.com/questions/62949119/how-to-use-google-cloud-shell-with-the-new-windows-terminal
	gcloudCLIDefaultPathWindows := fmt.Sprintf("%s\\Google\\Cloud SDK\\cloud_env.bat; %s\\Google\\Cloud SDK\\cloud_env.bat", os.Getenv("ProgramFiles(x86)"), os.Getenv("ProgramFiles"))

	// Execute GOOGLE CLOUD CLI to get project info
	var cliCmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cliCmd = exec.Command(fmt.Sprintf("%s\\system32\\cmd.exe", os.Getenv("windir")))
		cliCmd.Env = os.Environ()
		cliCmd.Env = append(cliCmd.Env, fmt.Sprintf("PATH=%s", gcloudCLIDefaultPathWindows))
		cliCmd.Args = append(cliCmd.Args, "/c", "gcloud")
	} else {
		cliCmd = exec.Command("gcloud")
		cliCmd.Env = os.Environ()
	}
	cliCmd.Args = append(cliCmd.Args, "config", "get-value", "project", "--format", "object")

	var stderr bytes.Buffer
	cliCmd.Stderr = &stderr

	output, err := cliCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("invoking gcloud CLI failed with the following error: %v", stderr)
	}

	plugin.Logger(ctx).Debug("getProjectFromCLI", "cmd_output", output)

	// Output will be '[]' if no project is set
	if len(output) < 1 {
		return nil, nil
	}

	project := types.ToString(output)
	plugin.Logger(ctx).Debug("getProjectFromCLI", "project", project)

	return &projectInfo{
		Project: project,
	}, nil
}

func getProjectFromConfig(connection *plugin.Connection) string {
	gcpConfig := GetConfig(connection)

	if gcpConfig.Project != nil {
		return *gcpConfig.Project
	}
	return ""
}

func base64DecodedData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data, err := base64.StdEncoding.DecodeString(types.SafeString(d.Value))
	// check if CorruptInputError or invalid UTF-8
	if err != nil {
		return nil, nil
	} else if !utf8.Valid(data) {
		return types.SafeString(d.Value), nil
	}
	return data, nil
}

// Set project values from config and return client options
func setSessionConfig(ctx context.Context, connection *plugin.Connection) []option.ClientOption {
	gcpConfig := GetConfig(connection)
	opts := []option.ClientOption{}

	// 'credential_file' in connection config is DEPRECATED, and will be removed in future release
	// use `credentials` instead
	if gcpConfig.Credentials != nil {
		contents, err := pathOrContents(*gcpConfig.Credentials)
		if err != nil {
			panic(err)
		}
		opts = append(opts, option.WithCredentialsJSON([]byte(contents)))
	} else if gcpConfig.CredentialFile != nil {
		opts = append(opts, option.WithCredentialsFile(*gcpConfig.CredentialFile))
	}

	if gcpConfig.ImpersonateServiceAccount != nil {
		ts, err := impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
			TargetPrincipal: *gcpConfig.ImpersonateServiceAccount,
			Scopes:          []string{"https://www.googleapis.com/auth/cloud-platform"},
		})
		if err != nil {
			panic(err)
		}

		opts = append(opts, option.WithTokenSource(ts))
	}

	// check if quota project is set via env var
	quotaProject := os.Getenv("GOOGLE_CLOUD_QUOTA_PROJECT")

	// check if quota project is set in config
	if gcpConfig.QuotaProject != nil {
		quotaProject = *gcpConfig.QuotaProject
	}

	if quotaProject != "" {
		opts = append(opts, option.WithQuotaProject(quotaProject))
	}

	return opts
}

// Returns the content of given file, or the inline JSON credential as it is
func pathOrContents(poc string) (string, error) {
	if len(poc) == 0 {
		return poc, nil
	}

	path := poc
	if path[0] == '~' {
		var err error
		path, err = homedir.Expand(path)
		if err != nil {
			return path, err
		}
	}

	// Check for valid file path
	if _, err := os.Stat(path); err == nil {
		contents, err := os.ReadFile(path)
		if err != nil {
			return string(contents), err
		}
		return string(contents), nil
	}

	// Return error if content is a file path and the file doesn't exist
	if len(path) > 1 && (path[0] == '/' || path[0] == '\\') {
		return "", fmt.Errorf("%s: no such file or dir", path)
	}

	// Return the inline content
	return poc, nil
}

// Get QualValueList as an list of items
func getListValues(listValue *proto.QualValueList) []string {
	values := make([]string, 0)
	for _, value := range listValue.Values {
		values = append(values, value.GetStringValue())
	}
	return values
}

/**
 * buildQueryFilter: To build gcp query filter from equal quals
 * Sample for gcp_compute_image table
 * select name, id, status, source_project, deprecation_state, family
 * from	gcp_morales_aaa.gcp_compute_image
 * where family in ('sles-12', 'sles-15') and deprecation_state = 'ACTIVE'
 * -------------------------------------------------------------------------
 * 	Column: family, Operator: "=", Value: "[sles-12 sles-15]"
 * 	Column: deprecation_state, Operator: "=", Value: "ACTIVE"
 * -------------------------------------------------------------------------
 *
 * Output: []string{"(family = "sles-12") OR (family = "sles-15")", "(deprecated.state = "ACTIVE")"}
 */
func buildQueryFilter(filterQuals []filterQualMap, equalQuals plugin.KeyColumnEqualsQualMap) []string {
	filters := []string{}

	for _, qual := range filterQuals {
		qualValue := equalQuals[qual.ColumnName]
		if qualValue != nil {
			switch qual.Type {
			case "string":

				// In case of a in caluse
				if qualValue.GetListValue() != nil {
					filter := ""
					for i, q := range qualValue.GetListValue().Values {
						if i == 0 {
							filter = fmt.Sprintf("(%s = \"%s\")", qual.PropertyPath, q.GetStringValue())
						} else {
							filter = fmt.Sprintf("%s OR (%s = \"%s\")", filter, qual.PropertyPath, q.GetStringValue())
						}
					}
					filters = append(filters, fmt.Sprintf("(%s)", filter))
				} else {
					filters = append(filters, fmt.Sprintf("(%s = \"%s\")", qual.PropertyPath, qualValue.GetStringValue()))
				}
			case "boolean":
				filters = append(filters, fmt.Sprintf("(%s = %t)", qual.PropertyPath, qualValue.GetBoolValue()))
			}
		}
	}

	return filters
}

/**
 * buildQueryFilter: To build gcp query filter from equal quals
 * Sample for gcp_compute_instance table
 * select name, id, machine_type_name, status, can_ip_forward, cpu_platform, deletion_protection, start_restricted, hostname
 * from gcp_morales_aaa.gcp_compute_instance
 * where
 *	status in ('TERMINATED', 'RUNNING') and
 *	cpu_platform = 'Intel Haswell' and
 *  not deletion_protection
 *  -----------------------STEAMPIPE QUAL INFO-----------------------------------------
 *  	Column: deletion_protection, Operator: '<>', Value: 'true'
 *  	Column: status, Operator: '=', Value: '[TERMINATED RUNNING]'
 *  	Column: cpu_platform, Operator: '=', Value: 'Intel Haswell'
 *  ----------------------------------------------------------------
 *
 * Output: []string{"(cpuPlatform = \"Intel Haswell\")", "((status = \"TERMINATED\") OR (status = \"RUNNING\"))", "(deletionProtection = false)"}
 *
 * This can be used for almost all the API's in GCP if it supports filter option
 */
func buildQueryFilterFromQuals(filterQuals []filterQualMap, equalQuals plugin.KeyColumnQualMap) []string {
	filters := []string{}

	for _, filterQualItem := range filterQuals {
		filterQual := equalQuals[filterQualItem.ColumnName]
		if filterQual == nil {
			continue
		}

		// Check only if filter qual map matches with optional column name
		if filterQual.Name == filterQualItem.ColumnName {
			if filterQual.Quals == nil {
				continue
			}

			for _, qual := range filterQual.Quals {
				if qual.Value != nil {
					value := qual.Value
					switch filterQualItem.Type {
					case "string":
						// In case of IN caluse
						if value.GetListValue() != nil {
							filter := ""
							for i, q := range value.GetListValue().Values {
								if i == 0 {
									filter = fmt.Sprintf("(%s = \"%s\")", filterQualItem.PropertyPath, q.GetStringValue())
								} else {
									filter = fmt.Sprintf("%s OR (%s = \"%s\")", filter, filterQualItem.PropertyPath, q.GetStringValue())
								}
							}
							filters = append(filters, fmt.Sprintf("(%s)", filter))
						} else {
							switch qual.Operator {
							case "=", "<>", "!=", ">", ",":
								filters = append(filters, fmt.Sprintf("(%s %s \"%s\")", filterQualItem.PropertyPath, GcpFilterOperatorMap[qual.Operator], value.GetStringValue()))
							case "<=", ">=":
								filters = append(filters, fmt.Sprintf("((%s = \"%s\") OR (%s %s \"%s\"))", filterQualItem.PropertyPath, value.GetStringValue(), filterQualItem.PropertyPath, GcpFilterOperatorMap[qual.Operator], value.GetStringValue()))
							}
						}
					case "boolean":
						boolValue := value.GetBoolValue()
						switch qual.Operator {
						case "<>":
							filters = append(filters, fmt.Sprintf("(%s = %t)", filterQualItem.PropertyPath, !boolValue))
						case "=":
							filters = append(filters, fmt.Sprintf("(%s = %t)", filterQualItem.PropertyPath, boolValue))
						}
					}
				}
			}

		}
	}

	return filters
}

type filterQualMap struct {
	ColumnName   string
	PropertyPath string
	Type         string
}

// Steampipe to GCP query filter map
//
// Filter sets the optional parameter "filter": A filter expression that
// filters resources listed in the response. The expression must specify
// the field name, a comparison operator, and the value that you want to
// use for filtering. The value must be a string, a number, or a
// boolean. The comparison operator must be either `=`, `!=`, `>`, or
// `<`.
var GcpFilterOperatorMap = map[string]string{
	"=":  "=",
	"<>": "!=",
	"!=": "!=",
	">":  ">",
	"<":  "<",
	"<=": "<", // Filter ((property=value) OR (property<value))
	">=": ">", // Filter ((property=value) OR (property>value))
}
