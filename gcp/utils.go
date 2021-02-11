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
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

type projectInfo struct {
	Project string `json:"project,omitempty"`
}

var projectName string

func init() {
	projectName = activeProject()
}

func activeProject() string {
	gcpProject := os.Getenv("GCP_PROJECT")
	sdkCoreProject := os.Getenv("CLOUDSDK_CORE_PROJECT")

	if gcpProject != "" {
		return gcpProject
	} else if sdkCoreProject != "" {
		return sdkCoreProject
	} else {
		projectDetails, err := getProjectFromCLI()
		if err != nil {
			panic(err)
		}
		gcpProject = projectDetails.Project
	}

	return gcpProject
}

func getLastPathElement(path string) string {
	if path == "" {
		return ""
	}

	pathItems := strings.Split(path, "/")
	return pathItems[len(pathItems)-1]
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

func getProjectFromCLI() (*projectInfo, error) {
	// This is the path that a developer can set to tell this class what the install path for Azure CLI is.
	const gcloudCLIPath = "AzureCLIPath"

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
		cliCmd.Env = append(cliCmd.Env, fmt.Sprintf("PATH=%s;%s", os.Getenv(gcloudCLIPath), gcloudCLIDefaultPathWindows))
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
		panic(err)
		// return nil, fmt.Errorf("Invoking gcloud CLI failed with the following error: %v", stderr)
	}

	project := types.ToString(output)

	return &projectInfo{
		Project: project,
	}, nil
}
