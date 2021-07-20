package gcp

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

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
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeBackendBuckets,
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

func listComputeBackendBuckets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeBackendBuckets")

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}
	project := projectData.Project

	resp := service.BackendBuckets.List(project)
	if err := resp.Pages(ctx, func(page *compute.BackendBucketList) error {
		for _, backendBucket := range page.Items {
			d.StreamListItem(ctx, backendBucket)
		}
		return nil
	}); err != nil {
		if IsForbiddenError(err) {
			return nil, nil
		}
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
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}
	project := projectData.Project

	name := d.KeyColumnQuals["name"].GetStringValue()

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
