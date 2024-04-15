package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/bigquery/v2"
)

//// TABLE DEFINITION

func tableGcpBigqueryTable(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_bigquery_table",
		Description: "GCP Bigquery Table",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"dataset_id", "table_id"}),
			Hydrate:    getBigqueryTable,
			Tags:       map[string]string{"service": "bigquery", "action": "tables.get"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listBigQueryDatasets,
			Hydrate:       listBigqueryTables,
			Tags:          map[string]string{"service": "bigquery", "action": "tables.list"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getBigqueryTable,
				Tags: map[string]string{"service": "bigquery", "action": "tables.get"},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "A descriptive name for this table, if one exists.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FriendlyName"),
			},
			{
				Name:        "table_id",
				Description: "The ID of the table resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TableReference.TableId"),
			},
			{
				Name:        "dataset_id",
				Description: "The ID of the dataset containing this table.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TableReference.DatasetId"),
			},
			{
				Name:        "id",
				Description: "An opaque ID uniquely identifying the table.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of table. Possible values are: TABLE, VIEW.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "A URL that can be used to access this resource again.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBigqueryTable,
			},
			{
				Name:        "creation_time",
				Description: "The time when this table was created, in milliseconds since the epoch.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CreationTime").Transform(transform.UnixMsToTimestamp),
			},
			{
				Name:        "description",
				Description: "A user-friendly description of this table.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBigqueryTable,
			},
			{
				Name:        "etag",
				Description: "A hash of the table metadata, used to ensure there were no concurrent modifications to the resource when attempting an update.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBigqueryTable,
			},
			{
				Name:        "expiration_time",
				Description: "The time when this table expires, in milliseconds since the epoch. If not present, the table will persist indefinitely. Expired tables will be deleted and their storage reclaimed. The defaultTableExpirationMs property of the encapsulating dataset can be used to set a default expirationTime on newly created tables.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ExpirationTime").Transform(transform.UnixMsToTimestamp),
			},
			{
				Name:        "kind",
				Description: "The type of resource ID.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kms_key_name",
				Description: "Describes the Cloud KMS encryption key that will be used to protect destination BigQuery table.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBigqueryTable,
				Transform:   transform.FromField("EncryptionConfiguration.KmsKeyName"),
			},
			{
				Name:        "last_modified_time",
				Description: "The time when this table was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("LastModifiedTime").Transform(transform.UnixMsToTimestamp),
				Hydrate:     getBigqueryTable,
			},
			{
				Name:        "num_bytes",
				Description: "The size of this table in bytes, excluding any data in the streaming buffer.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getBigqueryTable,
			},
			{
				Name:        "num_long_term_bytes",
				Description: "The number of bytes in the table that are considered 'long-term storage'.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getBigqueryTable,
			},
			{
				Name:        "num_physical_bytes",
				Description: "The physical size of this table in bytes, excluding any data in the streaming buffer.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getBigqueryTable,
			},
			{
				Name:        "num_rows",
				Description: "The number of rows of data in this table, excluding any data in the streaming buffer.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getBigqueryTable,
			},
			{
				Name:        "require_partition_filter",
				Description: "If set to true, queries over this table require a partition filter that can be used for partition elimination to be specified.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getBigqueryTable,
			},
			{
				Name:        "snapshot_time",
				Description: "The time at which the base table was snapshot.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getBigqueryTable,
				Transform:   transform.FromField("SnapshotDefinition.SnapshotTime"),
			},
			{
				Name:        "view_query",
				Description: "A query that BigQuery executes when the view is referenced.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("View.Query"),
				Hydrate:     getBigqueryTable,
			},
			{
				Name:        "view_use_legacy_sql",
				Description: "True if view is defined in legacy SQL dialect, false if in standard SQL.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("View.UseLegacySql"),
			},
			{
				Name:        "clustering_fields",
				Description: "One or more fields on which data should be clustered.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Clustering.Fields"),
			},
			{
				Name:        "external_data_configuration",
				Description: "Describes the data format, location, and other properties of a table stored outside of BigQuery.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBigqueryTable,
			},
			{
				Name:        "labels",
				Description: "The labels associated with this table. You can use these to organize and group your tables. Label keys and values can be no longer than 63 characters, can only contain lowercase letters, numeric characters, underscores and dashes. International characters are allowed. Label values are optional. Label keys must start with a letter and each label in the list must have a different key.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "materialized_view",
				Description: "Describes materialized view definition.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBigqueryTable,
			},
			{
				Name:        "range_partitioning",
				Description: "If specified, configures range partitioning for this table.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "schema_fields",
				Description: "Describes the fields in a table.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBigqueryTable,
				Transform:   transform.FromField("Schema.Fields"),
			},
			{
				Name:        "streaming_buffer",
				Description: "Contains information regarding this table's streaming buffer, if one is present.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBigqueryTable,
			},
			{
				Name:        "time_partitioning",
				Description: "If specified, configures time-based partitioning for this table.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "view",
				Description: "dditional details for a view.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(bigQueryTableTitle),
			},
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Labels"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(bigqueryTableAkas),
			},

			// GCP standard columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBigqueryTable,
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TableReference.ProjectId"),
			},
		},
	}
}

//// LIST FUNCTION

func listBigqueryTables(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("ListBigqueryTables")

	// Get a details of Cloud Dataset
	dataset := h.Item.(*bigquery.DatasetListDatasets)

	// Create Service Connection
	service, err := BigQueryService(ctx, d)
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

	resp := service.Tables.List(project, dataset.DatasetReference.DatasetId).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *bigquery.TableList) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, table := range page.Tables {
			d.StreamListItem(ctx, table)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				page.NextPageToken = ""
				return nil
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getBigqueryTable(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getBigqueryTable")

	// Create Service Connection
	service, err := BigQueryService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	var datasetID, id string
	if h.Item != nil {
		data := h.Item.(*bigquery.TableListTables)
		datasetID = data.TableReference.DatasetId
		id = data.TableReference.TableId
	} else {
		datasetID = d.EqualsQuals["dataset_id"].GetStringValue()
		id = d.EqualsQuals["table_id"].GetStringValue()
	}

	// Empty Check
	if id == "" || datasetID == "" {
		return nil, nil
	}

	resp, err := service.Tables.Get(project, datasetID, id).Do()
	if err != nil {
		return nil, err
	}
	return resp, err
}

//// TRANSFORM FUNCTIONS

func bigqueryTableAkas(_ context.Context, h *transform.TransformData) (interface{}, error) {
	data := tableID(h.HydrateItem)

	projectID := strings.Split(data, ":")[0]
	id := strings.Split(data, ":")[1]
	datasetId := strings.Split(id, ".")[0]
	id = strings.Split(id, ".")[1]

	akas := []string{"gcp://bigquery.googleapis.com/projects/" + projectID + "/datasets/" + datasetId + "/tables/" + id}
	return akas, nil
}

func bigQueryTableTitle(_ context.Context, h *transform.TransformData) (interface{}, error) {
	data := tableID(h.HydrateItem)
	name := tableName(h.HydrateItem)

	if len(name) > 0 {
		return name, nil
	}
	return strings.Split(strings.Split(data, ":")[1], ".")[1], nil
}

func tableID(item interface{}) string {
	switch item := item.(type) {
	case *bigquery.TableListTables:
		return item.Id
	case *bigquery.Table:
		return item.Id
	}
	return ""
}

func tableName(item interface{}) string {
	switch item := item.(type) {
	case *bigquery.TableListTables:
		return item.FriendlyName
	case *bigquery.Table:
		return item.FriendlyName
	}
	return ""
}
