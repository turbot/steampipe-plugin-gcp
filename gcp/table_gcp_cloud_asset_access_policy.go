package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"

	"google.golang.org/api/cloudasset/v1"
)

//// TABLE DEFINITION

func tableGcpCloudAssetAccessPolicy(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_cloud_asset_access_policy",
		Description: "GCP Cloud Asset Access Policy",
		List: &plugin.ListConfig{
			Hydrate: listCloudAssetAccessPolicies,
			Tags:    map[string]string{"service": "cloudasset", "action": "assets.listAccessPolicy"},
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
				Name:        "access_policy",
				Description: "An access policy is a container for all of your Access Context Manager resources. Only set for accesscontextmanager.googleapis.com/AccessPolicy assets.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "access_level",
				Description: "Access levels are used for permitting access to resources based on contextual information about the request. Only set for accesscontextmanager.googleapis.com/AccessLevel assets.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "service_perimeter",
				Description: "A service perimeter describes a set of Google Cloud resources which can freely import and export data amongst themselves. Only set for accesscontextmanager.googleapis.com/ServicePerimeter assets.",
				Type:        proto.ColumnType_JSON,
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

func listCloudAssetAccessPolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return listCloudAssetsByContentType(ctx, d, h, "ACCESS_POLICY", "gcp_cloud_asset_access_policy.listCloudAssetAccessPolicies", func(ctx context.Context, d *plugin.QueryData, asset *cloudasset.Asset) {
		d.StreamListItem(ctx, asset)
	})
}
