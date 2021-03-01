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

type instanceGroupInfo = struct {
	InstanceGroup *compute.InstanceGroup
	IsManaged     bool
}

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
				Transform:   transform.FromField("InstanceGroup.Name"),
			},
			{
				Name:        "id",
				Description: "The unique identifier for the resource.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("InstanceGroup.Id"),
			},
			{
				Name:        "is_managed",
				Description: "Indicates whether the instance group is managed, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "description",
				Description: "An user-specified, human-readable description of the instance group.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getComputeInstanceGroup,
				Transform:   transform.FromField("InstanceGroup.Description"),
			},
			{
				Name:        "kind",
				Description: "Type of the resource",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InstanceGroup.Kind"),
			},
			{
				Name:        "size",
				Description: "Specifies the total number of instances in the instance group.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("InstanceGroup.Size"),
			},
			{
				Name:        "creation_timestamp",
				Description: "Timestamp the instance group was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("InstanceGroup.CreationTimestamp"),
			},
			{
				Name:        "fingerprint",
				Description: "Specifies fingerprint of the named ports. The system uses this fingerprint to detect conflicts when multiple users change the named ports concurrently.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InstanceGroup.Fingerprint"),
			},
			{
				Name:        "network",
				Description: "The URL of the network to which all instances in the instance group belong.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InstanceGroup.Network"),
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for this resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InstanceGroup.SelfLink"),
			},
			{
				Name:        "subnetwork",
				Description: "The URL of the subnetwork to which all instances in the instance group belong.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InstanceGroup.Subnetwork"),
			},
			{
				Name:        "zone",
				Description: "The URL of the zone where the instance group is located.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InstanceGroup.Zone"),
			},
			// zone_name is a simpler view of the zone, without the full path
			{
				Name:        "zone_name",
				Description: "The zone name in which the instance resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InstanceGroup.Zone").Transform(lastPathElement),
			},
			{
				Name:        "named_ports",
				Description: "A list of assignments of name to a port number.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("InstanceGroup.NamedPorts"),
			},
			{
				Name:        "instances",
				Description: "A list of instances present inside the instance group.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getComputeInstanceGroupInstancesList,
				Transform:   transform.FromValue(),
			},

			// standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InstanceGroup.Name"),
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
				Transform:   transform.FromField("InstanceGroup.Zone").Transform(lastPathElement),
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

//// LIST FUNCTION

func listComputeInstanceGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeInstanceGroups")

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

	var managedInstanceGroups []string

	// List all managed instance groups
	resp := service.InstanceGroupManagers.AggregatedList(project)
	if err := resp.Pages(
		ctx,
		func(page *compute.InstanceGroupManagerAggregatedList) error {
			for _, item := range page.Items {
				for _, instanceGroup := range item.InstanceGroupManagers {
					managedInstanceGroups = append(managedInstanceGroups, instanceGroup.Name)
					d.StreamListItem(ctx, instanceGroupInfo{&compute.InstanceGroup{
						CreationTimestamp: instanceGroup.CreationTimestamp,
						Description:       instanceGroup.Description,
						Fingerprint:       instanceGroup.Fingerprint,
						Id:                instanceGroup.Id,
						Kind:              instanceGroup.Kind,
						Name:              instanceGroup.Name,
						NamedPorts:        instanceGroup.NamedPorts,
						SelfLink:          instanceGroup.SelfLink,
						Size:              calculateTotalInstance(instanceGroup),
						Zone:              instanceGroup.Zone,
					}, true})
				}
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	// List all unmanaged instance groups
	result := service.InstanceGroups.AggregatedList(project)
	if err := result.Pages(
		ctx,
		func(page *compute.InstanceGroupAggregatedList) error {
			for _, item := range page.Items {
				for _, instanceGroup := range item.InstanceGroups {
					if !compareData(managedInstanceGroups, instanceGroup.Name) {
						d.StreamListItem(ctx, instanceGroupInfo{instanceGroup, false})
					}
				}
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getComputeInstanceGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getComputeInstanceGroup")

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

	var instanceGroup *instanceGroupInfo
	var name string
	isManaged := true

	if h.Item != nil {
		name = h.Item.(instanceGroupInfo).InstanceGroup.Name
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
	}

	resp := service.InstanceGroupManagers.AggregatedList(project).Filter("name=" + name)
	if err := resp.Pages(
		ctx,
		func(page *compute.InstanceGroupManagerAggregatedList) error {
			for _, item := range page.Items {
				if len(item.InstanceGroupManagers) == 0 {
					isManaged = false
				}
				for _, i := range item.InstanceGroupManagers {
					instanceGroup = &instanceGroupInfo{&compute.InstanceGroup{
						CreationTimestamp: i.CreationTimestamp,
						Description:       i.Description,
						Fingerprint:       i.Fingerprint,
						Id:                i.Id,
						Kind:              i.Kind,
						Name:              i.Name,
						NamedPorts:        i.NamedPorts,
						SelfLink:          i.SelfLink,
						Size:              calculateTotalInstance(i),
						Zone:              i.Zone,
					},
						isManaged,
					}
				}
			}
			return nil
		}); err != nil {
		return nil, err
	}

	// Get unmanaged instance group
	if !isManaged {
		resp := service.InstanceGroups.AggregatedList(project).Filter("name=" + name)
		if err := resp.Pages(
			ctx,
			func(page *compute.InstanceGroupAggregatedList) error {
				for _, item := range page.Items {
					for _, i := range item.InstanceGroups {
						instanceGroup = &instanceGroupInfo{&compute.InstanceGroup{
							CreationTimestamp: i.CreationTimestamp,
							Description:       i.Description,
							Fingerprint:       i.Fingerprint,
							Id:                i.Id,
							Kind:              i.Kind,
							Name:              i.Name,
							NamedPorts:        i.NamedPorts,
							SelfLink:          i.SelfLink,
							Size:              i.Size,
							Zone:              i.Zone,
						},
							isManaged,
						}
					}
				}
				return nil
			}); err != nil {
			return nil, err
		}
	}

	// If the specified resource is not present, API does not return any not found errors
	if instanceGroup == nil {
		return nil, nil
	}

	return instanceGroup, nil
}

func getComputeInstanceGroupInstancesList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	data := h.Item.(*instanceGroupInfo)
	splittedTitle := strings.Split(data.InstanceGroup.SelfLink, "/")
	project := splittedTitle[6]
	zone := getLastPathElement(types.SafeString(data.InstanceGroup.Zone))

	// List instances for unmanaged group
	if !data.IsManaged {
		// Build param
		listInstanceRequest := &compute.InstanceGroupsListInstancesRequest{
			InstanceState: "ALL",
		}
		var instances []*compute.InstanceWithNamedPorts
		resp := service.InstanceGroups.ListInstances(project, zone, data.InstanceGroup.Name, listInstanceRequest)
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

	// List instances for managed group
	var instances []*compute.ManagedInstance
	resp := service.InstanceGroupManagers.ListManagedInstances(project, zone, data.InstanceGroup.Name)
	if err := resp.Pages(
		ctx,
		func(page *compute.InstanceGroupManagersListManagedInstancesResponse) error {
			for _, item := range page.ManagedInstances {
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
	data := d.HydrateItem.(*instanceGroupInfo)
	param := d.Param.(string)

	zone := getLastPathElement(types.SafeString(data.InstanceGroup.Zone))
	project := strings.Split(data.InstanceGroup.SelfLink, "/")[6]

	turbotData := map[string]interface{}{
		"Project": project,
		"Akas":    []string{"gcp://compute.googleapis.com/projects/" + project + "/zones/" + zone + "/instanceGroups/" + data.InstanceGroup.Name},
	}

	return turbotData[param], nil
}

func calculateTotalInstance(data *compute.InstanceGroupManager) int64 {
	i := data.CurrentActions
	return (i.None + i.Creating + i.CreatingWithoutRetries + i.Verifying + i.Recreating + i.Deleting + i.Abandoning + i.Restarting + i.Refreshing)
}

func compareData(managedGroups []string, latest string) bool {
	for _, group := range managedGroups {
		if group == latest {
			return true
		}
	}
	return false
}
