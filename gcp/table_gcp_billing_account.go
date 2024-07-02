package gcp

import (
	"context"

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
			Tags:    map[string]string{"service": "billing", "action": "accounts.get"},
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

func getBillingAccount(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var accountName string

	// Create Service Connection
	service, err := BillingService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_billing_account.getBillingAccount", "service_err", err)
		return nil, err
	}

	// Validate whether a user has provided any input, then skip the `GetBillingInfo` api
	if d.EqualsQualString("name") != "" {
		accountName = "billingAccounts/" + d.EqualsQualString("name")
	} else {

		// Fetch BillingInfo for the project, to get the billing account name
		// Get project details

		projectId, err := getProject(ctx, d, h)
		if err != nil {
			plugin.Logger(ctx).Error("gcp_billing_account.getBillingAccount", "cache_err", err)
			return nil, err
		}
		project := projectId.(string)

		resp, err := service.Projects.GetBillingInfo("projects/" + project).Do()
		if err != nil {
			plugin.Logger(ctx).Error("gcp_billing_account.getBillingAccount.GetBillingInfo", "api_err", err)
			return nil, err
		}

		if resp != nil && resp.BillingAccountName != "" {
			accountName = resp.BillingAccountName
		}
	}

	accResponse, err := service.BillingAccounts.Get(accountName).Do()
	if err != nil {
		plugin.Logger(ctx).Error("gcp_billing_account.getBillingAccount", "api_err", err)
		return nil, err
	}

	d.StreamListItem(ctx, accResponse)

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

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}

	project := projectId.(string)
	acc := h.Item.(*cloudbilling.BillingAccount)
	akas := []string{"gcp://cloudbilling.googleapis.com/projects/" + project + "/" + acc.Name}

	return akas, nil
}
