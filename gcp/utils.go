package gcp

import (
	"context"
	"os"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func activeProject() string {
	return os.Getenv("GCP_PROJECT")
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
	ColumnDescriptionAkas  = "Array of globally unique identifier strings (also known as) for the resource."
	ColumnDescriptionTags  = "A map of tags for the resource."
	ColumnDescriptionTitle = "Title of the resource."
)

// Get Standard Column Descriptions
func resourceInterfaceDescription(key string) string {
	switch key {
	case "akas":
		return "Array of globally unique identifier strings (also known as) for the resource."
	case "tags":
		return "A map of tags for the resource."
	case "title":
		return "Title of the resource."

	case "project":
		return "The GCP Project in which the resource is located."
	case "location":
		return "The GCP multi-region, region, or zone in which the resource is located."
	case "location_type":
		return "The type of location (multi-region, region, or zone)."
	}
	return ""
}

//// TRANSFORM FUNCTIONS

func lastPathElement(_ context.Context, d *transform.TransformData) (interface{}, error) {
	return getLastPathElement(types.SafeString(d.Value)), nil
}
