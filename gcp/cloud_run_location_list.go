package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// Based on the available documentation, there is no indication that the Cloud Run Admin API v1 supports specific locations that v2 does not.
// Both API versions are designed to operate across all regions where Cloud Run is available.
// The primary differences between v1 and v2 pertain to API design and compatibility, rather than regional support.
// Using the Cloud Run V1 API to list supported regions will not have any negative impact, even though the Cloud Run V2 API is used in the gcp_cloud_run_* tables for building the region matrix. The region data remains consistent across both API versions, ensuring accurate coverage without affecting functionality.

// https://cloud.google.com/run/docs/locations
// BuildCloudRunLocationList :: return a list of matrix items, one per region specified
func BuildCloudRunLocationList(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {

	// have we already created and cached the locations?
	locationCacheKey := "CloudRunLocation"
	if cachedData, ok := d.ConnectionManager.Cache.Get(locationCacheKey); ok {
		plugin.Logger(ctx).Trace("listlocationDetails:", cachedData.([]map[string]interface{}))
		return cachedData.([]map[string]interface{})
	}

	// Create Service Connection
	service, err := CloudRunServiceV1(ctx, d)
	if err != nil {
		return nil
	}

	// Get project details
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil
	}
	project := projectData.Project


	resp, err := service.Projects.Locations.List("projects/" + project).Do()
	if err != nil {
		return nil
	}
	// validate location list
	matrix := make([]map[string]interface{}, len(resp.Locations))
	for i, location := range resp.Locations {
		matrix[i] = map[string]interface{}{matrixKeyLocation: location.LocationId}
	}
	d.ConnectionManager.Cache.Set(locationCacheKey, matrix)
	return matrix
}
