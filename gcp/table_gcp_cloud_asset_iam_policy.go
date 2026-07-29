package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"

	"google.golang.org/api/cloudasset/v1"
)

//// TABLE DEFINITION

func tableGcpCloudAssetIamPolicy(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_cloud_asset_iam_policy",
		Description: "GCP Cloud Asset IAM Policy",
		List: &plugin.ListConfig{
			Hydrate: listCloudAssetIamPolicies,
			Tags:    map[string]string{"service": "cloudasset", "action": "assets.listIamPolicy"},
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
				Name:        "policy_version",
				Description: "Specifies the format of the policy.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("IamPolicy.Version"),
			},
			{
				Name:        "etag",
				Description: "Etag used for optimistic concurrency control of the policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("IamPolicy.Etag"),
			},
			{
				Name:        "bindings",
				Description: "Associates a list of members, or principals, with a role.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("IamPolicy.Bindings"),
			},
			{
				Name:        "audit_configs",
				Description: "Specifies cloud audit logging configuration for the policy.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("IamPolicy.AuditConfigs"),
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

func listCloudAssetIamPolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return listCloudAssetsByContentType(ctx, d, h, "IAM_POLICY", "gcp_cloud_asset_iam_policy.listCloudAssetIamPolicies", func(ctx context.Context, d *plugin.QueryData, asset *cloudasset.Asset) {
		d.StreamListItem(ctx, asset)
	})
}
