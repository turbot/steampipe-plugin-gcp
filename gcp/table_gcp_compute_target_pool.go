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

func tableGcpComputeTargetPool(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_target_pool",
		Description: "GCP Compute Target Pool",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeTargetPool,
			Tags:       map[string]string{"service": "compute", "action": "targetPools.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeTargetPools,
			KeyColumns: plugin.KeyColumnSlice{
				// String columns
				{Name: "session_affinity", Require: plugin.Optional, Operators: []string{"<>", "="}},
			},
			Tags: map[string]string{"service": "compute", "action": "targetPools.list"},
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
				Type:        proto.ColumnType_INT,
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
				Transform:   transform.FromP(gcpComputeTargetPoolTurbotData, "Akas"),
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
				Transform:   transform.FromP(gcpComputeTargetPoolTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeTargetPools(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeTargetPools")
	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	filterString := ""
	if d.EqualsQuals["session_affinity"] != nil {
		filterString = "sessionAffinity=" + d.EqualsQuals["session_affinity"].GetStringValue()
	}

	// Max limit is set as per documentation
	// https://pkg.go.dev/google.golang.org/api@v0.48.0/compute/v1?utm_source=gopls#TargetPoolsAggregatedListCall.MaxResults
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

	resp := service.TargetPools.AggregatedList(project).Filter(filterString).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *compute.TargetPoolAggregatedList) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, item := range page.Items {
			for _, targetPool := range item.TargetPools {
				d.StreamListItem(ctx, targetPool)

				// Check if context has been cancelled or if the limit has been hit (if specified)
				// if there is a limit, it will return the number of rows required to reach this limit
				if d.RowsRemaining(ctx) == 0 {
					page.NextPageToken = ""
					return nil
				}
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

	var targetPool compute.TargetPool
	name := d.EqualsQuals["name"].GetStringValue()

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
	if len(targetPool.Name) < 1 {
		return nil, nil
	}

	return &targetPool, nil
}

//// TRANSFORM FUNCTION

func gcpComputeTargetPoolTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	targetPool := d.HydrateItem.(*compute.TargetPool)
	param := d.Param.(string)

	region := getLastPathElement(types.SafeString(targetPool.Region))
	project := strings.Split(targetPool.SelfLink, "/")[6]

	turbotData := map[string]interface{}{
		"Project": project,
		"Akas":    []string{"gcp://compute.googleapis.com/projects/" + project + "/regions/" + region + "/targetPools/" + targetPool.Name},
	}

	return turbotData[param], nil
}
