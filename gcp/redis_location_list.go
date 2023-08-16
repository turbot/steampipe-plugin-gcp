package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"google.golang.org/api/iterator"
	"google.golang.org/genproto/googleapis/cloud/location"
)

const matrixKeyRedisLocation = "redis-location"

func BuildRedisLocationList(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {

	// have we already created and cached the locations?
	locationCacheKey := "RedisLocation"
	if cachedData, ok := d.ConnectionManager.Cache.Get(locationCacheKey); ok {
		plugin.Logger(ctx).Trace("listlocationDetails:", cachedData.([]map[string]interface{}))
		return cachedData.([]map[string]interface{})
	}

	// Create Service Connection
	service, err := RedisService(ctx, d)
	if err != nil {
		return nil
	}

	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil
	}
	project := projectData.Project

	req := &location.ListLocationsRequest{
		Name:     "projects/" + project,
		PageSize: 100,
	}

	it := service.ListLocations(ctx, req)
	// var matrix []map[string]interface{}
	matrix := []map[string]interface{}{}
	for {
		resp, err := it.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return nil
		}
		obj := map[string]interface{}{matrixKeyRedisLocation: resp.LocationId}
		matrix = append(matrix, obj)
	}

	d.ConnectionManager.Cache.Set(locationCacheKey, matrix)
	return matrix
}
