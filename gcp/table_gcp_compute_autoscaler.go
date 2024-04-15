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

func tableGcpComputeAutoscaler(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_autoscaler",
		Description: "GCP Compute Autoscaler",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeAutoscaler,
			Tags:       map[string]string{"service": "compute", "action": "autoscalers.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeAutoscaler,
			Tags:    map[string]string{"service": "compute", "action": "autoscalers.list"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "Name of the Autoscaler.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique identifier for this Autoscaler. This identifier is defined by the server.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "creation_timestamp",
				Description: "Timestamp when the Autoscaler was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "An optional description of this resource. Provide this property when you create the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "Type of the resource. Always compute#autoscaler for Autoscalers.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "recommended_size",
				Description: "Target recommended MIG size (number of instances) computed by autoscaler.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "region",
				Description: "The URL of the region where the Autoscaler resides. Only applicable for regional resources.",
				Type:        proto.ColumnType_STRING,
			},

			// region_name is a simpler view of the region, without the full path
			{
				Name:        "region_name",
				Description: "Name of the region where the Autoscaler resides. Only applicable for regional resources.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Region").Transform(lastPathElement),
			},
			{
				Name:        "self_link",
				Description: "Server-defined fully-qualified URL for this resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the autoscaler configuration.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "target",
				Description: "URL of the managed instance group that this autoscaler will scale.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "zone",
				Description: "URL of the zone where the Autoscaler resides.",
				Type:        proto.ColumnType_STRING,
			},

			// zone_name is a simpler view of the zone, without the full path
			{
				Name:        "zone_name",
				Description: "The zone name in which the Autoscaler resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Zone").Transform(lastPathElement),
			},

			// JSON columns
			{
				Name:        "autoscaling_policy",
				Description: "The configuration parameters for the autoscaling algorithm.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "scaling_schedule_status",
				Description: "Status information of existing scaling schedules.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "status_details",
				Description: "Human-readable details about the current state of the autoscaler.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
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
				Transform:   transform.From(autoscalerAka),
			},

			// GCP standard columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(autoscalerLocation, "Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(autoscalerLocation, "Project"),
			},
		},
	}
}

//// LIST FUNCTIONS

func listComputeAutoscaler(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Max limit is set as per documentation
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

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_compute_autoscaler.listComputeAutoscaler", "service_creation_err", err)
		return nil, err
	}

	resp := service.Autoscalers.AggregatedList(project).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *compute.AutoscalerAggregatedList) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, item := range page.Items {
			for _, autoscaler := range item.Autoscalers {
				d.StreamListItem(ctx, autoscaler)

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
		plugin.Logger(ctx).Error("gcp_compute_autoscaler.listComputeAutoscaler", "api_err", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getComputeAutoscaler(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	var autoscaler compute.Autoscaler
	name := d.EqualsQuals["name"].GetStringValue()
	if name == "" {
		return nil, nil
	}

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_compute_autoscaler.getComputeAutoscaler", "service_creation_err", err)
		return nil, err
	}

	resp := service.Autoscalers.AggregatedList(project).Filter("name=" + name)
	if err := resp.Pages(ctx, func(page *compute.AutoscalerAggregatedList) error {
		for _, item := range page.Items {
			for _, i := range item.Autoscalers {
				autoscaler = *i
			}
		}
		return nil
	},
	); err != nil {
		plugin.Logger(ctx).Error("gcp_compute_autoscaler.getComputeAutoscaler", "api_err", err)
		return nil, err
	}

	// If the specified resource is not present, API does not return any not found errors
	if len(autoscaler.Name) < 1 {
		return nil, nil
	}

	return &autoscaler, nil
}

//// TRANSFORM FUNCTIONS

func autoscalerAka(_ context.Context, d *transform.TransformData) (interface{}, error) {
	i := d.HydrateItem.(*compute.Autoscaler)

	zoneName := getLastPathElement(types.SafeString(i.Zone))
	regionName := getLastPathElement(types.SafeString(i.Region))
	project := strings.Split(i.SelfLink, "/")[6]
	autoscalerName := types.SafeString(i.Name)

	var akas []string
	if zoneName == "" {
		akas = []string{"gcp://compute.googleapis.com/projects/" + project + "/regions/" + regionName + "/autoscalers/" + autoscalerName}
	} else {
		akas = []string{"gcp://compute.googleapis.com/projects/" + project + "/zones/" + zoneName + "/autoscalers/" + autoscalerName}
	}

	return akas, nil
}

func autoscalerLocation(_ context.Context, d *transform.TransformData) (interface{}, error) {
	i := d.HydrateItem.(*compute.Autoscaler)
	param := d.Param.(string)

	zoneName := getLastPathElement(types.SafeString(i.Zone))
	regionName := getLastPathElement(types.SafeString(i.Region))
	project := strings.Split(i.SelfLink, "/")[6]

	locationData := map[string]string{
		"Type":     "ZONAL",
		"Location": zoneName,
		"Project":  project,
	}

	if zoneName == "" {
		locationData["Type"] = "REGIONAL"
		locationData["Location"] = regionName
	}

	return locationData[param], nil
}
