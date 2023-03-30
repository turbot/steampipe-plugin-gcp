package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableGcpBillingAccount(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_billing_account",
		Description: "GCP Billing Account",
		List: &plugin.ListConfig{
			Hydrate: getGcpBillingAccount,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The resource name for project billing account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "The display name given to the billing account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "master_billing_account",
				Description: "The resource name of the parent billing account, if any.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "open",
				Description: "Whether the billing account is open, and will therefore be charged for any usage on associated projects.",
				Type:        proto.ColumnType_BOOL,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromConstant("global"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Hydrate:     plugin.HydrateFunc(getProject).WithCache(),
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func getGcpBillingAccount(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Get project details
	getProjectCached := plugin.HydrateFunc(getProject).WithCache()
	projectId, err := getProjectCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	// Create Service Connection
	service, err := BillingService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Debug("gcp_billing_account.getGcpBillingAccount."+project, "service_err", err)
		return nil, err
	}

	resp, err := service.Projects.GetBillingInfo("projects/" + project).Do()
	if err != nil {
		plugin.Logger(ctx).Debug("gcp_billing_account.getGcpBillingAccount."+project, "api_err", err)
		return nil, err
	}

	if resp != nil && resp.BillingAccountName != "" {
		accResponse, err := service.BillingAccounts.Get(resp.BillingAccountName).Do()
		if err != nil {
			plugin.Logger(ctx).Debug("gcp_billing_account.getGcpBillingAccount."+project, "api_err", err)
			return nil, err
		}

		d.StreamListItem(ctx, accResponse)
	}

	return nil, nil
}
