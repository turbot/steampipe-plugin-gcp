package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"google.golang.org/api/serviceusage/v1"
)

//// TABLE DEFINITION

func tableGcpProjectService(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_project_service",
		Description: "GCP Project Service",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			ItemFromKey:       serviceNameFromKey,
			Hydrate:           getGcpProjectService,
			ShouldIgnoreError: isNotFoundError([]string{"404", "403"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listGcpProjectServices,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The resource name of the consumer and service",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(lastPathElement),
			},
			{
				Name:        "state",
				Description: "Specifies the state of the service",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "parent",
				Description: "The resource name of the consumer",
				Type:        proto.ColumnType_STRING,
			},

			// standard steampipe columns
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(projectServiceNameToAkas),
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
				Hydrate:     getProject,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// ITEM FROM KEY

func serviceNameFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	name := quals["name"].GetStringValue()
	item := &serviceusage.GoogleApiServiceusageV1Service{
		Name: name,
	}
	return item, nil
}

//// FETCH FUNCTIONS

func listGcpProjectServices(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := ServiceUsageService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}
	project := projectData.Project

	result := service.Services.List("projects/" + project)
	if err := result.Pages(
		ctx,
		func(page *serviceusage.ListServicesResponse) error {
			for _, service := range page.Services {
				d.StreamListItem(ctx, service)
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getGcpProjectService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getGcpProjectService")

	// Create Service Connection
	service, err := ServiceUsageService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}
	project := projectData.Project

	serviceData := h.Item.(*serviceusage.GoogleApiServiceusageV1Service)
	op, err := service.Services.Get("projects/" + project + "/services/" + serviceData.Name).Do()
	if err != nil {
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func projectServiceNameToAkas(_ context.Context, d *transform.TransformData) (interface{}, error) {
	service := d.HydrateItem.(*serviceusage.GoogleApiServiceusageV1Service)
	akas := []string{"gcp://serviceusage.googleapis.com/" + service.Name}

	return akas, nil
}
