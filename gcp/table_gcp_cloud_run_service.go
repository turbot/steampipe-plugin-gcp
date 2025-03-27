package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"google.golang.org/api/run/v2"
)

//// TABLE DEFINITION

func tableGcpCloudRunService(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_cloud_run_service",
		Description: "GCP Cloud Run Service",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "location"}),
			Hydrate:    getCloudRunService,
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudRunServices,
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "location",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildCloudRunLocationList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The fully qualified name of this Service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(lastPathElement),
			},
			{
				Name:        "client",
				Description: "Arbitrary identifier for the API client.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "client_version",
				Description: "Arbitrary version identifier for the API client.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "User-provided description of the Service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "The creation timestamp of the resource.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromGo().NullIfZero(),
			},
			{
				Name:        "delete_time",
				Description: "The deletion time.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromGo().NullIfZero(),
			},
			{
				Name:        "creator",
				Description: "Email address of the authenticated creator.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "etag",
				Description: "A system-generated fingerprint for this version of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "expire_time",
				Description: "For a deleted resource, the time after which it will be permamently deleted.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromGo().NullIfZero(),
			},
			{
				Name:        "generation",
				Description: "A number that monotonically increases every time the user modifies the desired state.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "ingress",
				Description: "Provides the ingress settings for this Service. On output, returns the currently observed ingress settings, or INGRESS_TRAFFIC_UNSPECIFIED if no revision is active.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "The server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     cloudRunServiceSelfLink,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "last_modifier",
				Description: "Email address of the last authenticated modifier.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "latest_created_revision",
				Description: "Name of the last created revision. See comments in `reconciling` for additional information on reconciliation process in Cloud Run.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "latest_ready_revision",
				Description: "Name of the latest revision that is serving traffic. See comments in `reconciling` for additional information on reconciliation process in Cloud Run.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "launch_stage",
				Description: "The launch stage as defined by Google Cloud Platform Launch Stages (https://cloud.google.com/terms/launch-stages). Cloud Run supports `ALPHA`, `BETA`, and `GA`. If no value is specified, GA is assumed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "observed_generation",
				Description: "The generation of this Service currently serving traffic.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "reconciling",
				Description: "Returns true if the Service is currently being acted upon by the system to bring it into the desired state.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "satisfies_pzs",
				Description: "Reserved for future use.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "traffic_tags_cleanup_threshold",
				Description: "Override the traffic tag threshold limit. Garbage collection will start cleaning up non-serving tagged traffic targets based on creation item. The default value is 2000.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "uid",
				Description: "Server assigned unique identifier for the trigger.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "update_time",
				Description: "The last-modified time.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromGo().NullIfZero(),
			},
			{
				Name:        "uri",
				Description: "The main URI in which this Service is serving traffic.",
				Type:        proto.ColumnType_STRING,
			},

			// JSON fields
			{
				Name:        "annotations",
				Description: "Unstructured key value map that may be set by external tools to store and arbitrary metadata.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "binary_authorization",
				Description: "Settings for the Binary Authorization feature.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "conditions",
				Description: "The Conditions of all other associated sub-resources.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "custom_audiences",
				Description: "One or more custom audiences that you want this service to support.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "labels",
				Description: "Unstructured key value map that can be used to organize and categorize objects.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "template",
				Description: "The template used to create revisions for this Service.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "terminal_condition",
				Description: "The Condition of this Service, containing its readiness status, and detailed error information in case it did not reach a serving state.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "traffic",
				Description: "Specifies how to distribute traffic over a collection of Revisions belonging to the Service. If traffic is empty or not provided, defaults to 100% traffic to the latest `Ready` Revision.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "traffic_statuses",
				Description: "Detailed status information for corresponding traffic targets.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "iam_policy",
				Description: "An Identity and Access Management (IAM) policy, which specifies access controls for Google Cloud resources.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudRunServiceIamPolicy,
				Transform:   transform.FromValue(),
			},

			// Standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(cloudRunServiceData, "Title"),
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
				Transform:   transform.FromP(cloudRunServiceData, "Akas"),
			},

			// Standard GCP columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(cloudRunServiceData, "Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(cloudRunServiceData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listCloudRunServices(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString("region")

	var location string
	matrixLocation := d.EqualsQualString(matrixKeyLocation)
	// Since, when the service API is disabled, matrixLocation value will be nil
	if matrixLocation != "" {
		location = matrixLocation
	}

	// Minimize API call as per given location
	if region != "" && region != location {
		return nil, nil
	}

	// Create Service Connection
	service, err := CloudRunService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_cloud_run_service.listCloudRunServices", "service_error", err)
		return nil, err
	}

	// Max limit is set as per documentation
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

	input := "projects/" + project + "/locations/" + location

	resp := service.Projects.Locations.Services.List(input).PageSize(*pageSize)
	if err := resp.Pages(ctx, func(page *run.GoogleCloudRunV2ListServicesResponse) error {
		for _, item := range page.Services {
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
		plugin.Logger(ctx).Error("gcp_cloud_run_service.listCloudRunServices", "api_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCloudRunService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := CloudRunService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_cloud_run_service.getCloudRunService", "service_error", err)
		return nil, err
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	serviceName := d.EqualsQuals["name"].GetStringValue()
	location := d.EqualsQuals["location"].GetStringValue()

	// Empty Check
	if serviceName == "" || location == "" {
		return nil, nil
	}

	input := "projects/" + project + "/locations/" + location + "/services/" + serviceName

	resp, err := service.Projects.Locations.Services.Get(input).Do()
	if err != nil {
		plugin.Logger(ctx).Error("gcp_cloud_run_service.getCloudRunService", "api_error", err)
		return nil, err
	}
	return resp, err
}

func getCloudRunServiceIamPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	data := h.Item.(*run.GoogleCloudRunV2Service)
	serviceName := strings.Split(data.Name, "/")[5]
	location := strings.Split(data.Name, "/")[3]

	// Create Service Connection
	service, err := CloudRunService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_cloud_run_service.getCloudRunServiceIamPolicy", "service_error", err)
		return nil, err
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	input := "projects/" + project + "/locations/" + location + "/services/" + serviceName

	resp, err := service.Projects.Locations.Services.GetIamPolicy(input).Do()
	if err != nil {
		plugin.Logger(ctx).Error("gcp_cloud_run_service.getCloudRunServiceIamPolicy", "api_error", err)
		return nil, err
	}

	return resp, err
}

//// TRANSFORM FUNCTIONS

func cloudRunServiceSelfLink(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(*run.GoogleCloudRunV2Service)

	var location string
	matrixLocation := d.EqualsQualString(matrixKeyLocation)
	// Since, when the service API is disabled, matrixLocation value will be nil
	if matrixLocation != "" {
		location = matrixLocation
	}

	projectID := strings.Split(data.Name, "/")[1]
	name := strings.Split(data.Name, "/")[5]

	selfLink := "https://run.googleapis.com/v2/projects/" + projectID + "/locations/" + location + "/services/" + name

	return selfLink, nil
}

//// TRANSFORM FUNCTIONS

func cloudRunServiceData(ctx context.Context, h *transform.TransformData) (interface{}, error) {
	data := h.HydrateItem.(*run.GoogleCloudRunV2Service)
	param := h.Param.(string)

	projectID := strings.Split(data.Name, "/")[1]
	name := strings.Split(data.Name, "/")[5]
	location := strings.Split(data.Name, "/")[3]

	turbotData := map[string]interface{}{
		"Project":  projectID,
		"Title":    name,
		"Location": location,
		"Akas":     []string{"gcp://run.googleapis.com/projects/" + projectID + "/services/" + name},
	}

	return turbotData[param], nil
}
