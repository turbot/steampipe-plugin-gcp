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

func tableGcpComputeRouter(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_router",
		Description: "GCP Compute Router",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeRouter,
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeRouters,
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
				Name:        "asn",
				Description: "Specifies the local BGP Autonomous System Number (ASN).",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Bgp.Asn"),
			},
			{
				Name:        "advertise_mode",
				Description: "An user-specified flag to indicate which mode to use for advertisement.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Bgp.AdvertiseMode"),
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
				Name:        "advertised_groups",
				Description: "An user-specified list of prefix groups to advertise in custom mode.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Bgp.AdvertisedGroups"),
			},
			{
				Name:        "advertised_ip_ranges",
				Description: "An user-specified ist of individual IP ranges to advertise in custom mode.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Bgp.AdvertisedIpRanges"),
			},
			{
				Name:        "bgp_peers",
				Description: "BGP information that must be configured into the routing stack to establish BGP peering.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Bgp.AdvertisedIpRanges"),
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
				Transform:   transform.From(gcpComputeRouterAka),
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

func listComputeRouters(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeRouters")
	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	project := activeProject()
	resp := service.Routers.AggregatedList(project)
	if err := resp.Pages(ctx, func(page *compute.RouterAggregatedList) error {
		for _, item := range page.Items {
			for _, router := range item.Routers {
				d.StreamListItem(ctx, router)
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
	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	var router compute.Router
	name := d.KeyColumnQuals["name"].GetStringValue()
	project := activeProject()

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

func gcpComputeRouterAka(_ context.Context, d *transform.TransformData) (interface{}, error) {
	router := d.HydrateItem.(*compute.Router)
	regionName := getLastPathElement(types.SafeString(router.Region))

	akas := []string{"gcp://compute.googleapis.com/projects/" + activeProject() + "/regions/" + regionName + "/routers/" + router.Name}

	return akas, nil
}
