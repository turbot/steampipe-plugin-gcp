package gcp

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	"google.golang.org/api/bigquery/v2"
)

func tableGcpBigqueryTable(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_bigquery_table",
		Description: "GCP Bigquery Table",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"dataset_id", "id"}),
			Hydrate:    getBigqueryTable,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listBigQueryDatasets,
			Hydrate:       listBigqueryTables,
		},
		Columns: []*plugin.Column{
			// commonly used columns
			{
				Name:        "creation_time",
				Description: "The time when this table was created, in milliseconds since the epoch.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CreationTime").Transform(transform.UnixMsToTimestamp),
			},
			{
				Name:        "dataset_id",
				Description: "The ID of the dataset containing this table.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TableReference.DatasetId"),
			},
			{
				Name:        "expiration_time",
				Description: "The time when this table expires, in milliseconds since the epoch. If not present, the table will persist indefinitely. Expired tables will be deleted and their storage reclaimed. The defaultTableExpirationMs property of the encapsulating dataset can be used to set a default expirationTime on newly created tables.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ExpirationTime").Transform(transform.UnixMsToTimestamp),
			},
			{
				Name:        "friendly_name",
				Description: "A descriptive name for this table.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FriendlyName"),
			},
			{
				Name:        "id",
				Description: "An opaque ID uniquely identifying the table.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id"),
			},
			{
				Name:        "kind",
				Description: "The type of resource ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Kind"),
			},
			{
				Name:        "labels",
				Description: "The labels associated with this table. You can use these to organize and group your tables. Label keys and values can be no longer than 63 characters, can only contain lowercase letters, numeric characters, underscores and dashes. International characters are allowed. Label values are optional. Label keys must start with a letter and each label in the list must have a different key.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Labels"),
			},
			{
				Name:        "range_partitioning",
				Description: "If specified, configures range partitioning for this table.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("RangePartitioning"),
			},
			{
				Name:        "table_reference",
				Description: "The type of resource ID.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("TableReference"),
			},
			{
				Name:        "time_partitioning",
				Description: "If specified, configures time-based partitioning for this table.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("TimePartitioning"),
			},
			{
				Name:        "type",
				Description: "The type of table. Possible values are: TABLE, VIEW.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Type"),
			},
			{
				Name:        "self_link",
				Description: "A URL that can be used to access this resource again.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBigqueryTable,
			},
			{
				Name:        "view",
				Description: "dditional details for a view.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("View"),
			},
			// standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id"),
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
				Transform:   transform.From(gcpBigqueryTableAkas),
			},
			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBigqueryTable,
				Transform:   transform.FromField("Location"),
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
	// Get project details
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}
	project := projectData.Project
	resp := service.Tables.List(project, dataset.DatasetReference.DatasetId)
	if err := resp.Pages(ctx, func(page *bigquery.TableList) error {
		for _, table := range page.Tables {
			d.StreamListItem(ctx, table)
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
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}
	project := projectData.Project
	var datasetID, id string
	if h.Item != nil {
		data := h.Item.(*bigquery.TableListTables)
		datasetID = data.TableReference.DatasetId
		id = data.TableReference.TableId
	} else {
		datasetID = d.KeyColumnQuals["dataset_id"].GetStringValue()
		id = d.KeyColumnQuals["id"].GetStringValue()
	}
	resp, err := service.Tables.Get(project, datasetID, id).Do()
	if err != nil {
		return nil, err
	}
	return resp, err
}

//// TRANSFORM FUNCTIONS
func gcpBigqueryTableAkas(ctx context.Context, h *transform.TransformData) (interface{}, error) {
	data := TableID(h.HydrateItem)
	projectID := strings.Split(data, ":")[0]
	id := strings.Split(data, ":")[1]
	akas := []string{"gcp://bigquery.googleapis.com/projects/" + projectID + "/tables/" + id}
	return akas, nil
}

func TableID(item interface{}) string {
	switch item := item.(type) {
	case *bigquery.TableListTables:
		return item.Id
	case *bigquery.Table:
		return item.Id
	}
	return ""
}
