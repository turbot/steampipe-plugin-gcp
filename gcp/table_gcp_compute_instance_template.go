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

func tableGcpComputeInstanceTemplate(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_instance_template",
		Description: "GCP Compute Instance Template",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeInstanceTemplate,
			Tags:       map[string]string{"service": "monitoring", "action": "instanceTemplates.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeInstanceTemplate,
			Tags:    map[string]string{"service": "monitoring", "action": "instanceTemplates.list"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "A friendly name that identifies the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "A unique identifier for this instance template. The server defines this identifier.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "creation_timestamp",
				Description: "The creation timestamp for this instance template.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "self_link",
				Description: "The server-defined URL for this instance template.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "The resource type, which is always compute#instanceTemplate for instance templates.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "An optional description of this resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_description",
				Description: "An optional text description for the instances that are created from these properties.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Description"),
			},
			{
				Name:        "instance_machine_type",
				Description: "The machine type to use for instances that are created from these properties.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.MachineType"),
			},
			{
				Name:        "instance_can_ip_forward",
				Description: "Enables instances created based on these properties to send packets with source IP addresses other than their own and receive packets with destination IP addresses other than their own.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.CanIpForward"),
			},
			{
				Name:        "instance_min_cpu_platform",
				Description: "Minimum cpu/platform to be used by instances. The instance may be scheduled on the specified or newer cpu/platform. Applicable values are the friendly names of CPU platforms, such as minCpuPlatform: \"Intel Haswell\" or minCpuPlatform: \"Intel Sandy Bridge\". For more information, read Specifying a Minimum CPU Platform.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.MinCpuPlatform"),
			},
			{
				Name:        "instance_private_ipv6_google_access",
				Description: "The private IPv6 google access type for VMs. If not specified, use INHERIT_FROM_SUBNETWORK as default. Possible values: \"ENABLE_BIDIRECTIONAL_ACCESS_TO_GOOGLE\", \"ENABLE_OUTBOUND_VM_ACCESS_TO_GOOGLE\", \"INHERIT_FROM_SUBNETWORK\"",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.PrivateIpv6GoogleAccess"),
			},
			{
				Name:        "source_instance",
				Description: "The URL of the source instance used to create the template.",
				Type:        proto.ColumnType_STRING,
			},
			// source_instance_name is a simpler view of the source_instance, without the full path
			{
				Name:        "source_instance_name",
				Description: "The URL of the source instance used to create the template.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SourceInstance").Transform(lastPathElement),
			},
			{
				Name:        "instance_disks",
				Description: "An array of disks that are associated with the instances that are created from these properties.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Disks"),
			},
			{
				Name:        "instance_guest_accelerators",
				Description: "A list of guest accelerator cards' type and count to use for instances created from these properties.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.GuestAccelerators"),
			},
			{
				Name:        "instance_metadata",
				Description: "The metadata key/value pairs to assign to instances that are created from these properties. These pairs can consist of custom metadata or predefined keys.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Metadata"),
			},
			{
				Name:        "instance_network_interfaces",
				Description: "An array of network access configurations for this interface.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.NetworkInterfaces"),
			},
			{
				Name:        "instance_reservation_affinity",
				Description: "Specifies the reservations that instances can consume from.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.ReservationAffinity"),
			},
			{
				Name:        "instance_resource_policies",
				Description: "Resource policies (names, not URLs) applied to instances created from these properties.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.ResourcePolicies"),
			},
			{
				Name:        "instance_scheduling",
				Description: "Specifies the scheduling options for the instances that are created from these properties.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Scheduling"),
			},
			{
				Name:        "instance_service_accounts",
				Description: "A list of service accounts with specified scopes. Access tokens for these service accounts are available to the instances that are created from these properties.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.ServiceAccounts"),
			},
			{
				Name:        "instance_shielded_instance_config",
				Description: "A set of Shielded Instance options.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.ShieldedInstanceConfig"),
			},
			{
				Name:        "instance_tags",
				Description: "A list of tags to apply to the instances that are created from these properties. The tags identify valid sources or targets for network firewalls.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Tags"),
			},
			{
				Name:        "labels",
				Description: "Labels to apply to instances that are created from these properties.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Labels"),
			},

			// common resource columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Labels"),
			},
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
				Transform:   transform.FromP(gcpComputeInstanceTemplateTurbotData, "Akas"),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromConstant("global"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpComputeInstanceTemplateTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeInstanceTemplate(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Max limit is set as per documentation
	// https://pkg.go.dev/google.golang.org/api@v0.48.0/compute/v1?utm_source=gopls#InstanceTemplatesListCall.MaxResults
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

	resp := service.InstanceTemplates.List(project).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *compute.InstanceTemplateList) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, template := range page.Items {
			d.StreamListItem(ctx, template)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				page.NextPageToken = ""
				return nil
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTION

func getComputeInstanceTemplate(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	name := d.EqualsQuals["name"].GetStringValue()

	// Error: pq: rpc error: code = Unknown desc = json: invalid use of ,string struct tag,
	// trying to unmarshal "projects/project/global/instanceTemplates/" into uint64
	if len(name) < 1 {
		return nil, nil
	}

	resp, err := service.InstanceTemplates.Get(project, name).Do()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

//// TRANSFORM FUNCTION

func gcpComputeInstanceTemplateTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	instanceTemplate := d.HydrateItem.(*compute.InstanceTemplate)
	param := d.Param.(string)

	project := strings.Split(instanceTemplate.SelfLink, "/")[6]

	turbotData := map[string]interface{}{
		"Project": project,
		"Akas":    []string{"gcp://compute.googleapis.com/projects/" + project + "/global/instanceTemplates/" + instanceTemplate.Name},
	}

	return turbotData[param], nil
}
