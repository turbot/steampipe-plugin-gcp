package gcp

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	"google.golang.org/api/monitoring/v3"
)

//// TABLE DEFINITION

func monitoringMetricColumns(columns []*plugin.Column) []*plugin.Column {
	return append(columns, commonMonitoringMetricColumns()...)
}

func commonMonitoringMetricColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "metadata",
			Description: "The associated monitored resource metadata.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "metric_type",
			Description: "The associated metric. A fully-specified metric used to identify the time series.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Metric.Type"),
		},
		{
			Name:        "metric_labels",
			Description: "The set of label values that uniquely identify this metric.",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Metric.Labels"),
		},
		{
			Name:        "metric_kind",
			Description: "The metric type.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "maximum",
			Description: "The maximum metric value for the data point.",
			Hydrate:     metricstatistic,
			Type:        proto.ColumnType_DOUBLE,
		},
		{
			Name:        "minimum",
			Description: "The minimum metric value for the data point.",
			Hydrate:     metricstatistic,
			Type:        proto.ColumnType_DOUBLE,
		},
		{
			Name:        "average",
			Description: "The average of the metric values that correspond to the data point.",
			Hydrate:     metricstatistic,
			Type:        proto.ColumnType_DOUBLE,
		},
		{
			Name:        "sample_count",
			Description: "The number of metric values that contributed to the aggregate value of this data point.",
			Hydrate:     metricstatistic,
			Type:        proto.ColumnType_DOUBLE,
		},
		{
			Name:        "sum",
			Description: "The sum of the metric values for the data point.",
			Hydrate:     metricstatistic,
			Type:        proto.ColumnType_DOUBLE,
		},
		{
			Name:        "resource",
			Description: "The associated monitored resource.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "unit",
			Description: "The data points of this time series. When listing time series, points are returned in reverse time order.When creating a time series, this field must contain exactly one point and the point's type must be the same as the value type of the associated metric. If the associated metric's descriptor must be auto-created, then the value type of the descriptor is determined by the point's type, which must be BOOL, INT64, DOUBLE, or DISTRIBUTION.",
			Type:        proto.ColumnType_STRING,
		},
	}
}

type monitorMetric struct {
	// Resource  Name
	DimensionValue string

	// The associated monitored resource metadata.
	Metadata *monitoring.MonitoredResourceMetadata

	// The associated metric. A fully-specified metric used to identify the time series.
	Metric *monitoring.Metric

	//Possible values:
	//   "METRIC_KIND_UNSPECIFIED" - Do not use this default value.
	//   "GAUGE" - An instantaneous measurement of a value.
	//   "DELTA" - The change in a value during a time interval.
	//   "CUMULATIVE" - A value accumulated over a time interval.
	MetricKind string

	// The maximum metric value for the data point.
	Maximum *float64

	// The minimum metric value for the data point.
	Minimum *float64

	// The average of the metric values that correspond to the data point.
	Average *float64

	// The data points of this time series. When listing time series, points are returned in reverse time order.When creating a time series, this field must contain exactly one point and the point's type must be the same as the value type of the associated metric. If the associated metric's descriptor must be auto-created, then the value type of the descriptor is determined by the point's type, which must be BOOL, INT64, DOUBLE, or DISTRIBUTION.
	Points []*monitoring.Point

	// The number of metric values that contributed to the aggregate value of this data point.
	SampleCount *float64

	// The sum of the metric values for the data point.
	Sum *float64

	// The associated monitored resource.
	Resource *monitoring.MonitoredResource

	// The units in which the metric value is reported. It is only applicable if the value_type is INT64, DOUBLE, or DISTRIBUTION. The unit defines the representation of the stored metric values.
	Unit string
}

func getMonitoringStartDateForGranularity(granularity string) time.Time {
	switch strings.ToUpper(granularity) {
	case "DAILY":
		// 1 year
		return time.Now().AddDate(-1, 0, 0)
	case "HOURLY":
		// 60 days
		return time.Now().AddDate(0, 0, -60)
	}
	// else 5 days
	return time.Now().AddDate(0, 0, -5)
}

func getMonitoringPeriodForGranularity(granularity string) string {
	switch strings.ToUpper(granularity) {
	case "DAILY":
		// 24 hours
		return "86400s"
	case "HOURLY":
		// 1 hour
		return "3600s"
	}
	// else 5 minutes
	return "300s"
}

func listMonitorMetricStatistics(ctx context.Context, d *plugin.QueryData, granularity string, metricType string, dimentionKey string, dimensionValue string, resourceName string) (*monitoring.ListTimeSeriesResponse, error) {

	plugin.Logger(ctx).Trace("listMonitorMetricStatistics")
	// Create Service Connection
	service, err := MonitoringService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}
	project := projectData.Project

	endTime := time.Now().Format(time.RFC3339)
	startTime := getMonitoringStartDateForGranularity(granularity).Format(time.RFC3339)

	period := getMonitoringPeriodForGranularity(granularity)

	filterString := "metric.type = " + metricType + " AND " + dimentionKey + dimensionValue

	resp := service.Projects.TimeSeries.List("projects/" + project).Filter(filterString).IntervalStartTime(startTime).IntervalEndTime(endTime).AggregationAlignmentPeriod(period)
	if err := resp.Pages(ctx, func(page *monitoring.ListTimeSeriesResponse) error {
		for _, metric := range page.TimeSeries {
			d.StreamLeafListItem(ctx, &monitorMetric{
				DimensionValue: strings.ReplaceAll(dimensionValue, "\"", ""),
				Metadata:       metric.Metadata,
				Metric:         metric.Metric,
				MetricKind:     metric.MetricKind,
				Points:         metric.Points,
				Resource:       metric.Resource,
				Unit:           metric.Unit,
			})
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

// Get metric statistic

func metricstatistic(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	points := h.Item.(*monitorMetric).Points
	var pointValues []float64
	for _, value := range points {
		pointValueType := value.Value

		// TODO: Need to handle BoolType, StringValue and DistributionValue

		if pointValueType.DoubleValue != nil {
			// val := strconv.FormatFloat(*pointValueType.DoubleValue, 'E', -1, 64)
			pointValues = append(pointValues, *pointValueType.DoubleValue)
		}

		if pointValueType.Int64Value != nil {
			// val := strconv.FormatInt(*pointValueType.Int64Value, 10)
			val := float64(*pointValueType.Int64Value)
			pointValues = append(pointValues, val)
		}

		if pointValueType.StringValue != nil {
			val, err := strconv.ParseFloat(*pointValueType.StringValue, 64)
			if err != nil {
				return nil, err
			}
			pointValues = append(pointValues, val)
		}

	}

	minValue := pointValues[0]
	maxValue := minValue

	m := make(map[string]float64)
	m["Maximum"] = maxValue
	m["Minimum"] = minValue
	m["Sum"] = float64(0)
	m["Average"] = float64(0)
	m["SampleCount"] = float64(0)
	
	sum := float64(0)
	// plugin.Logger(ctx).Trace("All points ===>", pointValues)

	for _, point := range pointValues {

		if point > maxValue {
			maxValue = point
			m["Maximum"] = point
		}
		if point < minValue {
			minValue = point
			m["Minimum"] = point
		}
		sum += point
	}
	m["Sum"] = sum

	sampleCount := float64(len(pointValues))
	m["SampleCount"] = sampleCount

	average := sum / sampleCount
	m["Average"] = average

	return m, nil

}
