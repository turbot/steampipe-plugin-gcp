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

func tableGcpComputeRegion(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_region",
		Description: "GCP Compute Region",
		List: &plugin.ListConfig{
			Hydrate: listComputeRegions,
			KeyColumns: plugin.KeyColumnSlice{
				// String columns
				{Name: "name", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "status", Require: plugin.Optional, Operators: []string{"<>", "="}},
			},
			Tags: map[string]string{"service": "compute", "action": "regions.list"},
		},
		Columns: []*plugin.Column{
			// commonly used columns
			{
				Name:        "name",
				Description: "The name of the region.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique identifier for the region.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "Status of the region, either UP or DOWN.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "Textual description of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "Type of the resource. Always compute#region for regions.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the region.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "quotas",
				Description: "Quotas assigned to this region.",
				Type:        proto.ColumnType_JSON,
			},
			// zone_names is a simpler view of zones, without the full path
			{
				Name:        "zone_names",
				Description: "A list of zones available in this region, in the form of zone_id.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(regionZoneNames),
			},
			{
				Name:        "zones",
				Description: "A list of zones available in this region, in the form of resource URLs.",
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
				Transform:   transform.FromP(gcpComputeRegionTurbotData, "Akas"),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpComputeRegionTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeRegions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeRegions")

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	filterQuals := []filterQualMap{
		{"name", "name", "string"},
		{"status", "status", "string"},
	}

	filters := buildQueryFilterFromQuals(filterQuals, d.Quals)
	filterString := ""
	if len(filters) > 0 {
		filterString = strings.Join(filters, " ")
	}

	// Max limit is set as per documentation
	// https://pkg.go.dev/google.golang.org/api@v0.48.0/compute/v1?utm_source=gopls#RegionsListCall.MaxResults
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

	resp := service.Regions.List(project).Filter(filterString).MaxResults(*pageSize)
	if err := resp.Pages(
		ctx,
		func(page *compute.RegionList) error {
			// apply rate limiting
			d.WaitForListRateLimit(ctx)

			for _, region := range page.Items {
				d.StreamListItem(ctx, region)

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

	return nil, nil
}

//// TRANSFORM FUNCTION

func regionZoneNames(_ context.Context, d *transform.TransformData) (interface{}, error) {
	region := d.HydrateItem.(*compute.Region)

	zoneNames := []string{}
	for _, zoneURL := range region.Zones {
		zoneNames = append(zoneNames, getLastPathElement(zoneURL))
	}

	return zoneNames, nil
}

func gcpComputeRegionTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	region := d.HydrateItem.(*compute.Region)
	param := d.Param.(string)
	project := strings.Split(region.SelfLink, "/")[6]

	turbotData := map[string]interface{}{
		"Project": project,
		"Akas":    []string{"gcp://compute.googleapis.com/projects/" + project + "/regions/" + region.Name},
	}

	return turbotData[param], nil
}
