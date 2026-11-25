package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
	"google.golang.org/api/networkservices/v1"
)

//// TABLE DEFINITION

func tableGcpNetworkServicesServiceLbPolicy(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_network_services_service_lb_policy",
		Description: "GCP Network Services Service LB Policy",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getNetworkServicesServiceLbPolicy,
			Tags:       map[string]string{"service": "networkservices", "action": "serviceLbPolicies.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listNetworkServicesServiceLbPolicies,
			Tags:    map[string]string{"service": "networkservices", "action": "serviceLbPolicies.list"},
		},
		GetMatrixItemFunc: BuildNetworkServicesLocationList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "Name of the ServiceLbPolicy resource, in the format projects/{project}/locations/{location}/serviceLbPolicies/{service_lb_policy_name}.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A free-text description of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "load_balancing_algorithm",
				Description: "The type of load balancing algorithm to be used, e.g. SPRAY_TO_WORLD, SPRAY_TO_REGION, WATERFALL_BY_REGION or WATERFALL_BY_ZONE.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "The timestamp when this resource was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "update_time",
				Description: "The timestamp when this resource was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(networkServicesServiceLbPolicyTurbotData, "SelfLink"),
			},
			{
				Name:        "auto_capacity_drain",
				Description: "Configuration to automatically move traffic away from unhealthy IG/NEG for the associated backend services.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "failover_config",
				Description: "Configuration related to health based failover.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "isolation_config",
				Description: "Configuration to provide isolation support for the associated backend service.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "labels",
				Description: "Set of label tags associated with the ServiceLbPolicy resource.",
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
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Labels"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(networkServicesServiceLbPolicyTurbotData, "Akas"),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(networkServicesServiceLbPolicyTurbotData, "Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(networkServicesServiceLbPolicyTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listNetworkServicesServiceLbPolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var location string
	matrixLocation := d.EqualsQualString(matrixKeyLocation)
	// When the Network Services API is disabled there are no locations, so this runs once with an empty location and the API call below returns a 403
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
	service, err := NetworkServicesService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_network_services_service_lb_policy.listNetworkServicesServiceLbPolicies", "service_error", err)
		return nil, err
	}

	// Get project details
	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	parent := "projects/" + project + "/locations/" + location

	resp := service.Projects.Locations.ServiceLbPolicies.List(parent).PageSize(*pageSize)
	if err := resp.Pages(ctx, func(page *networkservices.ListServiceLbPoliciesResponse) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, item := range page.ServiceLbPolicies {
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
		plugin.Logger(ctx).Error("gcp_network_services_service_lb_policy.listNetworkServicesServiceLbPolicies", "api_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getNetworkServicesServiceLbPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	matrixLocation := d.EqualsQualString(matrixKeyLocation)

	name := d.EqualsQualString("name")

	if name == "" {
		return nil, nil
	}

	// Restrict the API call for other locations
	if len(strings.Split(name, "/")) > 3 && strings.Split(name, "/")[3] != matrixLocation {
		return nil, nil
	}

	// Create service connection
	service, err := NetworkServicesService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_network_services_service_lb_policy.getNetworkServicesServiceLbPolicy", "service_error", err)
		return nil, err
	}

	resp, err := service.Projects.Locations.ServiceLbPolicies.Get(name).Do()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

//// TRANSFORM FUNCTIONS

func networkServicesServiceLbPolicyTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*networkservices.ServiceLbPolicy)
	param := d.Param.(string)

	splitName := strings.Split(data.Name, "/")

	turbotData := map[string]interface{}{
		"Project":  splitName[1],
		"Location": splitName[3],
		"SelfLink": "https://networkservices.googleapis.com/v1/" + data.Name,
		"Akas":     []string{"gcp://networkservices.googleapis.com/" + data.Name},
	}

	return turbotData[param], nil
}
