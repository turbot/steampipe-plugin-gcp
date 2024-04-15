package gcp

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"google.golang.org/api/serviceusage/v1"
)

//// TABLE DEFINITION

func tableGcpProjectService(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_project_service",
		Description: "GCP Project Service",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			Hydrate:           getGcpProjectService,
			Tags:              map[string]string{"service": "serviceusage", "action": "services.get"},
			ShouldIgnoreError: isIgnorableError([]string{"404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listGcpProjectServices,
			Tags:    map[string]string{"service": "serviceusage", "action": "services.list"},
			KeyColumns: plugin.KeyColumnSlice{
				// String columns
				{Name: "state", Require: plugin.Optional, Operators: []string{"<>", "="}},
			},
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

//// FETCH FUNCTIONS

func listGcpProjectServices(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := ServiceUsageService(ctx, d)
	if err != nil {
		return nil, err
	}

	filterString := ""
	if d.EqualsQuals["state"] != nil {
		filterString = "state:" + d.EqualsQuals["state"].GetStringValue()
	}

	// Max limit is set as per documentation
	// https://pkg.go.dev/google.golang.org/api@v0.48.0/serviceusage/v1?utm_source=gopls#ServicesListCall.PageSize
	pageSize := types.Int64(200)
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

	result := service.Services.List("projects/" + project).Filter(filterString).PageSize(*pageSize)
	if err := result.Pages(
		ctx,
		func(page *serviceusage.ListServicesResponse) error {
			// apply rate limiting
			d.WaitForListRateLimit(ctx)

			for _, service := range page.Services {
				d.StreamListItem(ctx, service)

				// Check if context has been cancelled or if the limit has been hit (if specified)
				// if there is a limit, it will return the number of rows required to reach this limit
				if d.RowsRemaining(ctx) == 0 {
					page.NextPageToken = ""
					return nil
				}
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

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)
	quals := d.EqualsQuals
	name := quals["name"].GetStringValue()
	op, err := service.Services.Get("projects/" + project + "/services/" + name).Do()
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
