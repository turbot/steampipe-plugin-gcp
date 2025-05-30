package gcp

import (
	"context"
	"encoding/json"
	"time"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"google.golang.org/api/logging/v2"
)

//// TABLE DEFINITION

func tableGcpLoggingLogEntry(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_logging_log_entry",
		Description: "GCP Logging Log Entry",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("insert_id"),
			Hydrate:    getGcpLoggingLogEntry,
			Tags:       map[string]string{"service": "logging", "action": "logEntries.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listGcpLoggingLogEntries,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "resource_type", Require: plugin.Optional},
				{Name: "severity", Require: plugin.Optional},
				{Name: "log_name", Require: plugin.Optional},
				{Name: "span_id", Require: plugin.Optional},
				{Name: "text_payload", Require: plugin.Optional},
				{Name: "receive_timestamp", Require: plugin.Optional, Operators: []string{"=", ">", "<", ">=", "<="}},
				{Name: "timestamp", Require: plugin.Optional, Operators: []string{"=", ">", "<", ">=", "<="}},
				{Name: "trace", Require: plugin.Optional},
				{Name: "operation_id", Require: plugin.Optional},
				{Name: "filter", Require: plugin.Optional, CacheMatch: "exact"},
			},
			Tags: map[string]string{"service": "logging", "action": "logEntries.list"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "log_name",
				Description: "The resource name of the log to which this log entry belongs to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "insert_id",
				Description: "A unique identifier for the log entry.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "filter",
				Type:        proto.ColumnType_STRING,
				Description: "The filter pattern for the search.",
				Transform:   transform.FromQual("filter"),
			},
			{
				Name:        "operation_id",
				Description: "An arbitrary operation identifier. Log entries with the same identifier are assumed to be part of the same operation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Operation.Id"),
			},
			{
				Name:        "receive_timestamp",
				Description: "The time the log entry was received by Logging.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "resource_type",
				Description: "The monitored resource type.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Resource.Type"),
			},
			{
				Name:        "severity",
				Description: "The severity of the log entry.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "span_id",
				Description: "The ID of the Cloud Trace (https://cloud.google.com/trace) span associated with the current operation in which the log is being written.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "text_payload",
				Description: "The log entry payload, represented as a Unicode string (UTF-8).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "timestamp",
				Description: "The time the event described by the log entry occurred.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "trace",
				Description: "The REST resource name of the trace being written to Cloud Trace (https://cloud.google.com/trace) in association with this log entry.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "trace_sampled",
				Description: "The sampling decision of the trace associated with the log entry.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "json_payload",
				Description: "The log entry payload, represented as a structure that is expressed as a JSON object.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(covertLogEntryByteArrayToJsonObject, "JsonPayload"),
			},
			{
				Name:        "proto_payload",
				Description: "The log entry payload, represented as a protocol buffer. Some Google Cloud Platform services use this field for their log entry payloads. The following protocol buffer types are supported; user-defined types are not supported: 'type.googleapis.com/google.cloud.audit.AuditLog' 'type.googleapis.com/google.appengine.logging.v1.RequestLog'",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(covertLogEntryByteArrayToJsonObject, "ProtoPayload"),
			},
			{
				Name:        "operation",
				Description: "Information about an operation associated with the log entry, if applicable.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "resource",
				Description: "The monitored resource that produced this log entry. Example: a log entry that reports a database error would be associated with the monitored resource designating the particular database that reported the error.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "split",
				Description: "Information indicating this LogEntry is part of a sequence of multiple log entries split from a single LogEntry.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "source_location",
				Description: "Source code location information associated with the log entry, if any.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "metadata",
				Description: "Auxiliary metadata for a MonitoredResource object. MonitoredResource objects contain the minimum set of information to uniquely identify a monitored resource instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "labels",
				Description: "A map of key, value pairs that provides additional information about the log entry. The labels can be user-defined or system-defined.User-defined labels are arbitrary key, value pairs that you can use to classify logs. System-defined labels are defined by GCP services for platform logs.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InsertId"),
			},
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Labels"),
			},

			// Standard GCP columns
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

func listGcpLoggingLogEntries(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := LoggingService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_logging_log_entry.listGcpLoggingLogEntries", "service_error", err)
		return nil, err
	}

	// Max limit isn't mentioned in the documentation
	// Default limit is set as 10000
	// 10000 seems to be a balanced limit, based on initial tests for retrieving 140k log entries: 5000 (124s), 10000 (88s), 20000 (84s).
	pageSize := types.Int64(10000)
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

	param := &logging.ListLogEntriesRequest{
		PageSize:   *pageSize,
		ProjectIds: []string{project},
	}

	filter := ""

	if d.EqualsQualString("filter") != "" {
		filter = d.EqualsQualString("filter")
	} else {
		filter = buildLoggingLogEntryFilterParam(d.Quals)
	}

	if filter != "" {
		param.Filter = filter
	}

	op := service.Entries.List(param)

	if err := op.Pages(
		ctx,
		func(page *logging.ListLogEntriesResponse) error {
			// apply rate limiting
			d.WaitForListRateLimit(ctx)

			for _, entry := range page.Entries {
				d.StreamListItem(ctx, entry)

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
		plugin.Logger(ctx).Error("gcp_logging_log_entry.listGcpLoggingLogEntries", "api_error", err)
		return nil, err
	}

	return nil, err
}

//// HYDRATE FUNCTION

func getGcpLoggingLogEntry(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := LoggingService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_logging_log_entry.getGcpLoggingLogEntry", "service_error", err)
		return nil, err
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	param := &logging.ListLogEntriesRequest{
		ProjectIds: []string{project},
	}

	insertId := d.EqualsQualString("insert_id")
	filter := ""

	if insertId != "" {
		filter = "insertId" + " = \"" + insertId + "\""
	}
	param.Filter = filter

	op, err := service.Entries.List(param).Do()

	if err != nil {
		plugin.Logger(ctx).Error("gcp_logging_log_entry.getGcpLoggingLogEntry", "api_error", err)
		return nil, err
	}

	if len(op.Entries) > 0 {
		return op.Entries[0], nil
	}

	return nil, nil
}

//// UTILITY FUNCTION

func buildLoggingLogEntryFilterParam(equalQuals plugin.KeyColumnQualMap) string {
	filter := ""

	filterQuals := []filterQualMap{
		{"resource_type", "resource.type", "string"},
		{"severity", "severity", "string"},
		{"log_name", "logName", "string"},
		{"span_id", "spanId", "string"},
		{"text_payload", "textPayload", "string"},
		{"trace", "trace", "string"},
		{"operation_id", "operation.id", "string"},
		{"receive_timestamp", "receiveTimestamp", "timestamp"},
		{"timestamp", "timestamp", "timestamp"},
	}

	for _, filterQualItem := range filterQuals {
		filterQual := equalQuals[filterQualItem.ColumnName]
		if filterQual == nil {
			continue
		}

		// Check only if filter qual map matches with optional column name
		if filterQual.Name == filterQualItem.ColumnName {
			if filterQual.Quals == nil {
				continue
			}
		}

		for _, qual := range filterQual.Quals {
			if qual.Value != nil {
				value := qual.Value
				switch filterQualItem.Type {
				case "string":
					if filter == "" {
						filter = filterQualItem.PropertyPath + " = \"" + value.GetStringValue() + "\""
					} else {
						filter = filter + " AND " + filterQualItem.PropertyPath + " = \"" + value.GetStringValue() + "\""
					}
				case "timestamp":
					propertyPath := filterQualItem.PropertyPath
					if filter == "" {
						switch qual.Operator {
						case "=":
							filter = propertyPath + " = \"" + value.GetTimestampValue().AsTime().Format(time.RFC3339) + "\""
						case ">":
							filter = propertyPath + " > \"" + value.GetTimestampValue().AsTime().Format(time.RFC3339) + "\""
						case "<":
							filter = propertyPath + " < \"" + value.GetTimestampValue().AsTime().Format(time.RFC3339) + "\""
						case ">=":
							filter = propertyPath + " >= \"" + value.GetTimestampValue().AsTime().Format(time.RFC3339) + "\""
						case "<=":
							filter = propertyPath + " <= \"" + value.GetTimestampValue().AsTime().Format(time.RFC3339) + "\""
						}
					} else {
						switch qual.Operator {
						case "=":
							filter = filter + " AND " + propertyPath + " = \"" + value.GetTimestampValue().AsTime().Format(time.RFC3339) + "\""
						case ">":
							filter = filter + " AND " + propertyPath + " > \"" + value.GetTimestampValue().AsTime().Format(time.RFC3339) + "\""
						case "<":
							filter = filter + " AND " + propertyPath + " < \"" + value.GetTimestampValue().AsTime().Format(time.RFC3339) + "\""
						case ">=":
							filter = filter + " AND " + propertyPath + " >= \"" + value.GetTimestampValue().AsTime().Format(time.RFC3339) + "\""
						case "<=":
							filter = filter + " AND " + propertyPath + " <= \"" + value.GetTimestampValue().AsTime().Format(time.RFC3339) + "\""
						}
					}
				}
			}
		}
	}
	return filter
}

//// TRANSFORM FUNCTION

func covertLogEntryByteArrayToJsonObject(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	entry := d.HydrateItem.(*logging.LogEntry)
	param := d.Param.(string)

	var protoPlayload interface{}
	var jsonPayload interface{}

	a, err := entry.ProtoPayload.MarshalJSON()
	if err != nil {
		return nil, err
	}
	b, err := entry.JsonPayload.MarshalJSON()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(a, &protoPlayload)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_logging_log_entry.covertLogEntryByteArrayToJsonObject.protoPlayload", err)
	}

	err = json.Unmarshal(b, &jsonPayload)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_logging_log_entry.covertLogEntryByteArrayToJsonObject.jsonPayload", err)
	}

	payload := map[string]interface{}{
		"JsonPayload":  jsonPayload,
		"ProtoPayload": protoPlayload,
	}

	return payload[param], nil
}
