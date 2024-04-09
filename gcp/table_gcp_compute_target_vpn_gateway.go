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

func tableGcpComputeTargetVpnGateway(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_target_vpn_gateway",
		Description: "GCP Compute Target VPN Gateway",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeTargetVpnGateway,
			Tags:       map[string]string{"service": "compute", "action": "targetVpnGateways.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeTargetVpnGateways,
			KeyColumns: plugin.KeyColumnSlice{
				// String columns
				{Name: "status", Require: plugin.Optional, Operators: []string{"<>", "="}},
			},
			Tags: map[string]string{"service": "compute", "action": "targetVpnGateways.list"},
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
				Transform:   transform.FromP(gcpComputeTargetVpnGatewayTurbotData, "Akas"),
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
				Transform:   transform.FromP(gcpComputeTargetVpnGatewayTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeTargetVpnGateways(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeTargetVpnGateways")
	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	filterString := ""
	if d.EqualsQuals["status"] != nil {
		filterString = "status=" + d.EqualsQuals["status"].GetStringValue()
	}

	// Max limit is set as per documentation
	// https://pkg.go.dev/google.golang.org/api@v0.48.0/compute/v1?utm_source=gopls#TargetVpnGatewaysAggregatedListCall.MaxResults
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

	resp := service.TargetVpnGateways.AggregatedList(project).Filter(filterString).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *compute.TargetVpnGatewayAggregatedList) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, item := range page.Items {
			for _, targetVpnGateway := range item.TargetVpnGateways {
				d.StreamListItem(ctx, targetVpnGateway)

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

func getComputeTargetVpnGateway(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

	var targetVpnGateway compute.TargetVpnGateway
	name := d.EqualsQuals["name"].GetStringValue()
	if name != "" {
		return nil, nil
	}

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

func gcpComputeTargetVpnGatewayTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	targetVpnGateway := d.HydrateItem.(*compute.TargetVpnGateway)
	param := d.Param.(string)

	region := getLastPathElement(types.SafeString(targetVpnGateway.Region))
	project := strings.Split(targetVpnGateway.SelfLink, "/")[6]

	turbotData := map[string]interface{}{
		"Project": project,
		"Akas":    []string{"gcp://compute.googleapis.com/projects/" + project + "/regions/" + region + "/targetVpnGateways/" + targetVpnGateway.Name},
	}

	return turbotData[param], nil
}
