package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// BuildTpuZoneList :: return a list of matrix items, one per supported TPU zone
func BuildTpuZoneList(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {
	// have we already created and cached the locations?
	locationCacheKey := "TPU"
	if cachedData, ok := d.ConnectionManager.Cache.Get(locationCacheKey); ok {
		plugin.Logger(ctx).Trace("listTpuLocationDetails:", cachedData.([]map[string]interface{}))
		return cachedData.([]map[string]interface{})
	}

	// Create Service Connection
	service, err := TPUService(ctx, d)
	if err != nil {
		return nil
	}

	// Get project details
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil
	}
	project := projectData.Project

	// List all locations that support TPUs
	parent := "projects/" + project
	resp, err := service.Projects.Locations.List(parent).Do()
	if err != nil {
		plugin.Logger(ctx).Error("BuildTpuZoneList", "Error listing TPU locations", err)
		return nil
	}

	// Create matrix items for each location
	matrix := make([]map[string]interface{}, len(resp.Locations))
	for i, location := range resp.Locations {
		matrix[i] = map[string]interface{}{
			matrixKeyZone: location.LocationId,
		}
	}

	// Cache the results
	d.ConnectionManager.Cache.Set(locationCacheKey, matrix)
	return matrix
}
