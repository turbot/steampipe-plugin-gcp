package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"google.golang.org/api/cloudasset/v1"
)

//// TABLE DEFINITION

func tableGcpAllResource(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_all_resources",
		Description: "List All GCP Resources",
		List: &plugin.ListConfig{
			Hydrate: listAllResourcees,
		},
		Columns: []*plugin.Column{
			{
				Name:        "display_name",
				Description: "A friendly name that identifies the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of this resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AssetType"),
			},
			{
				Name:        "name",
				Description: "The full resource name of this resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A user-specified, human-readable description of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "labels",
				Description: "A list of labels associated with this resource.",
				Type:        proto.ColumnType_JSON,
			},

			// standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(lastPathElement),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromConstant(activeProject()),
			},
		},
	}
}

//// LIST FUNCTION

func listAllResourcees(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	service, err := cloudasset.NewService(ctx)
	if err != nil {
		return nil, err
	}

	project := activeProject()
	resp := service.V1.SearchAllResources("projects/" + project)
	if err := resp.Pages(ctx, func(page *cloudasset.SearchAllResourcesResponse) error {
		for _, resource := range page.Results {
			d.StreamListItem(ctx, resource)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, err
}
