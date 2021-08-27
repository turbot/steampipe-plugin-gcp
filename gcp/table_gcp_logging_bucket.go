package gcp

import (
	"context"
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
			KeyColumns: plugin.AllColumns([]string{"name", "location"}),
			Hydrate:    getLoggingBucket,
		},
		List: &plugin.ListConfig{
			Hydrate:           listLoggingBuckets,
			ShouldIgnoreError: isIgnorableError([]string{"403"}),
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

	// '-' for all locations...
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

//// HYDRATE FUNCTIONS

func getLoggingBucket(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getLoggingBucket")

	// Create Service Connection
	service, err := LoggingService(ctx, d)
	if err != nil {
		return nil, err
	}

	bucketName := d.KeyColumnQuals["name"].GetStringValue()
	locationId := d.KeyColumnQuals["location"].GetStringValue()

	// Return nil, if no input provided
	if bucketName == "" || locationId == "" {
		return nil, nil
	}

	projectInfo, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}
	bucketNameWithLocation := "projects/" + projectInfo.Project + "/locations/" + locationId + "/buckets/" + bucketName

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
