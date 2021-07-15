package gcp

import (
	"context"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

//// TABLE DEFINITION

func tableGcpComputeVpnTunnel(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_vpn_tunnel",
		Description: "GCP Compute VPN Tunnel",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeVpnTunnel,
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeVpnTunnels,
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
				Description: "A user-specified, human-readable description of the vpn tunnel.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "Specifies the current status of the vpn tunnel.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_timestamp",
				Description: "The creation timestamp of the resource.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "detailed_status",
				Description: "Detailed status message for the VPN tunnel.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ike_version",
				Description: "Specifies the IKE protocol version to use when establishing the VPN tunnel with the peer VPN gateway.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "peer_external_gateway",
				Description: "The URL of the peer side external VPN gateway to which this VPN tunnel is connected.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "peer_external_gateway_interface",
				Description: "The interface ID of the external VPN gateway to which this VPN tunnel is connected.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "peer_gcp_gateway",
				Description: "The URL of the peer side HA GCP VPN gateway to which this VPN tunnel is connected.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "peer_ip",
				Description: "Specifies the IP address of the peer VPN gateway.",
				Type:        proto.ColumnType_IPADDR,
			},
			{
				Name:        "region",
				Description: "The URL of the region where the VPN tunnel resides.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "router",
				Description: "The URL of the router resource to be used for dynamic routing.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "The server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "shared_secret",
				Description: "Specifies the shared secret, used to set the secure session between the Cloud VPN gateway and the peer VPN gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "shared_secret_hash",
				Description: "Specifies the hash of the shared secret.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "target_vpn_gateway",
				Description: "The URL of the Target VPN gateway with which this VPN tunnel is associated.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpn_gateway",
				Description: "The URL of the VPN gateway with which this VPN tunnel is associated.",
				Type:        proto.ColumnType_STRING,
			},
			// simplified view of the vpn gateway, without the full path
			{
				Name:        "vpn_gateway_name",
				Description: "The URL of the VPN gateway with which this VPN tunnel is associated.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VpnGateway").Transform(lastPathElement),
			},
			{
				Name:        "vpn_gateway_interface",
				Description: "The interface ID of the VPN gateway with which this VPN tunnel is associated",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "local_traffic_selector",
				Description: "A list of local traffic selector to use when establishing the VPN tunnel with the peer VPN gateway. The value should be a CIDR formatted string.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "remote_traffic_selector",
				Description: "A list of remote traffic selector to use when establishing the VPN tunnel with the peer VPN gateway. The value should be a CIDR formatted string.",
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
				Hydrate:     getVpnTunnelAka,
				Transform:   transform.FromValue(),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Region").Transform(lastPathElement),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Hydrate:     getProject,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeVpnTunnels(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeVpnTunnels")
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

	resp := service.VpnTunnels.AggregatedList(project)
	if err := resp.Pages(ctx, func(page *compute.VpnTunnelAggregatedList) error {
		for _, item := range page.Items {
			for _, vpnTunnel := range item.VpnTunnels {
				d.StreamListItem(ctx, vpnTunnel)
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getComputeVpnTunnel(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

	var vpnTunnel compute.VpnTunnel
	name := d.KeyColumnQuals["name"].GetStringValue()

	resp := service.VpnTunnels.AggregatedList(project).Filter("name=" + name)
	if err := resp.Pages(
		ctx,
		func(page *compute.VpnTunnelAggregatedList) error {
			for _, item := range page.Items {
				for _, i := range item.VpnTunnels {
					vpnTunnel = *i
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

	// If the specified resource is not present, API does not return any not found errors
	if len(vpnTunnel.Name) < 1 {
		return nil, nil
	}

	return &vpnTunnel, nil
}

func getVpnTunnelAka(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	vpnTunnel := h.Item.(*compute.VpnTunnel)
	region := getLastPathElement(types.SafeString(vpnTunnel.Region))

	// Get project details
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}
	project := projectData.Project

	akas := []string{"gcp://compute.googleapis.com/projects/" + project + "/regions/" + region + "/vpnTunnels/" + vpnTunnel.Name}

	return akas, nil
}
