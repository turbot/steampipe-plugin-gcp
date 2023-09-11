package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"google.golang.org/api/iterator"
	"google.golang.org/genproto/googleapis/cloud/location"
)

const matrixKeyAIPlatformLocation = "aiplatform-location"

// BuildRedisLocationList :: return a list of matrix items, one per region specified
func BuildAIPlatformLocationList(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {

	// have we already created and cached the locations?
	locationCacheKey := "AIPlatformLocation"
	if cachedData, ok := d.ConnectionManager.Cache.Get(locationCacheKey); ok {
		plugin.Logger(ctx).Trace("listlocationDetails:", cachedData.([]map[string]interface{}))
		return cachedData.([]map[string]interface{})
	}

	// Create Service Connection
	service, err := AIplatformService(ctx, d)
	if err != nil {
		return nil
	}

	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil
	}
	project := projectData.Project
	if project == ""  {
		return nil
	}

	req := &location.ListLocationsRequest{
		Name:     "projects/" + project,
		PageSize: 100,
	}

	it := service.ListLocations(ctx, req)
	matrix := []map[string]interface{}{}
	for {
		resp, err := it.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return nil
		}
		obj := map[string]interface{}{matrixKeyAIPlatformLocation: resp.LocationId}
		matrix = append(matrix, obj)
	}

	d.ConnectionManager.Cache.Set(locationCacheKey, matrix)
	return matrix
}
