package gcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"
)

//// TABLE DEFINITION

func tableGcpMonitoringMetricStatistic(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name: "gcp_monitoring_metric_statistic",
		List: &plugin.ListConfig{
			Hydrate: listGcpMonitoringMetricStatistic,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "dimension_key", Require: plugin.Required, CacheMatch: query_cache.CacheMatchExact},
				{Name: "dimension_value", Require: plugin.Required, CacheMatch: query_cache.CacheMatchExact},
				{Name: "granularity", Require: plugin.Required, CacheMatch: query_cache.CacheMatchExact},
				{Name: "metric_type", Require: plugin.Required, CacheMatch: query_cache.CacheMatchExact},
				{Name: "location", Require: plugin.Optional},
			},
		},
		Columns: monitoringMetricColumns([]*plugin.Column{
			{
				Name:        "dimension_key",
				Description: "The name of the disk.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("dimension_key"),
			},
			{
				Name:        "granularity",
				Description: "The name of the disk.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("granularity"),
			},
			{
				Name:        "dimension_value",
				Description: "The name of the disk.",
				Type:        proto.ColumnType_STRING,
			},
		}),
	}
}

//// LIST FUNCTION

func listGcpMonitoringMetricStatistic(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	inputParam, err := buildMetricStatisticInputParam(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_monitoring_metric_statistic.listGcpMonitoringMetricStatistic", "param_error", err)
		return nil, err
	}

	return listMonitorMetricStatistics(ctx, d, h, inputParam.Granularity, inputParam.MetricType, inputParam.DimKey, inputParam.DimValue, "", inputParam.Location)
}

//// UTILITY FUNCTION

type MetricStatisticInput struct {
	DimKey      string
	DimValue    string
	Granularity string
	MetricType  string
	Location    string
}

// Build metric static input param
func buildMetricStatisticInputParam(_ context.Context, d *plugin.QueryData) (MetricStatisticInput, error) {
	dimKey := d.EqualsQualString("dimension_key")
	dimValue := d.EqualsQualString("dimension_value")
	granularity := d.EqualsQualString("granularity")
	metricType := d.EqualsQualString("metric_type")

	return MetricStatisticInput{
		DimKey:      fmt.Sprintf("%s = ", dimKey),
		DimValue:    fmt.Sprintf("\"%s\"", dimValue),
		Granularity: strings.ToUpper(granularity),
		MetricType:  fmt.Sprintf("\"%s\"", metricType),
		Location:    d.EqualsQualString("location"),
	}, nil
}
