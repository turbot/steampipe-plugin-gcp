package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	"google.golang.org/api/compute/v1"
)

func tableGcpComputeNetwork(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_network",
		Description: "GCP Compute Network",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeNetwork,
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeNetworks,
		},
		Columns: []*plugin.Column{
			// commonly used columns
			{
				Name:        "name",
				Description: "Name of the resource. Provided by the client when the resource is created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique identifier for the resource. This identifier is defined by the server.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "creation_timestamp",
				Description: "Creation timestamp in RFC3339 text format.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "auto_create_subnetworks",
				Description: "When set to true, the VPC network is created in auto mode. When set to false, the VPC network is created in custom mode.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "description",
				Description: "An optional description of this resource. Provide this field when you create the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "gateway_ipv4",
				Description: "The gateway address for default routing out of the network, selected by GCP",
				Type:        proto.ColumnType_IPADDR,
			},
			{
				Name:        "ipv4_range",
				Description: "The range of internal addresses that are legal on this network. Deprecated in favor of subnet mode networks. This range is a CIDR specification, for example: 192.168.0.0/16. Provided by the client when the network is created.",
				Type:        proto.ColumnType_CIDR,
			},
			{
				Name:        "kind",
				Description: "Type of the resource. Always compute#network for networks.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "mtu",
				Description: "Maximum Transmission Unit in bytes. The minimum value for this field is 1460 and the maximum value is 1500 bytes.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.From(networkMtu),
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "routing_mode",
				Description: "The network-wide routing mode to use. If set to REGIONAL, this network's Cloud Routers will only advertise routes with subnets of this network in the same region as the router. If set to GLOBAL, this network's Cloud Routers will advertise routes with all subnets of this network, across regions. Possible values: \"GLOBAL\"   \"REGIONAL\"",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RoutingConfig.RoutingMode"),
			},
			{
				Name:        "peerings",
				Description: "A list of network peerings for the resource. NetworkPeering: A network peering attached to a network resource. The message includes the peering name, peer network, peering state, and a flag indicating whether Google Compute Engine should automatically create routes for the peering",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "subnetworks",
				Description: "Server-defined fully-qualified URLs for all subnetworks in this VPC network.",
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
				Transform:   transform.From(networkAka),
			},

			// standard gcp columns
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromConstant(activeProject()),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeNetworks(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeNetworks")
	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	project := activeProject()
	resp := service.Networks.List(project)
	if err := resp.Pages(ctx, func(page *compute.NetworkList) error {
		for _, network := range page.Items {
			plugin.Logger(ctx).Trace("getComputeNetwork   ~~~~~~~~~~~~~~~", " DATA", network.Mtu)
			d.StreamListItem(ctx, network)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getComputeNetwork(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getComputeNetwork")

	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	project := activeProject()
	name := d.KeyColumnQuals["name"].GetStringValue()

	resp, err := service.Networks.Get(project, name).Do()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

//// TRANSFORM FUNCTIONS

func networkAka(_ context.Context, d *transform.TransformData) (interface{}, error) {
	network := d.HydrateItem.(*compute.Network)

	akas := []string{"gcp://compute.googleapis.com/projects/" + activeProject() + "/global/networks/" + network.Name}

	return akas, nil
}

func networkMtu(_ context.Context, d *transform.TransformData) (interface{}, error) {
	network := d.HydrateItem.(*compute.Network)

	if network.Mtu == 0 {
		return 1460, nil
	}

	return network.Mtu, nil
}
