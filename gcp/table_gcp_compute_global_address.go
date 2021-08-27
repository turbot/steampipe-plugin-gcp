package gcp

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"google.golang.org/api/compute/v1"
)

//// TABLE DEFINITION

func tableGcpComputeGlobalAddress(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_global_address",
		Description: "GCP Compute Global Address",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeGlobalAddress,
		},
		List: &plugin.ListConfig{
			Hydrate:           listComputeGlobalAddresses,
			ShouldIgnoreError: isIgnorableError([]string{"403"}),
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
				Transform:   transform.FromP(globalAddressSelfLinkToTurbotData, "Akas"),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromConstant("global"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(globalAddressSelfLinkToTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeGlobalAddresses(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
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

	resp := service.GlobalAddresses.List(project)
	if err := resp.Pages(ctx, func(page *compute.AddressList) error {
		for _, globalAddress := range page.Items {
			d.StreamListItem(ctx, globalAddress)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getComputeGlobalAddress(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

	name := d.KeyColumnQuals["name"].GetStringValue()

	// Error: pq: rpc error: code = Unknown desc = json: invalid use of ,string struct tag,
	// trying to unmarshal "projects/project/global/addresses/" into uint64
	if len(name) < 1 {
		return nil, nil
	}

	req, err := service.GlobalAddresses.Get(project, name).Do()
	if err != nil {
		return nil, err
	}

	return req, nil
}

//// TRANSFORM FUNCTIONS

func globalAddressSelfLinkToTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	globalAddress := d.HydrateItem.(*compute.Address)
	param := d.Param.(string)

	project := strings.Split(globalAddress.SelfLink, "/")[6]

	turbotData := map[string]interface{}{
		"Project": project,
		"Akas":    []string{"gcp://compute.googleapis.com/projects/" + project + "/global/addresses/" + globalAddress.Name},
	}

	return turbotData[param], nil
}
