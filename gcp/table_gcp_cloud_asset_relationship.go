package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"

	"google.golang.org/api/cloudasset/v1"
)

//// TABLE DEFINITION

func tableGcpCloudAssetRelationship(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_cloud_asset_relationship",
		Description: "GCP Cloud Asset Relationship",
		List: &plugin.ListConfig{
			Hydrate: listCloudAssetRelationships,
			Tags:    map[string]string{"service": "cloudasset", "action": "assets.listRelationship"},
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
				Name:        "relationship_type",
				Description: "The unique identifier of the relationship type, e.g., INSTANCE_TO_INSTANCEGROUP.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RelatedAsset.RelationshipType"),
			},
			{
				Name:        "related_asset_name",
				Description: "The full name of the related asset.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RelatedAsset.Asset"),
			},
			{
				Name:        "related_asset_type",
				Description: "The type of the related asset.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RelatedAsset.AssetType"),
			},
			{
				Name:        "related_ancestors",
				Description: "The ancestors of the related asset in Google Cloud resource hierarchy.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("RelatedAsset.Ancestors"),
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

func listCloudAssetRelationships(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// The RELATIONSHIP content type is only available to Security Command
	// Center Premium and Enterprise tier customers; the API returns an
	// authorization error otherwise.
	return listCloudAssetsByContentType(ctx, d, h, "RELATIONSHIP", "gcp_cloud_asset_relationship.listCloudAssetRelationships", func(ctx context.Context, d *plugin.QueryData, asset *cloudasset.Asset) {
		d.StreamListItem(ctx, asset)
	})
}
