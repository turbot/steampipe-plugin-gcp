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

func tableGcpComputeAddress(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_address",
		Description: "GCP Compute Address",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeAddress,
			Tags:       map[string]string{"service": "compute", "action": "addresses.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeAddresses,
			KeyColumns: plugin.KeyColumnSlice{
				// String columns
				{Name: "address_type", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "network_tier", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "purpose", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "status", Require: plugin.Optional, Operators: []string{"<>", "="}},
			},
			Tags: map[string]string{"service": "compute", "action": "addresses.list"},
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
				Transform:   transform.FromGo().NullIfZero(),
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
				Transform:   transform.FromP(addressSelfLinkToTurbotData, "Akas"),
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
				Transform:   transform.FromP(addressSelfLinkToTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeAddresses(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeAddresses")

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	filterQuals := []filterQualMap{
		{"address_type", "addressType", "string"},
		{"network_tier", "networkTier", "string"},
		{"purpose", "purpose", "string"},
		{"status", "status", "string"},
	}

	filters := buildQueryFilterFromQuals(filterQuals, d.Quals)
	filterString := ""
	if len(filters) > 0 {
		filterString = strings.Join(filters, " ")
	}

	// Max limit is set as per documentation
	// https://pkg.go.dev/google.golang.org/api@v0.48.0/compute/v1?utm_source=gopls#AddressesAggregatedListCall.MaxResults
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

	resp := service.Addresses.AggregatedList(project).Filter(filterString).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *compute.AddressAggregatedList) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, item := range page.Items {
			for _, address := range item.Addresses {
				d.StreamListItem(ctx, address)

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

func getComputeAddress(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

	var address compute.Address
	name := d.EqualsQuals["name"].GetStringValue()

	// Empty check
	if name == "" {
		return nil, nil
	}

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

	// If the specified resource is not present, API does not return any not found errors
	if len(address.Name) < 1 {
		return nil, nil
	}

	return &address, nil
}

//// TRANSFORM FUNCTIONS

func addressSelfLinkToTurbotData(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	address := d.HydrateItem.(*compute.Address)

	param := d.Param.(string)
	region := getLastPathElement(types.SafeString(address.Region))
	project := strings.Split(address.SelfLink, "/")[6]

	turbotData := map[string]interface{}{
		"Project": project,
		"Akas":    []string{"gcp://compute.googleapis.com/projects/" + project + "/regions/" + region + "/addresses/" + address.Name},
	}

	return turbotData[param], nil
}
