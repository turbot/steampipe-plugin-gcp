package gcp

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"google.golang.org/api/logging/v2"
)

//// TABLE DEFINITION
func tableGcpLoggingBucket(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_logging_bucket",
		Description: "GCP Logging Bucket",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getLoggingBucket,
		},
		List: &plugin.ListConfig{
			Hydrate: listLoggingBuckets,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The resource name of the bucket.",
				Type:        proto.ColumnType_STRING,
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
				Description: "Locked: Whether the bucket has been locked. The retention period on a locked bucket may not be changed. Locked buckets may only be deleted if they are empty.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "retention_days",
				Description: "Logs will be retained by default for this amount of time, after which they will automatically be deleted.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "update_time",
				Description: "he last update timestamp of the bucket.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromGo().NullIfZero(),
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
				Transform:   transform.From(bucketAka),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(bucketLocation),
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

//// LIST FUNCTION

func listLoggingBuckets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listLoggingBuckets")

	// Create service connection
	service, err := LoggingService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}

	project := projectData.Project

	resp := service.Projects.Locations.Buckets.List("projects/" + project + "/locations/-")
	if err := resp.Pages(
		ctx,
		func(page *logging.ListBucketsResponse) error {
			for _, bucket := range page.Buckets {
				d.StreamListItem(ctx, bucket)
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	return nil, err
}

//// HYDRATED FUNCTIONS

func getLoggingBucket(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getLoggingBucket")

	// Create Service Connection
	service, err := LoggingService(ctx, d)
	if err != nil {
		return nil, err
	}

	bucketName := d.KeyColumnQuals["name"].GetStringValue()

	// Match the bucket name pattern if we are doing query with where clause by passing name value
	matched, err := regexp.MatchString(`projects/.+/locations/.+/buckets/.+`, bucketName)

	if !matched {
		return nil, errors.New("Bucket name should match 'projects/[PROJECT_ID]/locations/[LOCATION_ID]/buckets/[BUCKET_ID]'")
	}

	op, err := service.Projects.Locations.Buckets.Get(bucketName).Do()
	if err != nil {
		return nil, err
	}

	return op, nil

}

//// TRANSFORM FUNCTIONS

func bucketLocation(_ context.Context, d *transform.TransformData) (interface{}, error) {
	bucket := d.HydrateItem.(*logging.LogBucket)
	location := strings.Split(bucket.Name, "/")[3]
	if location != "" {
		return location, nil
	}
	return "", nil
}

func bucketAka(_ context.Context, d *transform.TransformData) (interface{}, error) {
	bucketName := d.HydrateItem.(*logging.LogBucket).Name
	if bucketName != "" {
		return []string{"gcp://logging.googleapis.com/" + bucketName}, nil
	}
	return nil, nil
}
