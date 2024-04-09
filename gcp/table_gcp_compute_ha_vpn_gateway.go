package gcp

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"google.golang.org/api/compute/v1"
)

//// TABLE DEFINITION

func tableGcpComputeHaVpnGateway(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_ha_vpn_gateway",
		Description: "GCP Compute VPN Gateway",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeHaVpnGateway,
			Tags:       map[string]string{"service": "compute", "action": "vpnGateways.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeHaVpnGateways,
			Tags:    map[string]string{"service": "compute", "action": "vpnGateways.list"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getComputeHaVpnGateway,
				Tags: map[string]string{"service": "compute", "action": "vpnGateways.get"},
			},
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
				Name:        "creation_timestamp",
				Description: "The creation timestamp of the resource.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "A user-specified, human-readable description of the vpn gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "label_fingerprint",
				Description: "To see the latest fingerprint, make a get() request to retrieve an VpnGateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "network",
				Description: "URL of the network to which this VPN gateway is attached.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "region",
				Description: "The URL of the region where the VPN gateway resides.",
				Type:        proto.ColumnType_STRING,
			},

			// region_name is a simpler view of the region, without the full path
			{
				Name:        "region_name",
				Description: "Name of the region where the VPN gateway resides. Only applicable for regional resources.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Region").Transform(lastPathElement),
			},
			{
				Name:        "self_link",
				Description: "The server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpn_connections",
				Description: "List of VPN connection for this VpnGateway.",
				Hydrate:     getComputeHaVpnGatewayVpnConnections,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "labels",
				Description: "Labels for this resource.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "vpn_interfaces",
				Description: "The list of VPN interfaces associated with this VPN gateway.",
				Type:        proto.ColumnType_JSON,
			},

			//  Steampipe standard columns
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
				Hydrate:     getVpnGatewayAka,
				Transform:   transform.FromValue(),
			},

			// GCP standard columns
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

func listComputeHaVpnGateways(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_compute_ha_vpn_gateway.listComputeHaVpnGateways", "service_creation_err", err)
		return nil, err
	}

	resp := service.VpnGateways.AggregatedList(project).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *compute.VpnGatewayAggregatedList) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, item := range page.Items {
			for _, vpnGateway := range item.VpnGateways {
				d.StreamListItem(ctx, vpnGateway)

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
		plugin.Logger(ctx).Error("gcp_compute_ha_vpn_gateway.listComputeHaVpnGateways", "api_err", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getComputeHaVpnGateway(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	var vpnGateway compute.VpnGateway
	name := d.EqualsQuals["name"].GetStringValue()
	// Empty check
	if name != "" {
		return nil, nil
	}

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_compute_ha_vpn_gateway.getComputeHaVpnGateway", "service_creation_err", err)
		return nil, err
	}

	resp := service.VpnGateways.AggregatedList(project).Filter("name=" + name)
	if err := resp.Pages(
		ctx,
		func(page *compute.VpnGatewayAggregatedList) error {
			for _, item := range page.Items {
				for _, i := range item.VpnGateways {
					vpnGateway = *i
				}
			}
			return nil
		},
	); err != nil {
		plugin.Logger(ctx).Error("gcp_compute_ha_vpn_gateway.getComputeHaVpnGateway", "api_err", err)
		return nil, err
	}

	// If the specified resource is not present, API does not return any not found errors
	if len(vpnGateway.Name) < 1 {
		return nil, nil
	}

	return &vpnGateway, nil
}

func getComputeHaVpnGatewayVpnConnections(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	vpnGateway := h.Item.(*compute.VpnGateway)
	region := getLastPathElement(types.SafeString(vpnGateway.Region))

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_compute_ha_vpn_gateway.getComputeHaVpnGatewayStatus", "service_creation_err", err)
		return nil, err
	}

	resp, err := service.VpnGateways.GetStatus(project, region, vpnGateway.Name).Do()
	if err != nil {
		plugin.Logger(ctx).Error("gcp_compute_ha_vpn_gateway.getComputeHaVpnGatewayStatus", "api_err", err)
		return nil, err
	}
	return resp.Result.VpnConnections, nil
}

func getVpnGatewayAka(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	vpnGateway := h.Item.(*compute.VpnGateway)
	region := getLastPathElement(types.SafeString(vpnGateway.Region))

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_compute_ha_vpn_gateway.getVpnGatewayAka", "cache_err", err)
		return nil, err
	}
	project := projectId.(string)

	akas := []string{"gcp://compute.googleapis.com/projects/" + project + "/regions/" + region + "/vpnGateways/" + vpnGateway.Name}

	return akas, nil
}
