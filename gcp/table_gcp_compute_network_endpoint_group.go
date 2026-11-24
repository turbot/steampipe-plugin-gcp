package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
	"google.golang.org/api/compute/v1"
)

//// TABLE DEFINITION

func tableGcpComputeNetworkEndpointGroup(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_network_endpoint_group",
		Description: "GCP Compute Network Endpoint Group",
		List: &plugin.ListConfig{
			Hydrate: listComputeNetworkEndpointGroups,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "name", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "network_endpoint_type", Require: plugin.Optional, Operators: []string{"<>", "="}},
			},
			Tags: map[string]string{"service": "compute", "action": "networkEndpointGroups.list"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "Name of the resource. Provided by the client when the resource is created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique identifier for the resource. This identifier is defined by the server.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "creation_timestamp",
				Description: "Creation timestamp in RFC3339 text format.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "An optional description of this resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "Type of the resource. Always compute#networkEndpointGroup for network endpoint group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "network_endpoint_type",
				Description: "Type of network endpoints in this network endpoint group, e.g. GCE_VM_IP_PORT, GCE_VM_IP, NON_GCP_PRIVATE_IP_PORT, INTERNET_FQDN_PORT, INTERNET_IP_PORT, SERVERLESS or PRIVATE_SERVICE_CONNECT.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "size",
				Description: "Number of network endpoints in the network endpoint group.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "default_port",
				Description: "The default port used if the port number is not specified in the network endpoint.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "network",
				Description: "The URL of the network to which all network endpoints in the NEG belong.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "subnetwork",
				Description: "Optional URL of the subnetwork to which all network endpoints in the NEG belong.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "zone",
				Description: "The URL of the zone where the network endpoint group is located. Set only for zonal network endpoint groups.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "region",
				Description: "The URL of the region where the network endpoint group is located. Set only for regional network endpoint groups.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "psc_target_service",
				Description: "The target service url used to set up private service connection to a Google API or a PSC Producer Service Attachment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "annotations",
				Description: "Metadata defined as annotations on the network endpoint group.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "app_engine",
				Description: "Only valid when networkEndpointType is SERVERLESS. Only one of cloudRun, appEngine or cloudFunction may be set.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "cloud_function",
				Description: "Only valid when networkEndpointType is SERVERLESS. Only one of cloudRun, appEngine or cloudFunction may be set.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "cloud_run",
				Description: "Only valid when networkEndpointType is SERVERLESS. Only one of cloudRun, appEngine or cloudFunction may be set.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "psc_data",
				Description: "Private Service Connect data. Only valid when networkEndpointType is PRIVATE_SERVICE_CONNECT.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "location_type",
				Description: "Location type where the network endpoint group resides. Possible values are ZONAL, REGIONAL and GLOBAL.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpComputeNetworkEndpointGroupLocation, "Type"),
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
				Transform:   transform.From(gcpComputeNetworkEndpointGroupAka),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpComputeNetworkEndpointGroupLocation, "Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpComputeNetworkEndpointGroupLocation, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

// listComputeNetworkEndpointGroups lists zonal NEGs with networkEndpointGroups.aggregatedList,
// global NEGs with globalNetworkEndpointGroups.list and regional NEGs with
// regionNetworkEndpointGroups.list in every region, since regional NEGs have no aggregated list.
// Rows are de-duplicated by self link in case the aggregated list also returns them.
func listComputeNetworkEndpointGroups(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_compute_network_endpoint_group.listComputeNetworkEndpointGroups", "service_error", err)
		return nil, err
	}

	filterQuals := []filterQualMap{
		{"name", "name", "string"},
		{"network_endpoint_type", "networkEndpointType", "string"},
	}

	filters := buildQueryFilterFromQuals(filterQuals, d.Quals)
	filterString := ""
	if len(filters) > 0 {
		filterString = strings.Join(filters, " ")
	}

	// Max limit is set as per documentation
	// https://pkg.go.dev/google.golang.org/api/compute/v1#NetworkEndpointGroupsAggregatedListCall.MaxResults
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

	seen := map[string]bool{}
	done := false
	stream := func(items []*compute.NetworkEndpointGroup) {
		for _, item := range items {
			if done || seen[item.SelfLink] {
				continue
			}
			seen[item.SelfLink] = true
			d.StreamListItem(ctx, item)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				done = true
			}
		}
	}

	// Zonal (and any other scope returned by the aggregated list)
	aggResp := service.NetworkEndpointGroups.AggregatedList(project).Filter(filterString).MaxResults(*pageSize)
	if err := aggResp.Pages(ctx, func(page *compute.NetworkEndpointGroupAggregatedList) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, item := range page.Items {
			stream(item.NetworkEndpointGroups)
		}
		if done {
			page.NextPageToken = ""
		}
		return nil
	}); err != nil {
		plugin.Logger(ctx).Error("gcp_compute_network_endpoint_group.listComputeNetworkEndpointGroups", "aggregated_list_error", err)
		return nil, err
	}
	if done {
		return nil, nil
	}

	// Global
	globalResp := service.GlobalNetworkEndpointGroups.List(project).Filter(filterString).MaxResults(*pageSize)
	if err := globalResp.Pages(ctx, func(page *compute.NetworkEndpointGroupList) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		stream(page.Items)
		if done {
			page.NextPageToken = ""
		}
		return nil
	}); err != nil {
		plugin.Logger(ctx).Error("gcp_compute_network_endpoint_group.listComputeNetworkEndpointGroups", "global_list_error", err)
		return nil, err
	}
	if done {
		return nil, nil
	}

	// Regional
	regions, err := service.Regions.List(project).Do()
	if err != nil {
		plugin.Logger(ctx).Error("gcp_compute_network_endpoint_group.listComputeNetworkEndpointGroups", "region_list_error", err)
		return nil, err
	}
	for _, region := range regions.Items {
		regionResp := service.RegionNetworkEndpointGroups.List(project, region.Name).Filter(filterString).MaxResults(*pageSize)
		if err := regionResp.Pages(ctx, func(page *compute.NetworkEndpointGroupList) error {
			// apply rate limiting
			d.WaitForListRateLimit(ctx)

			stream(page.Items)
			if done {
				page.NextPageToken = ""
			}
			return nil
		}); err != nil {
			plugin.Logger(ctx).Error("gcp_compute_network_endpoint_group.listComputeNetworkEndpointGroups", "region_list_error", err, "region", region.Name)
			return nil, err
		}
		if done {
			return nil, nil
		}
	}

	return nil, nil
}

//// TRANSFORM FUNCTIONS

func gcpComputeNetworkEndpointGroupAka(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*compute.NetworkEndpointGroup)
	project := strings.Split(data.SelfLink, "/")[6]
	zone := getLastPathElement(types.SafeString(data.Zone))
	region := getLastPathElement(types.SafeString(data.Region))

	switch {
	case zone != "":
		return []string{"gcp://compute.googleapis.com/projects/" + project + "/zones/" + zone + "/networkEndpointGroups/" + data.Name}, nil
	case region != "":
		return []string{"gcp://compute.googleapis.com/projects/" + project + "/regions/" + region + "/networkEndpointGroups/" + data.Name}, nil
	}
	return []string{"gcp://compute.googleapis.com/projects/" + project + "/global/networkEndpointGroups/" + data.Name}, nil
}

func gcpComputeNetworkEndpointGroupLocation(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*compute.NetworkEndpointGroup)
	param := d.Param.(string)

	zone := getLastPathElement(types.SafeString(data.Zone))
	region := getLastPathElement(types.SafeString(data.Region))
	project := strings.Split(data.SelfLink, "/")[6]

	locationData := map[string]string{
		"Type":     "GLOBAL",
		"Location": "global",
		"Project":  project,
	}

	switch {
	case zone != "":
		locationData["Type"] = "ZONAL"
		locationData["Location"] = zone
	case region != "":
		locationData["Type"] = "REGIONAL"
		locationData["Location"] = region
	}

	return locationData[param], nil
}
