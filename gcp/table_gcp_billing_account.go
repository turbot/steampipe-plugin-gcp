package gcp

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/cloudbilling/v1"
)

//// TABLE DEFINITION

func tableGcpBillingAccount(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_billing_account",
		Description: "GCP Billing Account",
		List: &plugin.ListConfig{
			KeyColumns: plugin.OptionalColumns([]string{"name"}),
			Hydrate:    getBillingAccount,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The resource name of the billing account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel().Transform(lastPathElement),
			},
			{
				Name:        "display_name",
				Description: "The display name given to the billing account. This name is displayed in the Google Cloud Console.",
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
			{
				Name:        "iam_policy",
				Description: "An IAM policy, which specifies access controls for the billing account.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBillingAccountIamPolicy,
				Transform:   transform.FromValue(),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBillingAccountAka,
				Transform:   transform.FromValue(),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromConstant("global"),
			},
		},
	}
}

//// LIST FUNCTION

func getBillingAccount(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create Service Connection
	service, err := BillingService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_billing_account.getBillingAccount", "service_err", err)
		return nil, err
	}

	// If a user has provided a billing account name, get it instead of listing all accounts
	if d.EqualsQualString("name") != "" {
		accountName := "billingAccounts/" + d.EqualsQualString("name")

		resp, err := service.BillingAccounts.Get(accountName).Do()
		if err != nil {
			plugin.Logger(ctx).Error("gcp_billing_account.getBillingAccount.get", "api_err", err)
			return nil, err
		}

		d.StreamListItem(ctx, resp)
	} else {
		// Max limit is set as per documentation
		pageSize := types.Int64(100)
		limit := d.QueryContext.Limit
		if d.QueryContext.Limit != nil {
			if *limit < *pageSize {
				pageSize = limit
			}
		}
		resp := service.BillingAccounts.List().PageSize(*pageSize)
		if err := resp.Pages(ctx, func(page *cloudbilling.ListBillingAccountsResponse) error {
			for _, account := range page.BillingAccounts {
				d.StreamListItem(ctx, account)

				// Check if context has been cancelled or if the limit has been hit (if specified)
				// if there is a limit, it will return the number of rows required to reach this limit
				if d.RowsRemaining(ctx) == 0 {
					page.NextPageToken = ""
					return nil
				}
			}
			return nil
		}); err != nil {
			plugin.Logger(ctx).Error("gcp_billing_account.getBillingAccount.list", "api_err", err)
			return nil, err
		}
	}

	return nil, nil
}

func getBillingAccountIamPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	acc := h.Item.(*cloudbilling.BillingAccount)

	// Create Service Connection
	service, err := BillingService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_billing_account.getBillingAccountIamPolicy", "service_err", err)
		return nil, err
	}

	policy, err := service.BillingAccounts.GetIamPolicy(acc.Name).Do()
	if err != nil {
		plugin.Logger(ctx).Error("gcp_billing_account.getBillingAccountIamPolicy", "api_err", err)
		return nil, err
	}

	return policy, nil
}

func getBillingAccountAka(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	acc := h.Item.(*cloudbilling.BillingAccount)
	akas := []string{"gcp://cloudbilling.googleapis.com/" + acc.Name}

	return akas, nil
}
