package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	billingbudgets "google.golang.org/api/billingbudgets/v1beta1"
)

func tableGcpBillingAndUsage(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_cost_usage",
		Description: "GCP Billing - Billing and Usage",
		List: &plugin.ListConfig{
			Hydrate: listCostAndUsage,
		},
		Columns: costExplorerColumns([]*plugin.Column{
			{
				Name:        "dimension_1",
				Description: "Valid values are AZ, INSTANCE_TYPE, LINKED_ACCOUNT, OPERATION, PURCHASE_TYPE, SERVICE, USAGE_TYPE, PLATFORM, TENANCY, RECORD_TYPE, LEGAL_ENTITY_NAME, DEPLOYMENT_OPTION, DATABASE_ENGINE, CACHE_ENGINE, INSTANCE_TYPE_FAMILY, REGION, BILLING_ENTITY, RESERVATION_ID, SAVINGS_PLANS_TYPE, SAVINGS_PLAN_ARN, OPERATING_SYSTEM",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "dimension_2",
				Description: "Valid values are AZ, INSTANCE_TYPE, LINKED_ACCOUNT, OPERATION, PURCHASE_TYPE, SERVICE, USAGE_TYPE, PLATFORM, TENANCY, RECORD_TYPE, LEGAL_ENTITY_NAME, DEPLOYMENT_OPTION, DATABASE_ENGINE, CACHE_ENGINE, INSTANCE_TYPE_FAMILY, REGION, BILLING_ENTITY, RESERVATION_ID, SAVINGS_PLANS_TYPE, SAVINGS_PLAN_ARN, OPERATING_SYSTEM",
				Type:        proto.ColumnType_STRING,
			},

			// Quals columns - to filter the lookups
			{
				Name:        "granularity",
				Description: "",
				Type:        proto.ColumnType_STRING,
				Hydrate:     hydrateCostAndUsageQuals,
			},
			{
				Name:        "search_start_time",
				Description: "",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     hydrateCostAndUsageQuals,
			},
			{
				Name:        "search_end_time",
				Description: "",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     hydrateCostAndUsageQuals,
			},
			{
				Name:        "dimension_type_1",
				Description: "",
				Type:        proto.ColumnType_STRING,
				Hydrate:     hydrateCostAndUsageQuals,
			},
			{
				Name:        "dimension_type_2",
				Description: "",
				Type:        proto.ColumnType_STRING,
				Hydrate:     hydrateCostAndUsageQuals,
			},
		}),
	}
}

func listCostAndUsage(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	params := buildInputFromQuals(d.KeyColumnQuals)

	return streamCostAndUsage(ctx, d, params)
}

func buildInputFromQuals(keyQuals map[string]*proto.QualValue) *billingbudgets.GoogleCloudBillingBudgetsV1beta1Filter {
	// granularity := strings.ToUpper(keyQuals["granularity"].GetStringValue())
	// timeFormat := time.RFC3339
	// if granularity == "HOURLY" {
	// 	timeFormat = "2006-01-02T15:04:05Z"
	// }
	// endTime := time.Now().Format(timeFormat)
	// startTime := getCEStartDateForGranularity(granularity)

	// dim1 := strings.ToUpper(keyQuals["dimension_type_1"].GetStringValue())
	// dim2 := strings.ToUpper(keyQuals["dimension_type_2"].GetStringValue())

	params := &billingbudgets.GoogleCloudBillingBudgetsV1beta1Filter{
		CustomPeriod: &billingbudgets.GoogleCloudBillingBudgetsV1beta1CustomPeriod{
			StartDate: &billingbudgets.GoogleTypeDate{
				Year:  2021,
				Month: 6,
				Day:   1,
			},
			EndDate: &billingbudgets.GoogleTypeDate{
				Year:  2021,
				Month: 7,
				Day:   1,
			},
		},
	}

	// params := &costexplorer.GetCostAndUsageInput{
	// 	TimePeriod: &costexplorer.DateInterval{
	// 		Start: aws.String(startTime),
	// 		End:   aws.String(endTime),
	// 	},
	// 	Granularity: aws.String(granularity),
	// 	Metrics:     aws.StringSlice(AllCostMetrics()),
	// }
	// var groupings []*costexplorer.GroupDefinition
	// if dim1 != "" {
	// 	groupings = append(groupings, &costexplorer.GroupDefinition{
	// 		Type: aws.String("DIMENSION"),
	// 		Key:  aws.String(dim1),
	// 	})
	// }
	// if dim2 != "" {
	// 	groupings = append(groupings, &costexplorer.GroupDefinition{
	// 		Type: aws.String("DIMENSION"),
	// 		Key:  aws.String(dim2),
	// 	})
	// }
	// params.SetGroupBy(groupings)

	return params
}
