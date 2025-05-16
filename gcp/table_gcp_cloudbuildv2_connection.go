package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/cloudbuild/v2"
)

func tableGcpCloudbuildv2Connection(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_cloudbuildv2_connection",
		Description: "GCP Cloud Build v2 Connection",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "location"}),
			Hydrate:    getGcpCloudbuildv2Connection,
		},
		List: &plugin.ListConfig{
			Hydrate:       listGcpCloudbuildv2Connections,
			ParentHydrate: listGcpCloudbuildv2Locations,
		},
		Columns: []*plugin.Column{
			// Key columns
			{
				Name:        "name",
				Description: "The resource name of the connection.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "Time when the connection was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "update_time",
				Description: "Time when the connection was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "etag",
				Description: "This checksum is computed by the server based on the value of other fields.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "disabled",
				Description: "If disabled is set to true, functionality is disabled for this connection.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "reconciling",
				Description: "Indicates whether the connection is currently being reconciled. A connection is reconciling when there is a change in configuration that needs to be applied.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "installation_state",
				Description: "The installation state of the connection.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "github_config",
				Description: "Configuration for connections to github.com.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "github_enterprise_config",
				Description: "Configuration for connections to an instance of GitHub Enterprise.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "gitlab_config",
				Description: "Configuration for connections to gitlab.com or an instance of GitLab Enterprise.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "bitbucket_data_center_config",
				Description: "Configuration for connections to Bitbucket Data Center.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "bitbucket_cloud_config",
				Description: "Configuration for connections to Bitbucket Cloud.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "annotations",
				Description: "Annotations to allow clients to store small amounts of arbitrary data.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "iam_policy",
				Description: "An Identity and Access Management (IAM) policy for the connection.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGcpCloudbuildv2ConnectionIamPolicy,
				Transform:   transform.FromValue(),
			},

			// Standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(lastPathElement),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     gcpCloudbuildv2ConnectionAka,
				Transform:   transform.FromValue(),
			},

			// Standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpCloudbuildv2ConnectionLocation, "Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Hydrate:     getProject,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listGcpCloudbuildv2Locations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// have we already created and cached the locations?
	locationCacheKey := "CloudBuildV2"
	if cachedData, ok := d.ConnectionManager.Cache.Get(locationCacheKey); ok {
		locations := cachedData.([]string)
		for _, location := range locations {
			d.StreamListItem(ctx, map[string]string{"location": location})
		}
		return nil, nil
	}

	// Create service connection
	service, err := CloudBuildV2Service(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}
	project := projectData.Project

	resp := service.Projects.Locations.List("projects/" + project)
	if err != nil {
		return nil, err
	}

	var locations []string
	if err := resp.Pages(ctx, func(page *cloudbuild.ListLocationsResponse) error {
		for _, location := range page.Locations {
			locations = append(locations, location.LocationId)
			d.StreamListItem(ctx, map[string]string{"location": location.LocationId})
		}
		return nil
	}); err != nil {
		plugin.Logger(ctx).Error("listGcpCloudbuildv2Locations", "error", err)
		return nil, err
	}

	// Cache the location list
	d.ConnectionManager.Cache.Set(locationCacheKey, locations)

	return nil, nil
}

func listGcpCloudbuildv2Connections(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service connection
	service, err := CloudBuildV2Service(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details
	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	// Max limit isn't mentioned in the documentation
	// Default limit is set as 1000
	maxResults := types.Int64(1000)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *maxResults {
			maxResults = limit
		}
	}

	location := h.Item.(map[string]string)["location"]
	parent := "projects/" + project + "/locations/" + location

	resp := service.Projects.Locations.Connections.List(parent).PageSize(*maxResults)
	if err := resp.Pages(ctx, func(page *cloudbuild.ListConnectionsResponse) error {
		for _, connection := range page.Connections {
			d.StreamListItem(ctx, connection)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				page.NextPageToken = ""
				return nil
			}
		}
		return nil
	}); err != nil {
		plugin.Logger(ctx).Debug("listGcpCloudbuildv2Connections", "location", location, "error", err)
		return nil, nil // Skip locations that return errors
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getGcpCloudbuildv2Connection(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["name"].GetStringValue()
	location := d.EqualsQuals["location"].GetStringValue()

	// Create service connection
	service, err := CloudBuildV2Service(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details
	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	// If name is a full resource name, extract the connection name
	if strings.HasPrefix(name, "projects/") {
		parts := strings.Split(name, "/")
		if len(parts) == 6 {
			name = parts[5]
		}
	}

	fullName := "projects/" + project + "/locations/" + location + "/connections/" + name
	req := service.Projects.Locations.Connections.Get(fullName)
	connection, err := req.Do()
	if err != nil {
		plugin.Logger(ctx).Debug("getGcpCloudbuildv2Connection", "name", fullName, "error", err)
		return nil, err
	}

	return connection, nil
}

func getGcpCloudbuildv2ConnectionIamPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	connection := h.Item.(*cloudbuild.Connection)

	// Create service connection
	service, err := CloudBuildV2Service(ctx, d)
	if err != nil {
		return nil, err
	}

	req := service.Projects.Locations.Connections.GetIamPolicy(connection.Name)
	policy, err := req.Do()
	if err != nil {
		return nil, err
	}

	return policy, nil
}

//// TRANSFORM FUNCTIONS

func gcpCloudbuildv2ConnectionAka(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	connection := h.Item.(*cloudbuild.Connection)
	akas := []string{"gcp://cloudbuild.googleapis.com/" + connection.Name}
	return akas, nil
}

func gcpCloudbuildv2ConnectionLocation(_ context.Context, d *transform.TransformData) (interface{}, error) {
	connection := d.HydrateItem.(*cloudbuild.Connection)
	parts := strings.Split(connection.Name, "/")
	if len(parts) < 4 {
		plugin.Logger(context.Background()).Warn("gcpCloudbuildv2ConnectionLocation", "malformed_resource_name", connection.Name)
		return nil, nil
	}
	return parts[3], nil
}
