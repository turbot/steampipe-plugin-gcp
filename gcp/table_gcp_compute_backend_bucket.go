package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"google.golang.org/api/compute/v1"
)

//// TABLE DEFINITION

func tableGcpComputeBackendBucket(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_backend_bucket",
		Description: "GCP Compute Backend Bucket",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeBackendBucket,
			Tags:       map[string]string{"service": "compute", "action": "backendBuckets.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeBackendBuckets,
			KeyColumns: plugin.KeyColumnSlice{
				// String columns
				{Name: "bucket_name", Require: plugin.Optional, Operators: []string{"<>", "="}},

				// Boolean columns
				{Name: "enable_cdn", Require: plugin.Optional, Operators: []string{"<>", "="}},
			},
			Tags: map[string]string{"service": "compute", "action": "backendBuckets.list"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "A friendly name that identifies the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique identifier for the resource.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "description",
				Description: "A user-specified, human-readable description of the backend bucket.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "enable_cdn",
				Description: "Specifies whether the Cloud CDN is enabled for this backend bucket, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "bucket_name",
				Description: "Specifies the name of the cloud storage bucket.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_timestamp",
				Description: "The creation timestamp of the resource.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "signed_url_cache_max_age_sec",
				Description: "Specifies the maximum number of seconds the response to a signed URL request will be considered fresh.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("CdnPolicy.SignedUrlCacheMaxAgeSec"),
			},
			{
				Name:        "self_link",
				Description: "The server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "signed_url_key_names",
				Description: "A list od names of the keys for signing request URLs.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CdnPolicy.SignedUrlKeyNames"),
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
				Transform:   transform.FromP(backendBucketSelfLinkToTurbotData, "Akas"),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Region").Transform(lastPathElement),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(backendBucketSelfLinkToTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeBackendBuckets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeBackendBuckets")

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	filterQuals := []filterQualMap{
		{"bucket_name", "bucketName", "string"},
		{"enable_cdn", "enableCdn", "boolean"},
	}

	filters := buildQueryFilterFromQuals(filterQuals, d.Quals)
	filterString := ""
	if len(filters) > 0 {
		filterString = strings.Join(filters, " ")
	}

	// Max limit is set as per documentation
	// https://pkg.go.dev/google.golang.org/api@v0.48.0/compute/v1#BackendBucketsListCall.MaxResults
	pageSize := types.Int64(500)
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

	resp := service.BackendBuckets.List(project).Filter(filterString).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *compute.BackendBucketList) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, backendBucket := range page.Items {
			d.StreamListItem(ctx, backendBucket)

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

func getComputeBackendBucket(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := ComputeService(ctx, d)
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

	// Error: pq: rpc error: code = Unknown desc = json: invalid use of ,string struct tag,
	// trying to unmarshal "projects/project/global/backendBuckets/" into uint64
	if len(name) < 1 {
		return nil, nil
	}

	req, err := service.BackendBuckets.Get(project, name).Do()
	if err != nil {
		return nil, err
	}

	return req, nil
}

//// TRANSFORM FUNCTIONS

func backendBucketSelfLinkToTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	backendBucket := d.HydrateItem.(*compute.BackendBucket)
	param := d.Param.(string)

	project := strings.Split(backendBucket.SelfLink, "/")[6]

	turbotData := map[string]interface{}{
		"Project": project,
		"Akas":    []string{"gcp://compute.googleapis.com/projects/" + project + "/global/backendBuckets/" + backendBucket.Name},
	}

	return turbotData[param], nil
}
