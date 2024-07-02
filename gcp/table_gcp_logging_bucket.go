package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"google.golang.org/api/logging/v2"
)

//// TABLE DEFINITION

func tableGcpLoggingBucket(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_logging_bucket",
		Description: "GCP Logging Bucket",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "location"}),
			Hydrate:    getLoggingBucket,
			Tags:       map[string]string{"service": "logging", "action": "buckets.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listLoggingBuckets,
			Tags:    map[string]string{"service": "logging", "action": "buckets.list"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The resource name of the bucket.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(lastPathElement),
			},
			{
				Name:        "self_link",
				Description: "The server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(loggingBucketTurbotData, "SelfLink"),
			},
			{
				Name:        "create_time",
				Description: "The creation timestamp of the bucket. This is not set for any of the default buckets.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromGo().NullIfZero(),
			},
			{
				Name:        "description",
				Description: "Describes this bucket.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The bucket lifecycle state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "locked",
				Description: "Specifies whether the bucket has been locked, or not. The retention period on a locked bucket may not be changed. Locked buckets may only be deleted if they are empty.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "retention_days",
				Description: "Logs will be retained by default for this amount of time, after which they will automatically be deleted.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "update_time",
				Description: "The last update timestamp of the bucket.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromGo().NullIfZero(),
			},

			// GCP standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(lastPathElement),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(loggingBucketTurbotData, "Akas"),
			},

			// GCP standard columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(loggingBucketTurbotData, "Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(loggingBucketTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listLoggingBuckets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listLoggingBuckets")

	// Create service connection
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

	// '-' for all locations...
	resp := service.Projects.Locations.Buckets.List("projects/" + project + "/locations/-").PageSize(*pageSize)
	if err := resp.Pages(
		ctx,
		func(page *logging.ListBucketsResponse) error {
			// apply rate limiting
			d.WaitForListRateLimit(ctx)

			for _, bucket := range page.Buckets {
				d.StreamListItem(ctx, bucket)

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

func getLoggingBucket(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getLoggingBucket")

	// Create Service Connection
	service, err := LoggingService(ctx, d)
	if err != nil {
		return nil, err
	}

	bucketName := d.EqualsQuals["name"].GetStringValue()
	locationId := d.EqualsQuals["location"].GetStringValue()

	// Return nil, if no input provided
	if bucketName == "" || locationId == "" {
		return nil, nil
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	bucketNameWithLocation := "projects/" + project + "/locations/" + locationId + "/buckets/" + bucketName

	op, err := service.Projects.Locations.Buckets.Get(bucketNameWithLocation).Do()
	if err != nil {
		return nil, err
	}

	return op, nil

}

//// TRANSFORM FUNCTIONS

func loggingBucketTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*logging.LogBucket)
	param := d.Param.(string)

	// Fetch data from name
	splittedTitle := strings.Split(data.Name, "/")

	turbotData := map[string]interface{}{
		"Project":  splittedTitle[1],
		"Location": splittedTitle[3],
		"SelfLink": "https://logging.googleapis.com/v2/" + data.Name,
		"Akas":     []string{"gcp://logging.googleapis.com/" + data.Name},
	}

	return turbotData[param], nil
}
