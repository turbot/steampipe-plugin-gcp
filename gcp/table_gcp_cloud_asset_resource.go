package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"

	"google.golang.org/api/cloudasset/v1"
)

//// TABLE DEFINITION

func tableGcpCloudAssetResource(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_cloud_asset_resource",
		Description: "GCP Cloud Asset Resource",
		List: &plugin.ListConfig{
			Hydrate: listCloudAssetResources,
			Tags:    map[string]string{"service": "cloudasset", "action": "assets.listResource"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The full name of the asset.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "asset_type",
				Description: "The type of the asset.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "update_time",
				Description: "The last update timestamp of an asset.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "version",
				Description: "The API version of the resource, e.g., v1.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Resource.Version"),
			},
			{
				Name:        "discovery_document_uri",
				Description: "The URL of the discovery document containing the resource's JSON schema.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Resource.DiscoveryDocumentUri"),
			},
			{
				Name:        "discovery_name",
				Description: "The JSON schema name listed in the discovery document, e.g., Project.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Resource.DiscoveryName"),
			},
			{
				Name:        "resource_url",
				Description: "The REST URL for accessing the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Resource.ResourceUrl"),
			},
			{
				Name:        "parent",
				Description: "The full name of the immediate parent of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Resource.Parent"),
			},
			{
				Name:        "location",
				Description: "The location of the resource in Google Cloud, such as its zone and region.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Resource.Location"),
			},
			{
				Name:        "data",
				Description: "The content of the resource, in which some sensitive fields are removed and may not be present.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Resource.Data"),
			},
			{
				Name:        "ancestors",
				Description: "The ancestry path of an asset in Google Cloud resource hierarchy.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},

			// Standard GCP columns
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

func listCloudAssetResources(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return listCloudAssetsByContentType(ctx, d, h, "RESOURCE", "gcp_cloud_asset_resource.listCloudAssetResources", func(ctx context.Context, d *plugin.QueryData, asset *cloudasset.Asset) {
		d.StreamListItem(ctx, asset)
	})
}
