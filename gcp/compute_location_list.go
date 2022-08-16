package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v3/connection"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

func init() {
	pluginQueryData = &plugin.QueryData{
		ConnectionManager: connection.NewManager(),
	}
}


// BuildregionList :: return a list of matrix items, one per region specified
// https://cloud.google.com/dataproc/docs/concepts/regional-endpoints
func BuildComputeLocationList(ctx context.Context, connection *plugin.Connection) []map[string]interface{} {

	pluginQueryData.Connection = connection

	// have we already created and cached the locations?
	locationCacheKey := "Compute"
	if cachedData, ok := pluginQueryData.ConnectionManager.Cache.Get(locationCacheKey); ok {
		plugin.Logger(ctx).Trace("listlocationDetails:", cachedData.([]map[string]interface{}))
		return cachedData.([]map[string]interface{})
	}

	// Create Service Connection
	service, err := ComputeService(ctx, pluginQueryData)
	if err != nil {
		return nil
	}

	// Get project details
	projectData, err := activeProject(ctx, pluginQueryData)
	if err != nil {
		return nil
	}
	project := projectData.Project

	resp, err := service.Regions.List(project).Do()
	if err != nil {
		return nil
	}

	// validate location list
	matrix := make([]map[string]interface{}, len(resp.Items))
	for i, location := range resp.Items {
		matrix[i] = map[string]interface{}{matrixKeyLocation: location}
	}
	pluginQueryData.ConnectionManager.Cache.Set(locationCacheKey, matrix)
	return matrix
}
