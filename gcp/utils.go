package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

var projectName string

func init() {
	projectName = ""
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

func getProject(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	projectData, err := activeProject(ctx, d.ConnectionManager)
	if err != nil {
		return nil, err
	}

	return projectData.Project, nil
}
