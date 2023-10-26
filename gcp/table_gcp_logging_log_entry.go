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
				{Name: "log_entry_operation_id", Require: plugin.Optional},
				{Name: "filter", Require: plugin.Optional, CacheMatch: "exact"},
			},
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
				Name:        "log_entry_operation_first",
				Description: "Set this to True if this is the first log entry in the operation.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Operation.First"),
			},
			{
				Name:        "log_entry_operation_last",
				Description: "Set this to True if this is the last log entry in the operation.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Operation.Last"),
			},
			{
				Name:        "filter",
				Type:        proto.ColumnType_STRING,
				Description: "The filter pattern for the search.",
				Transform:   transform.FromQual("filter"),
			},
			{
				Name:        "log_entry_operation_id",
				Description: "An arbitrary operation identifier. Log entries with the same identifier are assumed to be part of the same operation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Operation.Id"),
			},
			{
				Name:        "log_entry_operation_producer",
				Description: "An arbitrary producer identifier.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Operation.Producer"),
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
				Name:        "split_index",
				Description: "The index of this LogEntry in the sequence of split log entries.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Split.Index"),
			},
			{
				Name:        "total_splits",
				Description: "The total number of log entries that the original LogEntry was split into.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Split.TotalSplits"),
			},
			{
				Name:        "split_uid",
				Description: "A globally unique identifier for all log entries in a sequence of split log entries.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Split.Uid"),
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
				Name:        "resource_labels",
				Description: "Values for all of the labels listed in the associated monitored resource descriptor.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Resource.Labels"),
			},
			{
				Name:        "source_location",
				Description: "Source code location information associated with the log entry, if any.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InsertId"),
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
				Hydrate:     plugin.HydrateFunc(getProject).WithCache(),
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
	// Default limit is set as 1000
	pageSize := types.Int64(1000)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *pageSize {
			pageSize = limit
		}
	}

	// Get project details
	getProjectCached := plugin.HydrateFunc(getProject).WithCache()
	projectId, err := getProjectCached(ctx, d, h)
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
	getProjectCached := plugin.HydrateFunc(getProject).WithCache()
	projectId, err := getProjectCached(ctx, d, h)
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
		{"log_entry_operation_id", "operation.id", "string"},
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

// // TRANSFORM FUNCTION

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
