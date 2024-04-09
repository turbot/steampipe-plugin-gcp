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

//// TABLE DEFINITION

func tableGcpComputeVpnTunnel(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_vpn_tunnel",
		Description: "GCP Compute VPN Tunnel",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeVpnTunnel,
			Tags:       map[string]string{"service": "compute", "action": "vpnTunnels.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeVpnTunnels,
			KeyColumns: plugin.KeyColumnSlice{
				// String columns
				{Name: "vpn_gateway", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "status", Require: plugin.Optional, Operators: []string{"<>", "="}},
			},
			Tags: map[string]string{"service": "compute", "action": "vpnTunnels.list"},
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

func listComputeVpnTunnels(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeVpnTunnels")
	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	filterQuals := []filterQualMap{
		{"vpn_gateway", "vpnGateway", "string"},
		{"status", "status", "string"},
	}

	filters := buildQueryFilterFromQuals(filterQuals, d.Quals)
	filterString := ""
	if len(filters) > 0 {
		filterString = strings.Join(filters, " ")
	}

	// Max limit is set as per documentation
	// https://pkg.go.dev/google.golang.org/api@v0.48.0/compute/v1?utm_source=gopls#VpnTunnelsAggregatedListCall.MaxResults
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

	resp := service.VpnTunnels.AggregatedList(project).Filter(filterString).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *compute.VpnTunnelAggregatedList) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, item := range page.Items {
			for _, vpnTunnel := range item.VpnTunnels {
				d.StreamListItem(ctx, vpnTunnel)

				// Check if context has been cancelled or if the limit has been hit (if specified)
				// if there is a limit, it will return the number of rows required to reach this limit
				if d.RowsRemaining(ctx) == 0 {
					page.NextPageToken = ""
					return nil
				}
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

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	var vpnTunnel compute.VpnTunnel
	name := d.EqualsQuals["name"].GetStringValue()

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

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	akas := []string{"gcp://compute.googleapis.com/projects/" + project + "/regions/" + region + "/vpnTunnels/" + vpnTunnel.Name}

	return akas, nil
}
