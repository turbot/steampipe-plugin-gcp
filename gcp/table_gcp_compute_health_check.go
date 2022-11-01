package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
	"google.golang.org/api/compute/v1"
)

func tableGcpComputeHealthCheck(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_health_check",
		Description: "GCP Compute Health Check",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeHealthCheck,
		},
		List: &plugin.ListConfig{
			Hydrate:           listComputeHealthCheck,
			ShouldIgnoreError: isIgnorableError([]string{"403"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "Name of the health check.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique identifier for this health check. This identifier is defined by the server.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "check_interval_sec",
				Description: "How often (in seconds) to send a health check.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "creation_timestamp",
				Description: "Timestamp when the health check was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "An optional description of this resource. Provide this property when you create the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "healthy_threshold",
				Description: "A so-far unhealthy instance will be marked healthy after this many consecutive successes.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "Type of the resource. Always compute#healthCheck for health checks.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "region",
				Description: "The URL of the region where the health check resides. Only applicable for regional resources.",
				Type:        proto.ColumnType_STRING,
			},

			// region_name is a simpler view of the region, without the full path
			{
				Name:        "region_name",
				Description: "Name of the region where the health check resides. Only applicable for regional resources.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Region").Transform(lastPathElement),
			},
			{
				Name:        "self_link",
				Description: "Server-defined fully-qualified URL for this resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "timeout_sec",
				Description: "How long (in seconds) to wait before claiming failure.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "type",
				Description: "Specifies the type of the healthCheck, either TCP, SSL, HTTP, HTTPS or HTTP2.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "unhealthy_threshold",
				Description: "A so-far healthy instance will be marked unhealthy after this many consecutive failures.",
				Type:        proto.ColumnType_INT,
			},

			// JSON columns
			{
				Name:        "grpc_health_check",
				Description: "The details for the gRPC type health check.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "http2_health_check",
				Description: "The details for the HTTP2 type health check.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Http2HealthCheck"),
			},
			{
				Name:        "http_health_check",
				Description: "The details for the HTTP type health check.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "https_health_check",
				Description: "The details for the HTTPS type health check.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "log_config",
				Description: "The logging configuration details on this health check.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "ssl_health_check",
				Description: "The details for the SSL type health check.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tcp_health_check",
				Description: "The details for the TCP type health check.",
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
				Transform:   transform.From(healthCheckAka),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(healthCheckLocation, "Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(healthCheckLocation, "Project"),
			},
		},
	}
}

//// LIST FUNCTIONS

func listComputeHealthCheck(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Max limit is set as per documentation
	pageSize := types.Int64(500)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *pageSize {
			pageSize = limit
		}
	}

	// Get project details
	getProjectCached := plugin.HydrateFunc(getProject).WithCache()
	projectId, err := getProjectCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Debug("gcp_compute_health_check.listComputeHealthCheck", "service_creation_err", err)
		return nil, err
	}

	resp := service.HealthChecks.AggregatedList(project).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *compute.HealthChecksAggregatedList) error {
		for _, item := range page.Items {
			for _, check := range item.HealthChecks {
				d.StreamListItem(ctx, check)

				// Check if context has been cancelled or if the limit has been hit (if specified)
				// if there is a limit, it will return the number of rows required to reach this limit
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					page.NextPageToken = ""
					return nil
				}
			}
		}
		return nil
	}); err != nil {
		plugin.Logger(ctx).Debug("gcp_compute_health_check.listComputeHealthCheck", "api_err", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getComputeHealthCheck(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Get project details
	getProjectCached := plugin.HydrateFunc(getProject).WithCache()
	projectId, err := getProjectCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	var check compute.HealthCheck
	name := d.KeyColumnQuals["name"].GetStringValue()

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Debug("gcp_compute_health_check.getComputeHealthCheck", "service_creation_err", err)
		return nil, err
	}

	resp := service.HealthChecks.AggregatedList(project).Filter("name=" + name)
	if err := resp.Pages(ctx, func(page *compute.HealthChecksAggregatedList) error {
		for _, item := range page.Items {
			for _, i := range item.HealthChecks {
				check = *i
			}
		}
		return nil
	},
	); err != nil {
		plugin.Logger(ctx).Debug("gcp_compute_health_check.getComputeHealthCheck", "api_err", err)
		return nil, err
	}

	// If the specified resource is not present, API does not return any not found errors
	if len(check.Name) < 1 {
		return nil, nil
	}

	return &check, nil
}

//// TRANSFORM FUNCTIONS

func healthCheckAka(_ context.Context, d *transform.TransformData) (interface{}, error) {
	i := d.HydrateItem.(*compute.HealthCheck)

	regionName := getLastPathElement(types.SafeString(i.Region))
	project := strings.Split(i.SelfLink, "/")[6]
	healthCheckName := types.SafeString(i.Name)

	akas := []string{"gcp://compute.googleapis.com/projects/" + project + "/regions/" + regionName + "/healthChecks/" + healthCheckName}

	if regionName == "" {
		akas = []string{"gcp://compute.googleapis.com/projects/" + project + "/global/healthChecks/" + healthCheckName}
	}

	return akas, nil
}

func healthCheckLocation(_ context.Context, d *transform.TransformData) (interface{}, error) {
	i := d.HydrateItem.(*compute.HealthCheck)
	param := d.Param.(string)

	regionName := getLastPathElement(types.SafeString(i.Region))
	project := strings.Split(i.SelfLink, "/")[6]

	locationData := map[string]string{
		"Type":     "REGIONAL",
		"Location": regionName,
		"Project":  project,
	}

	if regionName == "" {
		locationData["Type"] = "GLOBAL"
		locationData["Location"] = "global"
	}

	return locationData[param], nil
}
