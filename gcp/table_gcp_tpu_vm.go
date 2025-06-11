package gcp

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	tpu "google.golang.org/api/tpu/v2"
)

func tableGcpTpuVM(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_tpu_vm",
		Description: "GCP TPU VMs are specialized hardware accelerators designed to speed up specific machine learning workloads.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getTpuVM,
			Tags:       map[string]string{"service": "tpu", "action": "nodes.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listTpuVMs,
			Tags:    map[string]string{"service": "tpu", "action": "nodes.list"},
		},
		Columns: []*plugin.Column{
			// Key columns
			{
				Name:        "name",
				Description: "The name of the TPU node.",
				Type:        proto.ColumnType_STRING,
			},

			// Other columns
			{
				Name:        "id",
				Description: "The unique identifier for the TPU node.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "description",
				Description: "The user-supplied description of the TPU node.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "accelerator_type",
				Description: "The type of TPU accelerator.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The current state of the TPU node.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "health",
				Description: "The health status of the TPU node.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "health_description",
				Description: "If the TPU node is unhealthy, this contains more detailed information about why.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "network",
				Description: "The name of the network that the TPU node is connected to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cidr_block",
				Description: "The CIDR block that the TPU node will use when selecting an IP address.",
				Type:        proto.ColumnType_CIDR,
			},
			{
				Name:        "service_account",
				Description: "The service account used to run the TPU node.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "The time when the TPU node was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "api_version",
				Description: "The API version that created this node.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "runtime_version",
				Description: "The runtime version running in the node.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "scheduling_config",
				Description: "Sets the scheduling options for the TPU instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "network_endpoints",
				Description: "The network endpoints where the TPU node can be accessed.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "labels",
				Description: "Resource labels to represent user provided metadata.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "metadata",
				Description: "Custom metadata to apply to the TPU node.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "multislice_node",
				Description: "Whether the Node belongs to a Multislice group.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "queued_resource",
				Description: "The qualified name of the QueuedResource that requested this Node.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "data_disks",
				Description: "The additional data disks for the Node.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "symptoms",
				Description: "The Symptoms that have occurred to the TPU Node.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "zone",
				Description: "The GCP zone where the TPU node is located.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpTpuVMTurbotData, "Zone"),
			},
			{
				Name:        "network_config",
				Description: "The network configuration for the TPU node.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "accelerator_config",
				Description: "The accelerator configuration for the TPU node.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "shielded_instance_config",
				Description: "The shielded instance configuration for the TPU node.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags",
				Description: "A map of tags for the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Labels"),
			},

			// Standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(lastPathElement),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(gcpTpuVMTurbotData, "Akas"),
			},

			// Standard GCP columns
			{
				Name:        "project",
				Description: "The GCP project ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpTpuVMTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listTpuVMs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := TPUService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_tpu_vm.listTpuVMs", "connection_error", err)
		return nil, err
	}

	// Get project details
	projectId, err := getProject(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_tpu_vm.listTpuVMs", "project_error", err)
		return nil, err
	}

	// Use locations/- to get all TPUs across all regions in a single request
	parent := "projects/" + projectId.(string) + "/locations/-"

	resp := service.Projects.Locations.Nodes.List(parent)
	if err := resp.Pages(ctx, func(page *tpu.ListNodesResponse) error {
		for _, node := range page.Nodes {
			d.StreamListItem(ctx, node)

			// Check if context has been cancelled or if the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil
			}
		}
		return nil
	}); err != nil {
		plugin.Logger(ctx).Error("gcp_tpu_vm.listTpuVMs", "api_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getTpuVM(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := TPUService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_tpu_vm.getTpuVM", "connection_error", err)
		return nil, err
	}

	name := d.EqualsQualString("name")
	if len(name) < 1 {
		return nil, nil
	}

	node, err := service.Projects.Locations.Nodes.Get(name).Do()
	if err != nil {
		plugin.Logger(ctx).Error("gcp_tpu_vm.getTpuVM", "api_error", err)
		return nil, err
	}

	return node, nil
}

//// TRANSFORM FUNCTIONS

func gcpTpuVMTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	node := d.HydrateItem.(*tpu.Node)
	param := d.Param.(string)

	// Extract project and zone from the node name
	// Format: projects/{project}/locations/{zone}/nodes/{node}
	parts := strings.Split(node.Name, "/")
	if len(parts) != 6 {
		return nil, nil
	}

	project := parts[1]
	zone := parts[3]

	turbotData := map[string]interface{}{
		"Project": project,
		"Zone":    zone,
		"Akas":    []string{"gcp://tpu.googleapis.com/" + node.Name},
	}

	return turbotData[param], nil
}
