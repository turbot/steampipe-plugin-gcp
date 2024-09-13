package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)



// BuildDataprocMetastoreLocationList :: return a list of matrix items, one per region specified
func BuildDataprocMetastoreLocationList(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {

	// have we already created and cached the locations?
	locationCacheKey := "DataprocMetastore"
	if cachedData, ok := d.ConnectionManager.Cache.Get(locationCacheKey); ok {
		plugin.Logger(ctx).Trace("listlocationDetails:", cachedData.([]map[string]interface{}))
		return cachedData.([]map[string]interface{})
	}

	// Create Service Connection
	service, err := DataprocMetastoreService(ctx, d)
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
	var matrix []map[string]interface{}
	for _, location := range resp.Locations {
		matrix = append(matrix, map[string]interface{}{matrixKeyLocation: location.LocationId})
	}
	d.ConnectionManager.Cache.Set(locationCacheKey, matrix)
	return matrix
}
