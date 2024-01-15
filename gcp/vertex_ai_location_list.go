package gcp

import (
	"context"
	"strings"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"google.golang.org/api/iterator"
	"google.golang.org/genproto/googleapis/cloud/location"
)

func BuildVertexAILocationListByClientType(clientType string) func(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {
	return BuildVertexAILocationList(clientType)
}

func BuildVertexAILocationList(clientType string) func(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {
	return func(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {
		// have we already created and cached the locations?
		locationCacheKey := "BuildVertexAILocationList" + clientType
		if cachedData, ok := d.ConnectionManager.Cache.Get(locationCacheKey); ok {
			plugin.Logger(ctx).Trace("listlocationDetails:", cachedData.([]map[string]interface{}))
			return cachedData.([]map[string]interface{})
		}

		// Create Service Connection
		service, err := AIService(ctx, d, clientType)
		if err != nil {
			return nil
		}

		// Get project details
		projectData, err := activeProject(ctx, d)
		if err != nil {
			return nil
		}
		project := projectData.Project

		var resourceLocations []*location.Location
		input := &location.ListLocationsRequest{
			Name: "projects/" + project,
		}

		switch clientType {
		case "Endpoint":
			resp := service.Endpoint.ListLocations(ctx, input)
			resourceLocations = append(resourceLocations, iterateLocationResponse(resp)...)
		case "Dataset":
			resp := service.Dataset.ListLocations(ctx, input)
			resourceLocations = append(resourceLocations, iterateLocationResponse(resp)...)
		case "Index":
			resp := service.Index.ListLocations(ctx, input)
			resourceLocations = append(resourceLocations, iterateLocationResponse(resp)...)
		case "Job":
			resp := service.Job.ListLocations(ctx, input)
			resourceLocations = append(resourceLocations, iterateLocationResponse(resp)...)
		}

		// validate location list
		matrix := make([]map[string]interface{}, len(resourceLocations))
		for i, location := range resourceLocations {
			matrix[i] = map[string]interface{}{matrixKeyLocation: location.LocationId}
		}
		d.ConnectionManager.Cache.Set(locationCacheKey, matrix)
		return matrix

	}
}

func iterateLocationResponse(response *aiplatform.LocationIterator) []*location.Location {
	var loc []*location.Location
	for {

		res, err := response.Next()
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				return nil
			}
			if err == iterator.Done {
				break
			}
		}
		loc = append(loc, res)
	}
	return loc
}
