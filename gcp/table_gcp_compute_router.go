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

func tableGcpComputeRouter(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_router",
		Description: "GCP Compute Router",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeRouter,
			Tags:       map[string]string{"service": "compute", "action": "routers.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeRouters,
			Tags:    map[string]string{"service": "compute", "action": "routers.list"},
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
				Name:        "description",
				Description: "A user-specified, human-readable description of the router.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_timestamp",
				Description: "The creation timestamp of the resource.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "network",
				Description: "The URI of the network to which this router belongs.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "The server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "bgp_advertise_mode",
				Description: "An user-specified flag to indicate which mode to use for advertisement.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Bgp.AdvertiseMode"),
			},
			{
				Name:        "bgp_asn",
				Description: "Specifies the local BGP Autonomous System Number (ASN).",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Bgp.Asn"),
			},
			{
				Name:        "bgp_advertised_groups",
				Description: "An user-specified list of prefix groups to advertise in custom mode.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Bgp.AdvertisedGroups"),
			},
			{
				Name:        "bgp_advertised_ip_ranges",
				Description: "User-specified list of individual IP ranges to advertise in custom mode. This field can only be populated if advertise_mode is CUSTOM and is advertised to all peers of the router. These IP ranges will be advertised in addition to any specified groups.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Bgp.AdvertisedIpRanges"),
			},
			{
				Name:        "bgp_peers",
				Description: "BGP information that must be configured into the routing stack to establish BGP peering. This information must specify the peer ASN and either the interface name, IP address, or peer IP address.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "interfaces",
				Description: "An list of router interfaces.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "nats",
				Description: "A list of NAT services created in this router.",
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
				Transform:   transform.FromP(gcpComputeRouterTurbotData, "Akas"),
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
				Transform:   transform.FromP(gcpComputeRouterTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeRouters(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeRouters")

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Max limit is set as per documentation
	// https://pkg.go.dev/google.golang.org/api@v0.48.0/compute/v1?utm_source=gopls#RoutersAggregatedListCall.MaxResults
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

	resp := service.Routers.AggregatedList(project).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *compute.RouterAggregatedList) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, item := range page.Items {
			for _, router := range item.Routers {
				d.StreamListItem(ctx, router)

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

func getComputeRouter(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

	var router compute.Router
	name := d.EqualsQuals["name"].GetStringValue()

	resp := service.Routers.AggregatedList(project).Filter("name=" + name)
	if err := resp.Pages(
		ctx,
		func(page *compute.RouterAggregatedList) error {
			for _, item := range page.Items {
				for _, i := range item.Routers {
					router = *i
				}
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	// If the specified resource is not present, API does not return any not found errors
	// it returns the data as {<nil> []   0 []   []    {0 map[]} [] []}
	if len(router.Name) < 1 {
		return nil, nil
	}

	return &router, nil
}

//// TRANSFORM FUNCTIONS

func gcpComputeRouterTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	router := d.HydrateItem.(*compute.Router)
	param := d.Param.(string)

	region := getLastPathElement(types.SafeString(router.Region))
	project := strings.Split(router.SelfLink, "/")[6]

	turbotData := map[string]interface{}{
		"Project": project,
		"Akas":    []string{"gcp://compute.googleapis.com/projects/" + project + "/regions/" + region + "/routers/" + router.Name},
	}

	return turbotData[param], nil
}
