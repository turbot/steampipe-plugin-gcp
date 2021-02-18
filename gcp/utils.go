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
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
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
	// To call the set

	// have we already created and cached the session?
	serviceCacheKey := "gcp_project_name"

	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*projectInfo), nil
	}

	var err error
	var projectData *projectInfo
	gcpProject := os.Getenv("GCP_PROJECT")
	sdkCoreProject := os.Getenv("CLOUDSDK_CORE_PROJECT")

	if sdkCoreProject != "" {
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
			panic("\n\nYou must specify a Project. You can configure your project by setting \"CLOUDSDK_CORE_PROJECT\" or \"GCP_PROJECT\" environment variable.")
		}
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, projectData)
	plugin.Logger(ctx).Warn("activeProject", "projectData", projectData)

	return projectData, nil
}

func getProjectFromCLI() (*projectInfo, error) {
	const gcloudCLIPath = "/usr/lib/google-cloud-sdk/bin"

	// The default install paths are used to find Google cloud CLI.
	// This is for security, so that any path in the calling program's Path environment is not used to execute Google Cloud CLI.
	// https://stackoverflow.com/questions/62949119/how-to-use-google-cloud-shell-with-the-new-windows-terminal
	gcloudCLIDefaultPathWindows := fmt.Sprintf("%s\\Google\\Cloud SDK\\cloud_env.bat; %s\\Google\\Cloud SDK\\cloud_env.bat", os.Getenv("ProgramFiles(x86)"), os.Getenv("ProgramFiles"))

	// Default path for non-Windows.
	const gcloudCLIDefaultPath = "/bin:/sbin:/usr/bin:/usr/local/bin"

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

// Set project values from config
func setSessionConfig(connection *plugin.Connection) {
	gcpConfig := GetConfig(connection)

	if &gcpConfig != nil {
		if gcpConfig.Project != nil {
			os.Setenv("CLOUDSDK_CORE_PROJECT", *gcpConfig.Project)
		}
		if gcpConfig.CredentialFile != nil {
			os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", *gcpConfig.CredentialFile)
		}
	}
}
