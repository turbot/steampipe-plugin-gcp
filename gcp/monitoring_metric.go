package gcp

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/monitoring/v3"
)

//// TABLE DEFINITION

func monitoringMetricColumns(columns []*plugin.Column) []*plugin.Column {
	return append(columns, commonMonitoringMetricColumns()...)
}

func commonMonitoringMetricColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "metric_type",
			Description: "The associated metric. A fully-specified metric used to identify the time series.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Metric.Type"),
		},
		{
			Name:        "metric_kind",
			Description: "The metric type.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "maximum",
			Description: "The maximum metric value for the data point.",
			Type:        proto.ColumnType_DOUBLE,
		},
		{
			Name:        "minimum",
			Description: "The minimum metric value for the data point.",
			Type:        proto.ColumnType_DOUBLE,
		},
		{
			Name:        "average",
			Description: "The average of the metric values that correspond to the data point.",
			Type:        proto.ColumnType_DOUBLE,
		},
		{
			Name:        "sample_count",
			Description: "The number of metric values that contributed to the aggregate value of this data point.",
			Type:        proto.ColumnType_DOUBLE,
		},
		{
			Name:        "sum",
			Description: "The sum of the metric values for the data point.",
			Type:        proto.ColumnType_DOUBLE,
		},
		{
			Name:        "timestamp",
			Description: "The time stamp used for the data point.",
			Type:        proto.ColumnType_TIMESTAMP,
			Transform:   transform.FromField("TimeStamp"),
		},
		{
			Name:        "unit",
			Description: "The data points of this time series. When listing time series, points are returned in reverse time order.When creating a time series, this field must contain exactly one point and the point's type must be the same as the value type of the associated metric. If the associated metric's descriptor must be auto-created, then the value type of the descriptor is determined by the point's type, which must be BOOL, INT64, DOUBLE, or DISTRIBUTION.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "metadata",
			Description: "The associated monitored resource metadata.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "metric_labels",
			Description: "The set of label values that uniquely identify this metric.",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Metric.Labels"),
		},
		{
			Name:        "resource",
			Description: "The associated monitored resource.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "location",
			Description: ColumnDescriptionLocation,
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "project",
			Description: ColumnDescriptionProject,
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

	// The time stamp used for the data point.
	TimeStamp *string

	// The associated monitored resource.
	Resource *monitoring.MonitoredResource

	// The units in which the metric value is reported. It is only applicable if the value_type is INT64, DOUBLE, or DISTRIBUTION. The unit defines the representation of the stored metric values.
	Unit string

	// The GCP multi-region, region, or zone in which the resource is located.
	Location string

	// The GCP Project in which the resource is located.
	Project string
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

func getIncrementalTimeAsPerGranularity(granularity string) time.Duration {
	switch granularity {
	case "DAILY":
		return 86400
	case "HOURLY":
		return 3600
	default:
		return 300
	}
}

func listMonitorMetricStatistics(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, granularity string, metricType string, dimensionKey string, dimensionValue string, resourceName string, location string) (*monitoring.ListTimeSeriesResponse, error) {
	plugin.Logger(ctx).Trace("listMonitorMetricStatistics")

	// Create Service Connection
	service, err := MonitoringService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	endTime := time.Now().Format(time.RFC3339)
	startTime := getMonitoringStartDateForGranularity(granularity).Format(time.RFC3339)
	period := getMonitoringPeriodForGranularity(granularity)

	filterString := "metric.type = " + metricType + " AND " + dimensionKey + dimensionValue

	resp := service.Projects.TimeSeries.List("projects/" + project).Filter(filterString).IntervalStartTime(startTime).IntervalEndTime(endTime).AggregationAlignmentPeriod(period)
	if err := resp.Pages(ctx, func(page *monitoring.ListTimeSeriesResponse) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, metric := range page.TimeSeries {
			statistics, _ := metricstatistic(granularity, metric.Points, ctx)
			for _, statistic := range statistics {
				d.StreamLeafListItem(ctx, &monitorMetric{
					DimensionValue: strings.ReplaceAll(dimensionValue, "\"", ""),
					Metadata:       metric.Metadata,
					Metric:         metric.Metric,
					MetricKind:     metric.MetricKind,
					Points:         metric.Points,
					Maximum:        &statistic.Maximum,
					Minimum:        &statistic.Minimum,
					Average:        &statistic.Average,
					SampleCount:    &statistic.SampleCount,
					Sum:            &statistic.Sum,
					TimeStamp:      &statistic.TimeStamp,
					Resource:       metric.Resource,
					Unit:           metric.Unit,
					Location:       location,
					Project:        project,
				})
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

type PointWithTimeStamp struct {
	// Point Value
	Point float64

	// Time stamp of the point value
	TimeStamp string
}

type Statistics struct {
	Maximum     float64
	Minimum     float64
	Sum         float64
	Average     float64
	SampleCount float64
	TimeStamp   string
}

// Get metric statistic
func metricstatistic(granularity string, points []*monitoring.Point, ctx context.Context) ([]*Statistics, error) {
	var pointValues []*PointWithTimeStamp
	var statistics []*Statistics

	// Form an array with required data of points
	for _, value := range points {
		pointValueType := value.Value
		timeStamp := value.Interval.StartTime

		// TODO: Need to handle BoolType, StringValue and DistributionValue
		if pointValueType.DoubleValue != nil {
			pointValues = append(pointValues, &PointWithTimeStamp{Point: *pointValueType.DoubleValue, TimeStamp: timeStamp})
		}

		if pointValueType.Int64Value != nil {
			val := float64(*pointValueType.Int64Value)
			pointValues = append(pointValues, &PointWithTimeStamp{Point: val, TimeStamp: timeStamp})
		}

		if pointValueType.StringValue != nil {
			val, err := strconv.ParseFloat(*pointValueType.StringValue, 64)
			if err != nil {
				return nil, err
			}
			pointValues = append(pointValues, &PointWithTimeStamp{Point: val, TimeStamp: timeStamp})
		}
	}

	// Initialize max and min value with first point value
	var sum, average, sampleCount float64
	minValue := pointValues[0].Point
	maxValue := minValue

	startTime := pointValues[0].TimeStamp
	var timeDiff float64
	var pointCount, pointIndex int64
	var diffCheckExecuted bool

	// Iterate the points value and extract statistics
	for _, point := range pointValues {
		timeDiff = checkTimeDiff(point.TimeStamp, startTime)
		plugin.Logger(ctx).Trace("Time Diff", timeDiff)

		// Check time duration between start time and current point time stamp
		interval, _ := strconv.ParseFloat(strings.ReplaceAll(getMonitoringPeriodForGranularity(granularity), "s", ""), 64)
		diffCheckExecuted = false

		// Check time diff(DAILY, HOURLY) and push the details to statistics
		if timeDiff >= interval {
			sampleCount = float64(pointCount)
			average = sum / sampleCount
			statistics = append(statistics, &Statistics{
				Maximum:     maxValue,
				Minimum:     minValue,
				Sum:         sum,
				Average:     average,
				SampleCount: sampleCount,
				TimeStamp:   startTime,
			})
			maxValue, minValue = pointValues[pointCount].Point, pointValues[pointCount].Point
			pointCount, sum, diffCheckExecuted = 0, 0.0, true

			// Set the time interval as per granularity
			currentStartTime, _ := time.Parse(time.RFC3339, startTime)
			startTime = currentStartTime.Add(-time.Second * getIncrementalTimeAsPerGranularity(granularity)).Format(time.RFC3339)
		}

		if point.Point > maxValue {
			maxValue = point.Point
		}
		if point.Point < minValue {
			minValue = point.Point
		}

		sum += point.Point
		pointCount++
		pointIndex++
	}

	// Left over points which is not with in the same time interval
	if pointIndex == int64(len(pointValues)) && !diffCheckExecuted {
		sampleCount = float64(pointCount)
		average = sum / sampleCount
		statistics = append(statistics, &Statistics{
			Maximum:     maxValue,
			Minimum:     minValue,
			Sum:         sum,
			Average:     average,
			SampleCount: sampleCount,
			TimeStamp:   startTime,
		})
	}

	return statistics, nil
}

// Check time difference in second
func checkTimeDiff(startTime string, endTime string) float64 {
	dt1, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		return 0
	}
	dt2, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		return 0
	}

	return dt2.Sub(dt1).Seconds()
}
