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

func tableGcpComputeInstanceGroup(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_instance_group",
		Description: "GCP Compute Instance Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeInstanceGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeInstanceGroup,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "Name of the instance group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique identifier for this instance group. This identifier is defined by the server.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "creation_timestamp",
				Description: "Timestamp when the instance group was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "An optional description of this resource. Provide this property when you create the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "fingerprint",
				Description: "The fingerprint of the instance group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "Type of the resource. Always compute#instanceGroup for instance groups.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "network",
				Description: "The URL of the network to which all instances in the instance group belong.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo().NullIfZero(),
			},
			{
				Name:        "region",
				Description: "The URL of the region where the instance group resides. Only applicable for regional resources.",
				Type:        proto.ColumnType_STRING,
			},

			// region_name is a simpler view of the region, without the full path
			{
				Name:        "region_name",
				Description: "Name of the region where the instance group resides. Only applicable for regional resources.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Region").Transform(lastPathElement),
			},
			{
				Name:        "self_link",
				Description: "Server-defined fully-qualified URL for this resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "size",
				Description: "The total number of instances in the instance group.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "subnetwork",
				Description: "The URL of the subnetwork to which all instances in the instance group belong.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "zone",
				Description: "The URL of the zone where the instance group resides.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "named_ports",
				Description: "Assigns a name to a port number.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromGo().NullIfZero(),
			},
			{
				Name:        "instances",
				Description: "List of instances that belongs to this group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getComputeInstanceGroupInstances,
				Transform:   transform.FromValue(),
			},

			// zone_name is a simpler view of the zone, without the full path
			{
				Name:        "zone_name",
				Description: "The zone name in which the instance group resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Zone").Transform(lastPathElement),
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
				Transform:   transform.From(instanceGroupAka),
			},

			// GCP standard columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(instanceGroupLocation, "Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(instanceGroupLocation, "Project"),
			},
		},
	}
}

//// LIST FUNCTIONS

func listComputeInstanceGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
		plugin.Logger(ctx).Error("gcp_compute_instance_group.listComputeInstanceGroup", "service_creation_err", err)
		return nil, err
	}

	resp := service.InstanceGroups.AggregatedList(project).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *compute.InstanceGroupAggregatedList) error {
		for _, item := range page.Items {
			for _, group := range item.InstanceGroups {
				d.StreamListItem(ctx, group)

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
		plugin.Logger(ctx).Error("gcp_compute_instance_group.listComputeInstanceGroup", "api_err", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getComputeInstanceGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getComputeInstanceGroup")

	// Get project details
	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	var group compute.InstanceGroup
	name := d.EqualsQuals["name"].GetStringValue()

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_compute_instance_group.getComputeInstanceGroup", "service_creation_err", err)
		return nil, err
	}

	resp := service.InstanceGroups.AggregatedList(project).Filter("name=" + name)
	if err := resp.Pages(ctx, func(page *compute.InstanceGroupAggregatedList) error {
		for _, item := range page.Items {
			for _, i := range item.InstanceGroups {
				group = *i
			}
		}
		return nil
	},
	); err != nil {
		plugin.Logger(ctx).Error("gcp_compute_instance_group.getComputeInstanceGroup", "api_err", err)
		return nil, err
	}

	// If the specified resource is not present, API does not return any not found errors
	if len(group.Name) < 1 {
		return nil, nil
	}

	return &group, nil
}

func getComputeInstanceGroupInstances(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instanceGroup := h.Item.(*compute.InstanceGroup)

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_compute_instance_group.getComputeInstanceGroupInstances", "service_creation_err", err)
		return nil, err
	}

	resp, err := service.InstanceGroups.ListInstances(project, getLastPathElement(types.SafeString(instanceGroup.Zone)), instanceGroup.Name, &compute.InstanceGroupsListInstancesRequest{}).Do()

	if err != nil {
		plugin.Logger(ctx).Error("gcp_compute_instance_group.getComputeInstanceGroupInstances", "api_err", err)
		return nil, err
	}

	return &resp.Items, nil
}

//// TRANSFORM FUNCTIONS

func instanceGroupAka(_ context.Context, d *transform.TransformData) (interface{}, error) {
	i := d.HydrateItem.(*compute.InstanceGroup)

	zoneName := getLastPathElement(types.SafeString(i.Zone))
	regionName := getLastPathElement(types.SafeString(i.Region))
	project := strings.Split(i.SelfLink, "/")[6]
	instanceGroupName := types.SafeString(i.Name)

	var akas []string
	if zoneName == "" {
		akas = []string{"gcp://compute.googleapis.com/projects/" + project + "/regions/" + regionName + "/instanceGroups/" + instanceGroupName}
	} else {
		akas = []string{"gcp://compute.googleapis.com/projects/" + project + "/zones/" + zoneName + "/instanceGroups/" + instanceGroupName}
	}

	return akas, nil
}

func instanceGroupLocation(_ context.Context, d *transform.TransformData) (interface{}, error) {
	i := d.HydrateItem.(*compute.InstanceGroup)
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
