package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

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
			Hydrate: listComputeInstanceGroups,
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
				Description: "An user-specified, human-readable description of the instance group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "Type of the resource",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "size",
				Description: "Specifies the total number of instances in the instance group.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "creation_timestamp",
				Description: "Timestamp the instance group was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "fingerprint",
				Description: "Specifies fingerprint of the named ports. The system uses this fingerprint to detect conflicts when multiple users change the named ports concurrently.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "network",
				Description: "The URL of the network to which all instances in the instance group belong.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "region",
				Description: "The URL of the region where the instance group is located.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for this resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "subnetwork",
				Description: "The URL of the subnetwork to which all instances in the instance group belong.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "zone",
				Description: "The URL of the zone where the instance group is located.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "named_ports",
				Description: "A list of assignments of name to a port number.",
				Type:        proto.ColumnType_JSON,
			},

			// zone_name is a simpler view of the zone, without the full path
			{
				Name:        "zone_name",
				Description: "The zone name in which the instance resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Zone").Transform(lastPathElement),
			},
			// {
			// 	Name:        "instances",
			// 	Description: "A list of instances present inside the instance group.",
			// 	Type:        proto.ColumnType_JSON,
			// 	Hydrate:     getComputeInstanceGroupInstancesList,
			// 	Transform:   transform.FromValue(),
			// },

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
				Transform:   transform.FromP(computeInstanceGroupSelfLinkToTurbotData, "Akas"),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Zone").Transform(lastPathElement),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(computeInstanceGroupSelfLinkToTurbotData, "Project"),
			},
		},
	}
}

//// HYDRATE FUNCTIONS

func listComputeInstanceGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listComputeInstanceGroups")

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

	resp := service.InstanceGroups.AggregatedList(project)
	if err := resp.Pages(
		ctx,
		func(page *compute.InstanceGroupAggregatedList) error {
			for _, item := range page.Items {
				for _, instanceGroup := range item.InstanceGroups {
					d.StreamListItem(ctx, instanceGroup)
				}
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	return nil, nil
}

func getComputeInstanceGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getComputeInstanceGroup")

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

	var instanceGroup compute.InstanceGroup
	name := d.KeyColumnQuals["name"].GetStringValue()

	resp := service.InstanceGroups.AggregatedList(project).Filter("name=" + name)
	if err := resp.Pages(
		ctx,
		func(page *compute.InstanceGroupAggregatedList) error {
			for _, item := range page.Items {
				for _, i := range item.InstanceGroups {
					instanceGroup = *i
				}
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	// If the specified resource is not present, API does not return any not found errors
	if len(instanceGroup.Name) < 1 {
		return nil, nil
	}

	return &instanceGroup, nil
}

func getComputeInstanceGroupInstancesList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	instanceGroup := h.Item.(*compute.InstanceGroup)
	splittedTitle := strings.Split(instanceGroup.SelfLink, "/")

	// Build param
	project := splittedTitle[1]
	zone := getLastPathElement(types.SafeString(instanceGroup.Zone))
	listInstanceRequest := &compute.InstanceGroupsListInstancesRequest{
		InstanceState: "ALL",
	}

	var instances []*compute.InstanceWithNamedPorts
	resp := service.InstanceGroups.ListInstances(project, zone, instanceGroup.Name, listInstanceRequest)
	if err := resp.Pages(
		ctx,
		func(page *compute.InstanceGroupsListInstances) error {
			for _, item := range page.Items {
				instances = append(instances, item)
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	return instances, nil
}

//// TRANSFORM FUNCTIONS

func computeInstanceGroupSelfLinkToTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	instanceGroup := d.HydrateItem.(*compute.InstanceGroup)
	param := d.Param.(string)

	zone := getLastPathElement(types.SafeString(instanceGroup.Zone))
	project := strings.Split(instanceGroup.SelfLink, "/")[1]

	turbotData := map[string]interface{}{
		"Project": project,
		"Akas":    []string{"gcp://compute.googleapis.com/projects/" + project + "/zones/" + zone + "/instanceGroups/" + instanceGroup.Name},
	}

	return turbotData[param], nil
}
