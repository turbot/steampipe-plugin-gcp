package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"google.golang.org/api/compute/v1"
)

//// TABLE DEFINITION

func tableGcpComputeRoute(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_route",
		Description: "GCP Compute Route",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeRoute,
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeRoutes,
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
				Name:        "kind",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A user-specified, human-readable description of the image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "The server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_timestamp",
				Description: "The creation timestamp of the resource.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "dest_range",
				Description: "The destination range of outgoing packets that this route applies to.",
				Type:        proto.ColumnType_IPADDR,
			},
			{
				Name:        "network",
				Description: "Fully-qualified URL of the network that this route applies to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "next_hop_gateway",
				Description: "The URL to a gateway that should handle matching packets.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "next_hop_ilb",
				Description: "The URL to a forwarding rule of type loadBalancingScheme INTERNAL that should handle matching packets or the IP address of the forwarding Rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "next_hop_instance",
				Description: "The URL to an instance that should handle matching packets.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "next_hop_ip",
				Description: "The network IP address of an instance that should handle matching packets.",
				Type:        proto.ColumnType_IPADDR,
			},
			{
				Name:        "next_hop_network",
				Description: "The URL of the local network if it should handle matching packets.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "next_hop_peering",
				Description: "The network peering name that should handle matching packets, which should conform to RFC1035.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "next_hop_vpn_tunnel",
				Description: "The URL to a VpnTunnel that should handle matching packets.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "priority",
				Description: "The priority of this route.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "tags",
				Description: "A list of instance tags to which this route applies.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "warnings",
				Description: "A list of warning messages, if potential misconfigurations are detected for this route.",
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
				Transform:   transform.From(gcpComputeRouteAka),
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

func listComputeRoutes(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeRoutes")
	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	project := activeProject()
	resp := service.Routes.List(project)
	if err := resp.Pages(ctx, func(page *compute.RouteList) error {
		for _, route := range page.Items {
			d.StreamListItem(ctx, route)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getComputeRoute(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	name := d.KeyColumnQuals["name"].GetStringValue()
	project := activeProject()

	// Error: pq: rpc error: code = Unknown desc = json: invalid use of ,string struct tag,
	// trying to unmarshal "projects/project/global/routes/" into uint64
	if len(name) < 1 {
		return nil, nil
	}

	req, err := service.Routes.Get(project, name).Do()
	if err != nil {
		return nil, err
	}

	return req, nil
}

//// TRANSFORM FUNCTIONS

func gcpComputeRouteAka(_ context.Context, d *transform.TransformData) (interface{}, error) {
	route := d.HydrateItem.(*compute.Route)

	// Build resource aka
	akas := []string{"gcp://compute.googleapis.com/projects/" + activeProject() + "/global/routes/" + route.Name}

	return akas, nil
}
