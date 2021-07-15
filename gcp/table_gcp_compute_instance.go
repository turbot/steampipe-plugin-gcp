package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/go-kit/types"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func tableGcpComputeInstance(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_instance",
		Description: "GCP Compute Instance",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeInstance,
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeInstances,
		},
		Columns: []*plugin.Column{
			// commonly used columns
			{
				Name:        "name",
				Description: "The name of the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "machine_type_name",
				Description: "Name of the machine type resource for this instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MachineType").Transform(lastPathElement),
			},
			{
				Name:        "status",
				Description: "The status of the instance (PROVISIONING, STAGING, RUNNING, STOPPING, SUSPENDING, SUSPENDED, REPAIRING, and TERMINATED).",
				Type:        proto.ColumnType_STRING,
			},

			// other columns
			{
				Name:        "can_ip_forward",
				Description: "Allows this instance to send and receive packets with non-matching destination or source IPs. This is required if you plan to use this instance to forward routes.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "confidential_instance_config",
				Description: "Confidential VM detail for the instance, if applicable. Confidential VMs, now in beta, is the first product in Google Cloud’s Confidential Computing portfolio.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "cpu_platform",
				Description: "The CPU platform used by this instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_timestamp",
				Description: "Timestamp the instance was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "deletion_protection",
				Description: "Whether the resource should be protected against deletion.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "description",
				Description: "The instance description.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "disks",
				Description: "An Array of disks associated with this instance",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "display_device",
				Description: "Display device for the instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "fingerprint",
				Description: "Specifies a fingerprint for this resource, which is essentially a hash of the instance's contents and used for optimistic locking. The fingerprint is initially generated by Compute Engine and changes after every request to modify or update the instance. You must always provide an up-to-date fingerprint hash in order to update the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "guest_accelerators",
				Description: "A list of the type and count of accelerator cardsattached to the instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "hostname",
				Description: "The instance hostname.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The instance id.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "kind",
				Description: "Type of the resource. Always compute#instance for instances.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "label_fingerprint",
				Description: "A fingerprint for this request, which is essentially a hash of the label's contents and used for optimistic locking. The fingerprint is initially generated by Compute Engine and changes after every request to modify or update labels. You must always provide an up-to-date fingerprint hash in order to update or change labels.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "labels",
				Description: "Labels that apply to this instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:      "last_start_timestamp",
				Transform: transform.FromGo().NullIfZero(), Description: "Timestamp when the instance was last started.",
				Type: proto.ColumnType_TIMESTAMP,
			},
			{
				Name:      "last_stop_timestamp",
				Transform: transform.FromGo().NullIfZero(), Description: "Timestamp when the instance was last stopped.",
				Type: proto.ColumnType_TIMESTAMP,
			},
			{
				Name:      "last_suspended_timestamp",
				Transform: transform.FromGo().NullIfZero(), Description: "Timestamp when the instance was last suspended.",
				Type: proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "machine_type",
				Description: "Full or partial URL of the machine type resource for this instance, in the format: zones/zone/machineTypes/machine-type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "metadata",
				Description: "The metadata key/value pairs assigned to this instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "min_cpu_platform",
				Description: "Specifies a minimum CPU platform for the VM instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "network_interfaces",
				Description: "An array of network configurations for this instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "private_ipv6_google_access",
				Description: "The private IPv6 google access type for the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "reservation_affinity",
				Description: "Specifies the reservations that this instance can consume from.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "resource_policies",
				Description: "Resource policies applied to this instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "scheduling",
				Description: "The scheduling options for this instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for this resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_accounts",
				Description: "A list of service accounts, with their specified scopes, authorized for this instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "shielded_instance_config",
				Description: "Shielded instance configuration. Shielded VM provides verifiable integrity to prevent against malware and rootkits.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "shielded_instance_integrity_policy",
				Description: "Shielded instance integrity policy. Shielded instance configuration. Shielded VM provides verifiable integrity to prevent against malware and rootkits.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "start_restricted",
				Description: "Whether a VM has been restricted for start because Compute Engine has detected suspicious activity.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "status_message",
				Description: "An optional, human-readable explanation of the status.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "network_tags",
				Description: "Network tags applied to this instance. Network tags are used to identify valid sources or targets for network firewalls.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags.Items"),
			},
			{
				Name:        "zone",
				Description: "The zone in which the instance resides.",
				Type:        proto.ColumnType_STRING,
			},
			// zone_name is a simpler view of the zone, without the full path
			{
				Name:        "zone_name",
				Description: "The zone name in which the instance resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Zone").Transform(lastPathElement),
			},
			{
				Name:        "iam_policy",
				Description: "An Identity and Access Management (IAM) policy, which specifies access controls for Google Cloud resources. A `Policy` is a collection of `bindings`. A `binding` binds one or more `members` to a single `role`. Members can be user accounts, service accounts, Google groups, and domains (such as G Suite). A `role` is a named list of permissions; each `role` can be an IAM predefined role or a user-created custom role. For some types of Google Cloud resources, a `binding` can also specify a `condition`, which is a logical expression that allows access to a resource only if the expression evaluates to `true`.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getComputeInstanceIamPolicy,
				Transform:   transform.FromValue(),
			},

			// standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Labels"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(gcpComputeInstanceTurbotData, "Akas"),
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
				Transform:   transform.FromP(gcpComputeInstanceTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeInstances(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeInstances")

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

	resp := service.Instances.AggregatedList(project)
	if err := resp.Pages(
		ctx,
		func(page *compute.InstanceAggregatedList) error {
			for _, item := range page.Items {
				for _, instance := range item.Instances {
					d.StreamListItem(ctx, instance)
				}
			}
			return nil
		},
	); err != nil {
		if gerr, ok := err.(*googleapi.Error); ok {
			if helpers.StringSliceContains([]string{"403"}, types.ToString(gerr.Code)) {
				return nil, nil
			}
		}
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getComputeInstance(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getComputeInstance")

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

	var instance compute.Instance
	name := d.KeyColumnQuals["name"].GetStringValue()

	resp := service.Instances.AggregatedList(project).Filter("name=" + name)
	if err := resp.Pages(
		ctx,
		func(page *compute.InstanceAggregatedList) error {
			for _, item := range page.Items {
				for _, i := range item.Instances {
					instance = *i
				}
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	// If the specified resource is not present, API does not return any not found errors
	if len(instance.Name) < 1 {
		return nil, nil
	}

	return &instance, nil
}

func getComputeInstanceIamPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instance := h.Item.(*compute.Instance)

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	project := strings.Split(instance.SelfLink, "/")[6]
	zone := getLastPathElement(types.SafeString(instance.Zone))

	resp, err := service.Instances.GetIamPolicy(project, zone, instance.Name).Do()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

//// TRANSFORM FUNCTION

func gcpComputeInstanceTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	instance := d.HydrateItem.(*compute.Instance)
	param := d.Param.(string)

	zone := getLastPathElement(types.SafeString(instance.Zone))
	project := strings.Split(instance.SelfLink, "/")[6]

	turbotData := map[string]interface{}{
		"Project": project,
		"Akas":    []string{"gcp://compute.googleapis.com/projects/" + project + "/zones/" + zone + "/instances/" + instance.Name},
	}

	return turbotData[param], nil
}
