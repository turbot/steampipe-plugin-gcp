package gcp

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/billingbudgets/v1"
	"google.golang.org/api/cloudbilling/v1"
)

//// TABLE DEFINITION

func tableGcpBillingBudget(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_billing_budget",
		Description: "GCP Billing Budget",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "billing_account"}),
			Hydrate:    getBillingBudget,
			Tags:       map[string]string{"service": "billing", "action": "budgets.get"},
		},
		List: &plugin.ListConfig{
			KeyColumns:    plugin.OptionalColumns([]string{"billing_account"}),
			ParentHydrate: getBillingAccount,
			Hydrate:       listBillingBudgets,
			Tags:          map[string]string{"service": "billing", "action": "budgets.list"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The resource name of the budget.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Budget.Name").Transform(lastPathElement),
			},
			{
				Name:        "billing_account",
				Description: "The name given to the associated billing account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BillingAccount").Transform(lastPathElement),
			},
			{
				Name:        "display_name",
				Description: "The display name given to the budget.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Budget.DisplayName"),
			},

			{
				Name:        "budget_filter",
				Description: "Filters that define which resources are used to compute the actual spend against the budget amount.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Budget.BudgetFilter"),
			},
			{
				Name:        "last_period_amount",
				Description: "Use the last period's actual spend as the budget for the present period.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Budget.Amount.LastPeriodAmount"),
			},
			{
				Name:        "notifications_rule",
				Description: "Rules to apply to notifications sent based on budget spend and thresholds.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Budget.NotificationsRule"),
			},
			{
				Name:        "specified_amount",
				Description: "Use the last period's actual spend as the budget for the present period.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Budget.Amount.SpecifiedAmount"),
			},
			{
				Name:        "threshold_rules",
				Description: "Rules that trigger alerts (notifications of thresholds being crossed) when spend exceeds the specified percentages of the budget.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Budget.ThresholdRules"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Budget.DisplayName"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBillingBudgetAka,
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

type budgetInfo = struct {
	BillingAccount string
	Budget         *billingbudgets.GoogleCloudBillingBudgetsV1Budget
}

//// LIST FUNCTION

func listBillingBudgets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	acc := h.Item.(*cloudbilling.BillingAccount)

	// Validate - User input(if any) should match with the hydrated billing account
	if d.EqualsQualString("billing_account") != "" && "billingAccounts/"+d.EqualsQualString("billing_account") != acc.Name {
		return nil, nil
	}

	pageSize := types.Int64(100)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *pageSize {
			pageSize = limit
		}
	}

	// Create Service Connection
	service, err := BillingBudgetsService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_billing_budget.listBillingBudgets", "service_err", err)
		return nil, err
	}

	response := service.BillingAccounts.Budgets.List(acc.Name).PageSize(*pageSize)
	if err := response.Pages(
		ctx,
		func(page *billingbudgets.GoogleCloudBillingBudgetsV1ListBudgetsResponse) error {
			// apply rate limiting
			d.WaitForListRateLimit(ctx)

			for _, item := range page.Budgets {
				d.StreamListItem(ctx, budgetInfo{acc.Name, item})

				// Check if context has been cancelled or if the limit has been hit (if specified)
				// if there is a limit, it will return the number of rows required to reach this limit
				if d.RowsRemaining(ctx) == 0 {
					page.NextPageToken = ""
					return nil
				}
			}
			return nil
		},
	); err != nil {
		plugin.Logger(ctx).Error("gcp_billing_budget.listBillingBudgets", "api_err", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTION

func getBillingBudget(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["name"].GetStringValue()
	acc := d.EqualsQuals["billing_account"].GetStringValue()

	// Create Service Connection
	service, err := BillingBudgetsService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_billing_budget.getBillingBudget", "service_err", err)
		return nil, err
	}

	budget, err := service.BillingAccounts.Budgets.Get("billingAccounts/" + acc + "/budgets/" + name).Do()
	if err != nil {
		plugin.Logger(ctx).Error("gcp_billing_budget.getBillingBudget", "api_err", err)
		return nil, err
	}

	return budgetInfo{acc, budget}, nil
}

func getBillingBudgetAka(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_billing_budget.getBillingBudgetAka", "cache_err", err)
		return nil, err
	}

	project := projectId.(string)
	data := h.Item.(budgetInfo)
	akas := []string{"gcp://billingbudgets.googleapis.com/projects/" + project + "/" + data.Budget.Name}

	return akas, nil
}
