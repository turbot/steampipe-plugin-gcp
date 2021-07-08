package gcp

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	billingbudgets "google.golang.org/api/billingbudgets/v1beta1"
)

var costExplorerColumnDefs = []*plugin.Column{
	{
		Name:        "period_start",
		Description: "Start timestamp for this cost metric.",
		Type:        proto.ColumnType_TIMESTAMP,
	},
	{
		Name:        "period_end",
		Description: "End timestamp for this cost metric.",
		Type:        proto.ColumnType_TIMESTAMP,
	},

	{
		Name:        "estimated",
		Description: "Whether the result is estimated.",
		Type:        proto.ColumnType_BOOL,
	},
	{
		Name:        "blended_cost_amount",
		Description: "This cost metric reflects the average cost of usage across the consolidated billing family. If you use the consolidated billing feature in AWS Organizations, you can view costs using blended rates.",
		Type:        proto.ColumnType_DOUBLE,
	},
	{
		Name:        "blended_cost_unit",
		Description: "Unit type for blended costs.",
		Type:        proto.ColumnType_STRING,
	},

	{
		Name:        "unblended_cost_amount",
		Description: "Unblended costs represent your usage costs on the day they are charged to you. In finance terms, they represent your costs on a cash basis of accounting.",
		Type:        proto.ColumnType_DOUBLE,
	},
	{
		Name:        "unblended_cost_unit",
		Description: "Unit type for unblended costs.",
		Type:        proto.ColumnType_STRING,
	},

	{
		Name:        "net_unblended_cost_amount",
		Description: "This cost metric reflects the unblended cost after discounts.",
		Type:        proto.ColumnType_DOUBLE,
	},
	{
		Name:        "net_unblended_cost_unit",
		Description: "Unit type for net unblended costs.",
		Type:        proto.ColumnType_STRING,
	},

	{
		Name:        "amortized_cost_amount",
		Description: "This cost metric reflects the effective cost of the upfront and monthly reservation fees spread across the billing period. By default, Cost Explorer shows the fees for Reserved Instances as a spike on the day that you're charged, but if you choose to show costs as amortized costs, the costs are amortized over the billing period. This means that the costs are broken out into the effective daily rate. AWS estimates your amortized costs by combining your unblended costs with the amortized portion of your upfront and recurring reservation fees.",
		Type:        proto.ColumnType_DOUBLE,
	},
	{
		Name:        "amortized_cost_unit",
		Description: "Unit type for amortized costs.",
		Type:        proto.ColumnType_STRING,
	},

	{
		Name:        "net_amortized_cost_amount",
		Description: "This cost metric amortizes the upfront and monthly reservation fees while including discounts such as RI volume discounts.",
		Type:        proto.ColumnType_DOUBLE,
	},
	{
		Name:        "net_amortized_cost_unit",
		Description: "Unit type for net amortized costs.",
		Type:        proto.ColumnType_STRING,
	},

	{
		Name:        "usage_quantity_amount",
		Description: "The amount of usage that you incurred. NOTE: If you return the UsageQuantity metric, the service aggregates all usage numbers without taking into account the units. For example, if you aggregate usageQuantity across all of Amazon EC2, the results aren't meaningful because Amazon EC2 compute hours and data transfer are measured in different units (for example, hours vs. GB).",
		Type:        proto.ColumnType_DOUBLE,
	},
	{
		Name:        "usage_quantity_unit",
		Description: "Unit type for usage quantity.",
		Type:        proto.ColumnType_STRING,
	},

	{
		Name:        "normalized_usage_amount",
		Description: "The amount of usage that you incurred, in normalized units, for size-flexible RIs. The NormalizedUsageAmount is equal to UsageAmount multiplied by NormalizationFactor.",
		Type:        proto.ColumnType_DOUBLE,
	},
	{
		Name:        "normalized_usage_unit",
		Description: "Unit type for normalized usage.",
		Type:        proto.ColumnType_STRING,
	},
}

// append the common aws cost explorer columns onto the column list
func costExplorerColumns(columns []*plugin.Column) []*plugin.Column {
	return append(columns, costExplorerColumnDefs...)
}

//// LIST FUNCTION
// , params *costexplorer.GetCostAndUsageInput

func streamCostAndUsage(ctx context.Context, d *plugin.QueryData, params *billingbudgets.GoogleCloudBillingBudgetsV1beta1Filter) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("streamCostAndUsage")

	// Create session
	svc, err := BillingBudgetService(ctx, d)
	if err != nil {
		return nil, err
	}

	// cloudBillingSvc, err := CloudBillingService(ctx, d)

	services, err := svc.BillingAccounts.Budgets.List("njka").Do()

	// output, err :=

	// List call
	// for {
	// 	output, err := svc.(params)
	// 	if err != nil {
	// 		logger.Error("streamCostAndUsage", "err", err)
	// 		return nil, err
	// 	}

	// 	// stream the results...
	// 	for _, row := range buildCEMetricRows(ctx, output, d.KeyColumnQuals) {
	// 		d.StreamListItem(ctx, row)
	// 	}

	// 	// get more pages if there are any...
	// 	if output.NextPageToken == nil {
	// 		break
	// 	}
	// 	params.SetNextPageToken(*output.NextPageToken)
	// }

	return nil, nil
}

type CEQuals struct {
	// Quals stuff
	SearchStartTime *timestamp.Timestamp
	SearchEndTime   *timestamp.Timestamp
	Granularity     string
	DimensionType1  string
	DimensionType2  string
}

func hydrateCostAndUsageQuals(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("hydrateKeyQuals")
	//plugin.Logger(ctx).Warn("hydrateKeyQuals", "d.KeyColumnQuals", d.KeyColumnQuals)

	return &CEQuals{
		SearchStartTime: d.KeyColumnQuals["search_start_time"].GetTimestampValue(),
		SearchEndTime:   d.KeyColumnQuals["search_end_time"].GetTimestampValue(),
		Granularity:     d.KeyColumnQuals["granularity"].GetStringValue(),
		DimensionType1:  d.KeyColumnQuals["dimension_type_1"].GetStringValue(),
		DimensionType2:  d.KeyColumnQuals["dimension_type_2"].GetStringValue(),
	}, nil
}

func getCEStartDateForGranularity(granularity string) string {
	switch granularity {
	case "DAILY", "MONTHLY":
		// 1 year
		return time.Now().AddDate(-1, 0, 0).Format(time.RFC3339)
	case "HOURLY":
		// 13 days
		return time.Now().AddDate(0, 0, -13).Format(time.RFC3339)
	}
	return time.Now().AddDate(0, 0, -13).Format(time.RFC3339)
}
