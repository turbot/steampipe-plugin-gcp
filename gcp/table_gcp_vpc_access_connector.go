package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/vpcaccess/v1"
)

//// TABLE DEFINITION

func tableGcpVPCAccessConnector(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_vpc_access_connector",
		Description: "GCP VPC Access Connector",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getVPCAccessConnector,
		},
		List: &plugin.ListConfig{
			Hydrate: listVPCAccessConnectors,
		},
		GetMatrixItemFunc: BuildVPCAccessLocationList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The resource name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "State of the VPC access connector.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "network",
				Description: "Name of a VPC network.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ip_cidr_range",
				Description: "The range of internal addresses that follows RFC 4632 notation.",
				Type:        proto.ColumnType_CIDR,
			},
			{
				Name:        "machine_type",
				Description: "Machine type of VM Instance underlying connector.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "The server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(vpcAccessConnectorTurbotData, "SelfLink"),
			},
			{
				Name:        "max_instances",
				Description: "Maximum value of instances in autoscaling group underlying the connector.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "max_throughput",
				Description: "Maximum throughput of the connector in Mbps.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "min_instances",
				Description: "Minimum throughput of the connector in Mbps.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "min_throughput",
				Description: "Minimum throughput of the connector in Mbps.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "connected_projects",
				Description: "List of projects using the connector.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "subnet",
				Description: "The subnet in which to house the VPC Access Connector.",
				Type:        proto.ColumnType_JSON,
			},

			// standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(extractLastPartSeparatedByBackslash),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(vpcAccessConnectorTurbotData, "Akas"),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(vpcAccessConnectorTurbotData, "Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(vpcAccessConnectorTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listVPCAccessConnectors(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var location string
	matrixLocation := d.EqualsQualString(matrixKeyLocation)
	// Since, when the service API is disabled, matrixLocation value will be nil
	if matrixLocation != "" {
		location = matrixLocation
	}

	// Max limit is set as per documentation
	pageSize := types.Int64(500)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *pageSize {
			pageSize = limit
		}
	}

	// Create service connection
	service, err := VPCAccessService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_vpc_access_connector.listVPCAccessConnectors", "service_error", err)
		return nil, err
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	parent := "projects/" + project + "/locations/" + location

	resp := service.Projects.Locations.Connectors.List(parent).PageSize(*pageSize)
	if err := resp.Pages(ctx, func(page *vpcaccess.ListConnectorsResponse) error {
		for _, item := range page.Connectors {
			d.StreamListItem(ctx, item)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				page.NextPageToken = ""
				return nil
			}
		}
		return nil
	}); err != nil {
		plugin.Logger(ctx).Error("gcp_vpc_access_connector.listVPCAccessConnectors", "api_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getVPCAccessConnector(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	matrixLocation := d.EqualsQualString(matrixKeyLocation)

	name := d.EqualsQualString("name")

	if name == "" {
		return nil, nil
	}

	// Restrict the API call for other locations
	if len(strings.Split(name, "/")) > 2 && strings.Split(name, "/")[3] != matrixLocation {
		return nil, nil
	}

	// Create service connection
	service, err := VPCAccessService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_vpc_access_connector.getVPCAccessConnector", "service_error", err)
		return nil, err
	}

	resp, err := service.Projects.Locations.Connectors.Get(name).Do()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

//// TRANSFORM FUNCTIONS

func vpcAccessConnectorTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*vpcaccess.Connector)
	param := d.Param.(string)

	splitName := strings.Split(data.Name, "/")

	turbotData := map[string]interface{}{
		"Project":  splitName[1],
		"Location": splitName[3],
		"SelfLink": "https://vpcaccess.googleapis.com/v1/" + data.Name,
		"Akas":     []string{"gcp://vpcaccess.googleapis.com/" + data.Name},
	}

	return turbotData[param], nil
}
