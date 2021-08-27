package gcp

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"google.golang.org/api/compute/v1"
)

func tableGcpComputeZone(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_zone",
		Description: "GCP Compute Zone",
		List: &plugin.ListConfig{
			Hydrate:           listComputeZones,
			ShouldIgnoreError: isIgnorableError([]string{"403"}),
		},
		Columns: []*plugin.Column{
			// commonly used columns
			{
				Name:        "name",
				Description: "The name of the zone.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique identifier for the zone.",
				Type:        proto.ColumnType_STRING,
			},
			// region_name is a simpler view of region, without the full path
			{
				Name:        "region_name",
				Description: "Region name which hosts the zone.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Region").Transform(lastPathElement),
			},
			{
				Name:        "status",
				Description: "Status of the zone, either UP or DOWN.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "region",
				Description: "Full URL reference to the region which hosts the zone.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "Textual description of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "Type of the resource. Always compute#zone for zones.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the zone.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "available_cpu_platforms",
				Description: "Available cpu/platform selections for the zone.",
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
				Transform:   transform.FromP(gcpComputeZoneTurbotData, "Akas"),
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
				Transform:   transform.FromP(gcpComputeZoneTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeZones(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeZones")

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

	resp := service.Zones.List(project)
	if err := resp.Pages(
		ctx,
		func(page *compute.ZoneList) error {
			for _, zone := range page.Items {
				d.StreamListItem(ctx, zone)
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	return nil, nil
}

//// TRANSFORM FUNCTION

func gcpComputeZoneTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	zone := d.HydrateItem.(*compute.Zone)
	param := d.Param.(string)

	project := strings.Split(zone.SelfLink, "/")[6]
	region := getLastPathElement(zone.Region)

	turbotData := map[string]interface{}{
		"Project": project,
		"Akas":    []string{"gcp://compute.googleapis.com/projects/" + project + "/regions/" + region + "/zones/" + zone.Name},
	}

	return turbotData[param], nil
}
