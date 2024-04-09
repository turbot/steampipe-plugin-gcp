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

func tableGcpComputeNetwork(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_network",
		Description: "GCP Compute Network",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeNetwork,
			Tags:       map[string]string{"service": "compute", "action": "networks.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeNetworks,
			KeyColumns: plugin.KeyColumnSlice{
				// Boolean columns
				{Name: "auto_create_subnetworks", Require: plugin.Optional, Operators: []string{"<>", "="}},
			},
			Tags: map[string]string{"service": "compute", "action": "networks.list"},
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
				Type:        proto.ColumnType_INT,
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
				Transform:   transform.FromP(gcpComputeNetworkTurbotData, "Akas"),
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
				Transform:   transform.FromP(gcpComputeNetworkTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeNetworks(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeNetworks")

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	filterQuals := []filterQualMap{
		{"auto_create_subnetworks", "autoCreateSubnetworks", "boolean"},
	}

	filters := buildQueryFilterFromQuals(filterQuals, d.Quals)
	filterString := ""
	if len(filters) > 0 {
		filterString = strings.Join(filters, " ")
	}

	// Max limit is set as per documentation
	// https://pkg.go.dev/google.golang.org/api@v0.48.0/compute/v1?utm_source=gopls#NetworksListCall.MaxResults
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

	resp := service.Networks.List(project).Filter(filterString).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *compute.NetworkList) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, network := range page.Items {
			d.StreamListItem(ctx, network)

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

//// HYDRATE FUNCTIONS

func getComputeNetwork(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getComputeNetwork")

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

	resp, err := service.Networks.Get(project, name).Do()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

//// TRANSFORM FUNCTIONS

func gcpComputeNetworkTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	network := d.HydrateItem.(*compute.Network)
	param := d.Param.(string)

	project := strings.Split(network.SelfLink, "/")[6]

	turbotData := map[string]interface{}{
		"Project": project,
		"Akas":    []string{"gcp://compute.googleapis.com/projects/" + project + "/global/networks/" + network.Name},
	}

	return turbotData[param], nil
}

func networkMtu(_ context.Context, d *transform.TransformData) (interface{}, error) {
	network := d.HydrateItem.(*compute.Network)

	if network.Mtu == 0 {
		return 1460, nil
	}

	return network.Mtu, nil
}
