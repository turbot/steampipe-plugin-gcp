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

func tableGcpComputeZone(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_zone",
		Description: "GCP Compute Zone",
		List: &plugin.ListConfig{
			Hydrate: listComputeZones,
			KeyColumns: plugin.KeyColumnSlice{
				// String columns
				{Name: "name", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "status", Require: plugin.Optional, Operators: []string{"<>", "="}},
			},
			Tags: map[string]string{"service": "compute", "action": "zones.list"},
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

func listComputeZones(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeZones")

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
	// https://pkg.go.dev/google.golang.org/api@v0.48.0/compute/v1?utm_source=gopls#ZonesListCall.MaxResults
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

	resp := service.Zones.List(project).Filter(filterString).MaxResults(*pageSize)
	if err := resp.Pages(
		ctx,
		func(page *compute.ZoneList) error {
			// apply rate limiting
			d.WaitForListRateLimit(ctx)

			for _, zone := range page.Items {
				d.StreamListItem(ctx, zone)

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
