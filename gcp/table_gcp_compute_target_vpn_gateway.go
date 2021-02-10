package gcp

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"google.golang.org/api/compute/v1"
)

//// TABLE DEFINITION

func tableGcpComputeTargetVpnGateway(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_target_vpn_gateway",
		Description: "GCP Compute Target VPN Gateway",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeTargetVpnGateway,
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeTargetVpnGateways,
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
				Description: "A user-specified, human-readable description of the target vpn gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "Specifies the status of the VPN gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_timestamp",
				Description: "The creation timestamp of the resource.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "network",
				Description: "The URL of the network to which this VPN gateway is attached.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "region",
				Description: "The URL of the region where the target VPN gateway resides.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "The server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "forwarding_rules",
				Description: "A list of URLs to the ForwardingRule resources.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tunnels",
				Description: "A list of URLs to VpnTunnel resources.",
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
				Transform:   transform.From(gcpComputeTargetVpnGatewayAka),
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
				Transform:   transform.FromConstant(activeProject()),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeTargetVpnGateways(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeTargetVpnGateways")
	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	project := activeProject()
	resp := service.TargetVpnGateways.AggregatedList(project)
	if err := resp.Pages(ctx, func(page *compute.TargetVpnGatewayAggregatedList) error {
		for _, item := range page.Items {
			for _, targetVpnGateway := range item.TargetVpnGateways {
				d.StreamListItem(ctx, targetVpnGateway)
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getComputeTargetVpnGateway(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	var targetVpnGateway compute.TargetVpnGateway
	name := d.KeyColumnQuals["name"].GetStringValue()
	project := activeProject()

	resp := service.TargetVpnGateways.AggregatedList(project).Filter("name=" + name)
	if err := resp.Pages(
		ctx,
		func(page *compute.TargetVpnGatewayAggregatedList) error {
			for _, item := range page.Items {
				for _, i := range item.TargetVpnGateways {
					targetVpnGateway = *i
				}
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	// If the specified resource is not present, API does not return any not found errors
	if len(targetVpnGateway.Name) < 1 {
		return nil, nil
	}

	return &targetVpnGateway, nil
}

//// TRANSFORM FUNCTIONS

func gcpComputeTargetVpnGatewayAka(_ context.Context, d *transform.TransformData) (interface{}, error) {
	targetVpnGateway := d.HydrateItem.(*compute.TargetVpnGateway)
	regionName := getLastPathElement(types.SafeString(targetVpnGateway.Region))

	akas := []string{"gcp://compute.googleapis.com/projects/" + activeProject() + "/regions/" + regionName + "/targetVpnGateways/" + targetVpnGateway.Name}

	return akas, nil
}
