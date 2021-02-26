package gcp

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"google.golang.org/api/compute/v1"
)

func tableGcpComputeRegion(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_region",
		Description: "GCP Compute Region",
		List: &plugin.ListConfig{
			Hydrate: listComputeRegions,
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
				Description: "Status of the region, either UP or DOWN. Possible values: \"DOWN\" and \"UP\".",
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

func listComputeRegions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeRegions")

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

	resp := service.Regions.List(project)
	if err := resp.Pages(
		ctx,
		func(page *compute.RegionList) error {
			for _, region := range page.Items {
				d.StreamListItem(ctx, region)
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
