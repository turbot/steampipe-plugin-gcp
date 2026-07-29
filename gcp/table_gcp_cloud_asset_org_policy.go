package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"

	"google.golang.org/api/cloudasset/v1"
)

type cloudAssetOrgPolicy struct {
	Name       string
	AssetType  string
	Ancestors  []string
	UpdateTime string
	Policy     *cloudasset.GoogleCloudOrgpolicyV1Policy
}

//// TABLE DEFINITION

func tableGcpCloudAssetOrgPolicy(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_cloud_asset_org_policy",
		Description: "GCP Cloud Asset Organization Policy",
		List: &plugin.ListConfig{
			Hydrate: listCloudAssetOrgPolicies,
			Tags:    map[string]string{"service": "cloudasset", "action": "assets.listOrgPolicy"},
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
				Name:        "constraint_name",
				Description: "The name of the constraint the policy is configuring, e.g., constraints/serviceuser.services.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Policy.Constraint"),
			},
			{
				Name:        "policy_version",
				Description: "Version of the policy.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Policy.Version"),
			},
			{
				Name:        "policy_update_time",
				Description: "The time stamp this policy was previously updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Policy.UpdateTime"),
			},
			{
				Name:        "etag",
				Description: "An opaque tag indicating the current version of the policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Policy.Etag"),
			},
			{
				Name:        "list_policy",
				Description: "Policy for data governed by a list of values.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Policy.ListPolicy"),
			},
			{
				Name:        "boolean_policy",
				Description: "For boolean constraints, whether to enforce the constraint.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Policy.BooleanPolicy"),
			},
			{
				Name:        "restore_default",
				Description: "Restores the default behavior of the constraint; independent of constraint type.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Policy.RestoreDefault"),
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

func listCloudAssetOrgPolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// An asset can have multiple organization policies set directly on it, one
	// per constraint; stream one row per policy.
	return listCloudAssetsByContentType(ctx, d, h, "ORG_POLICY", "gcp_cloud_asset_org_policy.listCloudAssetOrgPolicies", func(ctx context.Context, d *plugin.QueryData, asset *cloudasset.Asset) {
		for _, policy := range asset.OrgPolicy {
			d.StreamListItem(ctx, cloudAssetOrgPolicy{
				Name:       asset.Name,
				AssetType:  asset.AssetType,
				Ancestors:  asset.Ancestors,
				UpdateTime: asset.UpdateTime,
				Policy:     policy,
			})
		}
	})
}
