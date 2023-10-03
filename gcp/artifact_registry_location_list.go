package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"google.golang.org/api/artifactregistry/v1"
)

// func init() {
// 	pluginQueryData = &plugin.QueryData{
// 		ConnectionManager: connection.NewManager(),
// 	}
// }

// BuildregionList :: return a list of matrix items, one per region specified
func BuildArtifactRegistryLocationList(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {

	// have we already created and cached the locations?
	locationCacheKey := "ArtifactRegistry"
	if cachedData, ok := d.ConnectionManager.Cache.Get(locationCacheKey); ok {
		plugin.Logger(ctx).Trace("listlocationDetails:", cachedData.([]map[string]interface{}))
		return cachedData.([]map[string]interface{})
	}

	// Create Service Connection
	service, err := ArtifactRegistryService(ctx, d)
	if err != nil {
		return nil
	}

	// Get project details
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil
	}
	project := projectData.Project

	resp := service.Projects.Locations.List("projects/" + project)
	if err != nil {
		return nil
	}

	var locations []*artifactregistry.Location

	if err := resp.Pages(ctx, func(page *artifactregistry.ListLocationsResponse) error {
		locations = append(locations, page.Locations...)
		return nil
	}); err != nil {
		return nil
	}

	// validate location list
	matrix := make([]map[string]interface{}, len(locations))
	for i, location := range locations {
		matrix[i] = map[string]interface{}{matrixKeyLocation: location.LocationId}
	}
	d.ConnectionManager.Cache.Set(locationCacheKey, matrix)
	return matrix
}
