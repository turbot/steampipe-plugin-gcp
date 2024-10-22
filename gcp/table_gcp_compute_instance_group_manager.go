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

func tableGcpComputeInstanceGroupManager(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_instance_group_manager",
		Description: "GCP Compute Instance Group Manager",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeInstanceGroupManager,
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeInstanceGroupManager,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the instance group manager.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique identifier for this instance group manager. This identifier is defined by the server.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "all_instances_config",
				Description: "Specifies configuration that overrides the instance template configuration for the group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromGo().NullIfZero(),
			},
			{
				Name:        "auto_healing_policies",
				Description: "The autohealing policy for this managed instance group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromGo().NullIfZero(),
			},
			{
				Name:        "base_instance_name",
				Description: "The base instance name is a prefix that you want to attach to the names of all VMs in a MIG.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_timestamp",
				Description: "The timestamp when the instance group manager was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "current_actions",
				Description: "The list of instance actions and the number of instances in this managed instance group that are scheduled for each of those actions.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromGo().NullIfZero(),
			},
			{
				Name:        "description",
				Description: "An optional description of this resource. Provide this property when you create the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "distribution_policy",
				Description: "The Policy specifying the intended distribution of managed instances across zones in a regional managed instance group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromGo().NullIfZero(),
			},
			{
				Name:        "fingerprint",
				Description: "The fingerprint of the instance group manager.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_group",
				Description: "The instance group that is managed by this group manager.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "instance_lifecycle_policy",
				Description: "The repair policy for this managed instance group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromGo().NullIfZero(),
			},
			{
				Name:        "instance_template",
				Description: "The URL of the instance template that is specified for this managed instance group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "The type of the resource. Always compute#instanceGroupManager for instance group managers.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "named_ports",
				Description: "The named ports configured for the Instance Groups complementary to this Instance Group Manager.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromGo().NullIfZero(),
			},
			{
				Name:        "region",
				Description: "The URL of the region where the instance group manager resides. Only applicable for regional resources.",
				Type:        proto.ColumnType_STRING,
			},

			// region_name is a simpler view of the region, without the full path
			{
				Name:        "region_name",
				Description: "The name of the region where the instance group manager resides. Only applicable for regional resources.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Region").Transform(lastPathElement),
			},
			{
				Name:        "self_link",
				Description: "The server-defined fully-qualified URL for this resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "stateful_policy",
				Description: "The stateful policy for this managed instance group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromGo().NullIfZero(),
			},
			{
				Name:        "status",
				Description: "The status of this managed instance group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromGo().NullIfZero(),
			},
			{
				Name:        "target_pools",
				Description: "The URLs for all TargetPool resources to which instances in the instanceGroup field are added.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromGo().NullIfZero(),
			},
			{
				Name:        "target_size",
				Description: "The target number of running instances for this managed instance group.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "update_policy",
				Description: "The update policy for this managed instance group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromGo().NullIfZero(),
			},
			{
				Name:        "versions",
				Description: "The versions of instance templates used by this managed instance group to create instances, useful for canary updates. Overrides the top-level instance template.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromGo().NullIfZero(),
			},
			{
				Name:        "zone",
				Description: "The URL of the zone where the instance group manager resides.",
				Type:        proto.ColumnType_STRING,
			},

			// zone_name is a simpler view of the zone, without the full path
			{
				Name:        "zone_name",
				Description: "The zone name in which the instance group manager resides.",
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
				Transform:   transform.From(instanceGroupManagerAka),
			},

			// GCP standard columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(instanceGroupManagerLocation, "Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(instanceGroupManagerLocation, "Project"),
			},
		},
	}
}

//// LIST FUNCTIONS

func listComputeInstanceGroupManager(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
		plugin.Logger(ctx).Error("gcp_compute_instance_group_manager.listComputeInstanceGroupManager", "service_creation_err", err)
		return nil, err
	}

	resp := service.InstanceGroupManagers.AggregatedList(project).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *compute.InstanceGroupManagerAggregatedList) error {
		for _, item := range page.Items {
			for _, group := range item.InstanceGroupManagers {
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
		plugin.Logger(ctx).Error("gcp_compute_instance_group_manager.listComputeInstanceGroupManager", "api_err", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getComputeInstanceGroupManager(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get project details
	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	var group compute.InstanceGroupManager
	name := d.EqualsQuals["name"].GetStringValue()

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_compute_instance_group_manager.getComputeInstanceGroupManager", "service_creation_err", err)
		return nil, err
	}

	resp := service.InstanceGroupManagers.AggregatedList(project).Filter("name=" + name)
	if err := resp.Pages(ctx, func(page *compute.InstanceGroupManagerAggregatedList) error {
		for _, item := range page.Items {
			for _, i := range item.InstanceGroupManagers {
				group = *i
			}
		}
		return nil
	},
	); err != nil {
		plugin.Logger(ctx).Error("gcp_compute_instance_group_manager.getComputeInstanceGroupManager", "api_err", err)
		return nil, err
	}

	// If the specified resource is not present, API does not return any not found errors
	if len(group.Name) < 1 {
		return nil, nil
	}

	return &group, nil
}

//// TRANSFORM FUNCTIONS

func instanceGroupManagerAka(_ context.Context, d *transform.TransformData) (interface{}, error) {
	i := d.HydrateItem.(*compute.InstanceGroupManager)

	zoneName := getLastPathElement(types.SafeString(i.Zone))
	regionName := getLastPathElement(types.SafeString(i.Region))
	project := strings.Split(i.SelfLink, "/")[6]
	instanceGroupManagerName := types.SafeString(i.Name)

	var akas []string
	if zoneName == "" {
		akas = []string{"gcp://compute.googleapis.com/projects/" + project + "/regions/" + regionName + "/instanceGroupManagers/" + instanceGroupManagerName}
	} else {
		akas = []string{"gcp://compute.googleapis.com/projects/" + project + "/zones/" + zoneName + "/instanceGroupManagers/" + instanceGroupManagerName}
	}

	return akas, nil
}

func instanceGroupManagerLocation(_ context.Context, d *transform.TransformData) (interface{}, error) {
	i := d.HydrateItem.(*compute.InstanceGroupManager)
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
