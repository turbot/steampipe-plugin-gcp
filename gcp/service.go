package gcp

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/connection"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

type projectInfo struct {
	Project string `json:"project,omitempty"`
}

func activeProject(ctx context.Context, connectionManager *connection.Manager) (*projectInfo, error) {

	// have we already created and cached the session?
	serviceCacheKey := "gcp_project_name"

	if cachedData, ok := connectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*projectInfo), nil
	}

	var projectData *projectInfo
	gcpProject := os.Getenv("GCP_PROJECT")
	sdkCoreProject := os.Getenv("CLOUDSDK_CORE_PROJECT")

	if gcpProject != "" {
		projectData = &projectInfo{
			Project: gcpProject,
		}
	} else if sdkCoreProject != "" {
		projectData = &projectInfo{
			Project: sdkCoreProject,
		}
	} else {
		projectData, err := getProjectFromCLI()
		if err != nil {
			panic("\n\nYou must specify a Project. You can configure your project by setting \"GCP_REGION\" or \"CLOUDSDK_CORE_PROJECT\" environment variable.")
		}
		return projectData, nil
	}

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
