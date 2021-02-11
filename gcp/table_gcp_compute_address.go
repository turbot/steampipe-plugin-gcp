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

func tableGcpComputeAddress(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_address",
		Description: "GCP Compute Address",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeAddress,
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeAddresses,
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
				Name:        "address",
				Description: "The static IP address represented by this resource.",
				Type:        proto.ColumnType_IPADDR,
			},
			{
				Name:        "address_type",
				Description: "The type of address to reserve, either INTERNAL or EXTERNAL.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ip_version",
				Description: "The IP version that will be used by this address.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A user-specified, human-readable description of the address.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_timestamp",
				Description: "The creation timestamp of the resource.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "kind",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "network",
				Description: "The URL of the network in which to reserve the address.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "network_tier",
				Description: "Specifies the networking tier used for configuring this address.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "prefix_length",
				Description: "Specifies the prefix length if the resource represents an IP range.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "purpose",
				Description: "Specifies the purpose of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "The server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the address.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "subnetwork",
				Description: "The URL of the subnetwork in which to reserve the address.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "users",
				Description: "A list of URLs of the resources that are using this address.",
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
				Transform:   transform.From(addressAka),
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
				Transform:   transform.FromConstant(projectName),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeAddresses(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeAddresses")
	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	project := projectName
	resp := service.Addresses.AggregatedList(project)
	if err := resp.Pages(ctx, func(page *compute.AddressAggregatedList) error {
		for _, item := range page.Items {
			for _, address := range item.Addresses {
				d.StreamListItem(ctx, address)
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getComputeAddress(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	var address compute.Address
	name := d.KeyColumnQuals["name"].GetStringValue()
	project := projectName

	resp := service.Addresses.AggregatedList(project).Filter("name=" + name)
	if err := resp.Pages(
		ctx,
		func(page *compute.AddressAggregatedList) error {
			for _, item := range page.Items {
				for _, i := range item.Addresses {
					address = *i
				}
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	return &address, nil
}

//// TRANSFORM FUNCTIONS

func addressAka(_ context.Context, d *transform.TransformData) (interface{}, error) {
	address := d.HydrateItem.(*compute.Address)
	regionName := getLastPathElement(types.SafeString(address.Region))

	akas := []string{"gcp://compute.googleapis.com/projects/" + activeProject() + "/regions/" + regionName + "/addresses/" + address.Name}

	return akas, nil
}
