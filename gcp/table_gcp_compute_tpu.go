package gcp

import (
	"context"
	"errors"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGcpComputeTpu(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_tpu",
		Description: "[DEPRECATED] GCP Compute TPUs are specialized hardware accelerators designed to speed up specific machine learning workloads.",
		List: &plugin.ListConfig{
			Hydrate: listComputeTpus,
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
			},

			// Standard GCP columns
			{
				Name:        "project",
				Description: "The GCP project ID.",
				Type:        proto.ColumnType_STRING,
			},
		},
	}
}

//// LIST FUNCTION

func listComputeTpus(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	err := errors.New("The gcp_compute_tpu table has been deprecated and removed, please use gcp_tpu_vm table instead.")
	return nil, err
}
