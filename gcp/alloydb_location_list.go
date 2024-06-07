package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// BuildregionList :: return a list of matrix items, one per region specified
func BuildAlloyDBLocationList(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {

	// have we already created and cached the locations?
	locationCacheKey := "AlloydbLocation"
	if cachedData, ok := d.ConnectionManager.Cache.Get(locationCacheKey); ok {
		plugin.Logger(ctx).Trace("listlocationDetails:", cachedData.([]map[string]interface{}))
		return cachedData.([]map[string]interface{})
	}

	// Create Service Connection
	service, err := AlloyDBService(ctx, d)
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
