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

func tableGcpBigQueryDataset(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_bigquery_dataset",
		Description: "GCP BigQuery Dataset",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("dataset_id"),
			Hydrate:    getBigQueryDataset,
			Tags:       map[string]string{"service": "bigquery", "action": "datasets.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listBigQueryDatasets,
			Tags:    map[string]string{"service": "bigquery", "action": "datasets.list"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getBigQueryDataset,
				Tags: map[string]string{"service": "bigquery", "action": "datasets.get"},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "A descriptive name for the dataset, if one exists.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FriendlyName"),
			},
			{
				Name:        "dataset_id",
				Description: "The ID of the dataset resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatasetReference.DatasetId"),
			},
			{
				Name:        "id",
				Description: "The fully-qualified, unique, opaque ID of the dataset.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "The type of the resource. This property always returns the value 'bigquery#dataset'.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The time when this dataset was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getBigQueryDataset,
				Transform:   transform.FromField("CreationTime").Transform(transform.UnixMsToTimestamp),
			},
			{
				Name:        "description",
				Description: "A user-friendly description of the dataset.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBigQueryDataset,
			},
			{
				Name:        "etag",
				Description: "A hash of the resource.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBigQueryDataset,
			},
			{
				Name:        "default_partition_expiration_ms",
				Description: "The default partition expiration for all partitioned tables in the dataset, in milliseconds.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getBigQueryDataset,
			},
			{
				Name:        "default_table_expiration_ms",
				Description: "The default lifetime of all tables in the dataset, in milliseconds.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getBigQueryDataset,
			},
			{
				Name:        "kms_key_name",
				Description: "Describes the Cloud KMS encryption key that will be used to protect destination BigQuery table.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBigQueryDataset,
				Transform:   transform.FromField("DefaultEncryptionConfiguration.KmsKeyName"),
			},
			{
				Name:        "last_modified_time",
				Description: "The date when this dataset or any of its tables was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getBigQueryDataset,
				Transform:   transform.FromField("LastModifiedTime").Transform(transform.UnixMsToTimestamp),
			},
			{
				Name:        "self_link",
				Description: "An URL that can be used to access the resource again.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBigQueryDataset,
			},
			{
				Name:        "access",
				Description: "An array of objects that define dataset access for one or more entities.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBigQueryDataset,
			},
			{
				Name:        "labels",
				Description: "A set of labels associated with this dataset.",
				Type:        proto.ColumnType_JSON,
			},

			// standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(bigQueryDatasetTitle),
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
				Transform:   transform.From(bigQueryDatasetAka),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatasetReference.ProjectId"),
			},
		},
	}
}

//// LIST FUNCTION

func listBigQueryDatasets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listBigQueryDatasets")

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

	resp := service.Datasets.List(project).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *bigquery.DatasetList) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, dataset := range page.Datasets {
			d.StreamListItem(ctx, dataset)

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

func getBigQueryDataset(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

	var id string
	if h.Item != nil {
		data := datasetID(h.Item)
		id = strings.Split(data, ":")[1]
	} else {
		id = d.EqualsQuals["dataset_id"].GetStringValue()
	}

	// check if id is empty
	if id == "" {
		return nil, nil
	}

	resp, err := service.Datasets.Get(project, id).Do()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

//// TRANSFORM FUNCTIONS

func bigQueryDatasetAka(ctx context.Context, h *transform.TransformData) (interface{}, error) {
	data := datasetID(h.HydrateItem)

	projectID := strings.Split(data, ":")[0]
	id := strings.Split(data, ":")[1]

	akas := []string{"gcp://bigquery.googleapis.com/projects/" + projectID + "/datasets/" + id}

	return akas, nil
}

func bigQueryDatasetTitle(ctx context.Context, h *transform.TransformData) (interface{}, error) {
	data := datasetID(h.HydrateItem)
	name := datasetName(h.HydrateItem)

	if len(name) > 0 {
		return name, nil
	}
	return strings.Split(data, ":")[1], nil
}

func datasetID(item interface{}) string {
	switch item := item.(type) {
	case *bigquery.DatasetListDatasets:
		return item.Id
	case *bigquery.Dataset:
		return item.Id
	}
	return ""
}

func datasetName(item interface{}) string {
	switch item := item.(type) {
	case *bigquery.DatasetListDatasets:
		return item.FriendlyName
	case *bigquery.Dataset:
		return item.FriendlyName
	}
	return ""
}
