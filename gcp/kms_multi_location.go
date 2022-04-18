package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v3/connection"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

var pluginQueryData *plugin.QueryData

func init() {
	pluginQueryData = &plugin.QueryData{
		ConnectionManager: connection.NewManager(),
	}
}

const matrixKeyLocation = "location"

// BuildregionList :: return a list of matrix items, one per region specified
func BuildLocationList(ctx context.Context, connection *plugin.Connection) []map[string]interface{} {

	pluginQueryData.Connection = connection

	// have we already created and cached the locations?
	locationCacheKey := "KMSLocation"
	if cachedData, ok := pluginQueryData.ConnectionManager.Cache.Get(locationCacheKey); ok {
		plugin.Logger(ctx).Trace("listlocationDetails:", cachedData.([]map[string]interface{}))
		return cachedData.([]map[string]interface{})
	}

	// Create Service Connection
	service, err := KMSService(ctx, pluginQueryData)
	if err != nil {
		return nil
	}

	// Get project details
	projectData, err := activeProject(ctx, pluginQueryData)
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
	pluginQueryData.ConnectionManager.Cache.Set(locationCacheKey, matrix)
	return matrix
}
