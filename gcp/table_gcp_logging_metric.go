package gcp

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"google.golang.org/api/logging/v2"
)

//// TABLE DEFINITION

func tableGcpLoggingMetric(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_logging_metric",
		Description: "GCP Logging Metric",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getGcpLoggingMetric,
			Tags:       map[string]string{"service": "logging", "action": "logMetrics.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listGcpLoggingMetrics,
			Tags:    map[string]string{"service": "logging", "action": "logMetrics.list"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The client-assigned metric identifier.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A user-specified, human-readable description of the metric.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version",
				Description: "The API version that created or updated this metric.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "filter",
				Description: "An advanced logs filter used to match log entries.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "The creation timestamp of the metric.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "exponential_buckets_options_growth_factor",
				Description: "Specifies the growth factor of a bucket with exponential sequence, used to create a histogram for the distribution. Value must be greater than 1.",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("BucketOptions.ExponentialBuckets.GrowthFactor"),
			},
			{
				Name:        "exponential_buckets_options_num_finite_buckets",
				Description: "Specifies the finite buckets, used to create a histogram for the distribution. Value must be greater than 0.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("BucketOptions.ExponentialBuckets.NumFiniteBuckets"),
			},
			{
				Name:        "exponential_buckets_options_scale",
				Description: "Specifies the scale. Value must be greater than 0.",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("BucketOptions.ExponentialBuckets.Scale"),
			},
			{
				Name:        "linear_buckets_options_num_finite_buckets",
				Description: "Specifies the number of finite buckets of linear sequence. Value must be greater than 0.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("BucketOptions.LinearBuckets.NumFiniteBuckets"),
			},
			{
				Name:        "linear_buckets_options_offset",
				Description: "Specifies the lower bound of the first bucket.",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("BucketOptions.LinearBuckets.Offset"),
			},
			{
				Name:        "linear_buckets_options_width",
				Description: "Specifies the width, used to create a histogram for the distribution.",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("BucketOptions.LinearBuckets.Width"),
			},
			{
				Name:        "metric_descriptor_display_name",
				Description: "A concise name for the metric, which can be displayed in user interfaces.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MetricDescriptor.DisplayName"),
			},
			{
				Name:        "metric_descriptor_metric_kind",
				Description: "The kind of the metric. Possible values are 'METRIC_KIND_UNSPECIFIED', 'GAUGE', 'DELTA', 'CUMULATIVE'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MetricDescriptor.MetricKind"),
			},
			{
				Name:        "metric_descriptor_type",
				Description: "The metric type, including its DNS name in prefix.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MetricDescriptor.Type"),
			},
			{
				Name:        "metric_descriptor_unit",
				Description: "The units in which the metric value is reported.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MetricDescriptor.Unit"),
			},
			{
				Name:        "metric_descriptor_value_type",
				Description: "The value type of the measurement.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MetricDescriptor.ValueType"),
			},
			{
				Name:        "update_time",
				Description: "The last update timestamp of the metric.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "value_extractor",
				Description: "A value_extractor is required when using a distribution logs-based metric to extract the values to record from a log entry.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "explicit_buckets_options_bounds",
				Description: "Specifies the bounds, used to create a histogram for the distribution.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("BucketOptions.ExplicitBuckets.Bounds"),
			},
			{
				Name:        "label_extractors",
				Description: "A map from a label key string to an extractor expression which is used to extract data from a log entry field and assign as the label value.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "metric_descriptor_labels",
				Description: "The set of labels that can be used to instance of this metric type.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("MetricDescriptor.Labels"),
			},

			// standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     metricNameToAkas,
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

//// FETCH FUNCTIONS

func listGcpLoggingMetrics(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := LoggingService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Max limit isn't mentioned in the documentation
	// Default limit is set as 1000
	pageSize := types.Int64(1000)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *pageSize {
			pageSize = limit
		}
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	resp := service.Projects.Metrics.List("projects/" + project).PageSize(*pageSize)
	if err := resp.Pages(
		ctx,
		func(page *logging.ListLogMetricsResponse) error {
			// apply rate limiting
			d.WaitForListRateLimit(ctx)

			for _, metric := range page.Metrics {
				d.StreamListItem(ctx, metric)

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
		return nil, err
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getGcpLoggingMetric(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getGcpLoggingMetric")

	// Create Service Connection
	service, err := LoggingService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)
	name := d.EqualsQuals["name"].GetStringValue()

	op, err := service.Projects.Metrics.Get("projects/" + project + "/metrics/" + name).Do()
	if err != nil {
		plugin.Logger(ctx).Debug("getGcpLoggingMetric__", "ERROR", err)
		return nil, err
	}

	// If the name has been passed as empty string, API does not returns any error
	if len(op.Name) < 1 {
		return nil, nil
	}

	return op, nil
}

func metricNameToAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	metric := h.Item.(*logging.LogMetric)

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	akas := []string{"gcp://logging.googleapis.com/projects/" + project + "/metrics/" + metric.Name}
	return akas, nil
}
