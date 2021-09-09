package gcp

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
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

func getProject(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cacheKey := "getGCPProjectInfo"
	var err error
	var projectData *projectInfo
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		projectData = cachedData.(*projectInfo)
	} else {
		// To set the config argument for the connection in a project
		setSessionConfig(d.Connection)
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
	serviceCacheKey := "gcp_project_name"

	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*projectInfo), nil
	}

	var err error
	var projectData *projectInfo
	gcpProject := os.Getenv("GCP_PROJECT")
	sdkCoreProject := os.Getenv("CLOUDSDK_CORE_PROJECT")
	projectFromConfig := getProjectFromConfig(d.Connection)

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
		projectData, err = getProjectFromCLI()
		if err != nil {
			panic("\n\n'project' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
		}
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, projectData)
	plugin.Logger(ctx).Info("activeProject", "Project", projectData.Project)

	return projectData, nil
}

func getProjectFromCLI() (*projectInfo, error) {
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
		return nil, fmt.Errorf("Invoking gcloud CLI failed with the following error: %v", stderr)
	}

	project := types.ToString(output)

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

// Set project values from config and return client options
func setSessionConfig(connection *plugin.Connection) []option.ClientOption {
	gcpConfig := GetConfig(connection)
	opts := []option.ClientOption{}

	if gcpConfig.CredentialFile != nil {
		opts = append(opts, option.WithCredentialsFile(*gcpConfig.CredentialFile))
	}
	if gcpConfig.ImpersonateServiceAccount != nil {
		opts = append(opts, option.ImpersonateCredentials(*gcpConfig.ImpersonateServiceAccount))
	}
	return opts
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

func GetBoolQualValue(quals plugin.KeyColumnQualMap, columnName string) (value *bool, exists bool) {
	exists = false
	if quals[columnName] == nil {
		return nil, exists
	}

	if quals[columnName].Quals == nil {
		return nil, exists
	}

	for _, qual := range quals[columnName].Quals {
		if qual.Value != nil {
			value := qual.Value
			boolValue := value.GetBoolValue()
			switch qual.Operator {
			case "<>":
				return types.Bool(!boolValue), true
			case "=":
				return types.Bool(boolValue), true
			}
			break
		}
	}
	return nil, exists
}
