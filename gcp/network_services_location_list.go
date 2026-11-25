package gcp

import (
	"context"
	"regexp"

	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"google.golang.org/api/networkservices/v1"
)

// BuildNetworkServicesLocationList :: return a list of matrix items, one per Network Services location
func BuildNetworkServicesLocationList(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {

	// have we already created and cached the locations?
	locationCacheKey := "BuildNetworkServicesLocationList"
	if cachedData, ok := d.ConnectionCache.Get(ctx, locationCacheKey); ok {
		plugin.Logger(ctx).Debug("BuildNetworkServicesLocationList:", cachedData.([]map[string]interface{}))
		return cachedData.([]map[string]interface{})
	}

	// Create Service Connection
	service, err := NetworkServicesService(ctx, d)
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

	var locations []*networkservices.Location

	if err := resp.Pages(ctx, func(page *networkservices.ListLocationsResponse) error {
		locations = append(locations, page.Locations...)
		return nil
	}); err != nil {
		return nil
	}

	// The locations list also returns zones (for example us-central1-a) and
	// other location kinds, but serviceLbPolicies only exist in regions and
	// global; listing any other location fails with 400 "Malformed name".
	matrix := make([]map[string]interface{}, 0, len(locations))
	for _, location := range locations {
		if !isPolicyLocationID(location.LocationId) {
			continue
		}
		matrix = append(matrix, map[string]interface{}{matrixKeyLocation: location.LocationId})
	}
	if err := d.ConnectionCache.Set(ctx, locationCacheKey, matrix); err != nil {
		plugin.Logger(ctx).Warn("gcp.cache", "cache_key", locationCacheKey, "error", err)
	}
	return matrix
}

var regionLocationIDPattern = regexp.MustCompile(`^[a-z]+-[a-z]+[0-9]+$`)

// isPolicyLocationID reports whether a location ID is global or a region
// (us-central1), the only locations that hold serviceLbPolicies.
func isPolicyLocationID(id string) bool {
	return id == "global" || regionLocationIDPattern.MatchString(id)
}
