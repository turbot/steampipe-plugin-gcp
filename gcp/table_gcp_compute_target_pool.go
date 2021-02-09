package gcp

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	"google.golang.org/api/compute/v1"
)

func tableGcpComputeTargetPool(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_target_pool",
		Description: "GCP Compute Traget Pool",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeTargetPool,
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeTargetPools,
		},
		Columns: []*plugin.Column{

			{
				Name:        "name",
				Description: "Name of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique identifier for the resource.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "description",
				Description: "An optional description of this resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_timestamp",
				Description: "Timestamp when the target pool was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "kind",
				Description: "Type of the resource. Always compute#targetPool for target pools.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "region",
				Description: "URL of the region where the target pool resides.",
				Type:        proto.ColumnType_STRING,
			},
			// region_name is a simpler view of the zone, without the full path
			{
				Name:        "region_name",
				Description: "The region name where the target pool resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Region").Transform(lastPathElement),
			},
			{
				Name:        "session_affinity",
				Description: "Session affinity option, must be one of the following values: (CLIENT_IP | CLIENT_IP_PORT_PROTO | CLIENT_IP_PROTO | GENERATED_COOKIE | HEADER_FIELD | HTTP_COOKIE | NONE )",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "health_checks",
				Description: "The URL of the HttpHealthCheck resource. A member instance in this pool is considered healthy if and only if the health checks pass. An empty list means all member instances will be considered healthy at all times.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "instances",
				Description: "A list of resource URLs to the virtual machine instances serving this pool. They must live in zones contained in the same region as this pool.",
				Type:        proto.ColumnType_JSON,
			},
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
				Transform:   transform.From(targetPoolAka),
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
				Transform:   transform.FromConstant(activeProject()),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeTargetPools(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeTargetPools")
	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	project := activeProject()
	resp := service.TargetPools.AggregatedList(project)
	if err := resp.Pages(ctx, func(page *compute.TargetPoolAggregatedList) error {
		for _, item := range page.Items {
			for _, targetPool := range item.TargetPools {
				d.StreamListItem(ctx, targetPool)
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTION

func getComputeTargetPool(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	var targetPool compute.TargetPool
	name := d.KeyColumnQuals["name"].GetStringValue()
	project := activeProject()

	resp := service.TargetPools.AggregatedList(project).Filter("name=" + name)
	if err := resp.Pages(
		ctx,
		func(page *compute.TargetPoolAggregatedList) error {
			for _, item := range page.Items {
				for _, i := range item.TargetPools {
					targetPool = *i
				}
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	// If the specified resource is not present, API does not return any not found errors
	// it returns the data as {<nil> []   0 []   []    {0 map[]} [] []}
	if len(targetPool.Name) < 1 {
		return nil, nil
	}

	return &targetPool, nil
}

//// TRANSFORM FUNCTION

func targetPoolAka(_ context.Context, d *transform.TransformData) (interface{}, error) {
	targetPool := d.HydrateItem.(*compute.TargetPool)
	regionName := getLastPathElement(types.SafeString(targetPool.Region))

	akas := []string{"gcp://compute.googleapis.com/projects/" + activeProject() + "/regions/" + regionName + "/targetPools/" + targetPool.Name}

	return akas, nil
}
