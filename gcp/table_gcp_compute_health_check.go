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

func tableGcpComputeHealthCheck(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_health_check",
		Description: "GCP Compute Health Check",
		List: &plugin.ListConfig{
			Hydrate: listComputeHealthChecks,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "name", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "type", Require: plugin.Optional, Operators: []string{"<>", "="}},
			},
			Tags: map[string]string{"service": "compute", "action": "healthChecks.list"},
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
				Description: "Type of the resource. Always compute#healthCheck for health checks.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "Specifies the type of the healthCheck, either TCP, SSL, HTTP, HTTPS, HTTP2 or GRPC.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "check_interval_sec",
				Description: "How often (in seconds) to send a health check.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "timeout_sec",
				Description: "How long (in seconds) to wait before claiming failure.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "healthy_threshold",
				Description: "A so-far unhealthy instance will be marked healthy after this many consecutive successes.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "unhealthy_threshold",
				Description: "A so-far healthy instance will be marked unhealthy after this many consecutive failures.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "region",
				Description: "URL of the region where the health check resides. This field applies only to the regional resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "grpc_health_check",
				Description: "The gRPC health check configuration.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "grpc_tls_health_check",
				Description: "The gRPC with TLS health check configuration.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "http_health_check",
				Description: "The HTTP health check configuration.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "http2_health_check",
				Description: "The HTTP/2 health check configuration.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "https_health_check",
				Description: "The HTTPS health check configuration.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "ssl_health_check",
				Description: "The SSL health check configuration.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tcp_health_check",
				Description: "The TCP health check configuration.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "log_config",
				Description: "Configures logging on this health check.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "source_regions",
				Description: "The list of cloud regions from which health checks are performed. Only applies to global health checks used by DNS routing policies.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "location_type",
				Description: "Location type where the health check resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpComputeHealthCheckLocation, "Type"),
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
				Transform:   transform.From(gcpComputeHealthCheckAka),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpComputeHealthCheckLocation, "Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpComputeHealthCheckLocation, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeHealthChecks(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeHealthChecks")

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	filterQuals := []filterQualMap{
		{"name", "name", "string"},
		{"type", "type", "string"},
	}

	filters := buildQueryFilterFromQuals(filterQuals, d.Quals)
	filterString := ""
	if len(filters) > 0 {
		filterString = strings.Join(filters, " ")
	}

	// Max limit is set as per documentation
	// https://pkg.go.dev/google.golang.org/api/compute/v1#HealthChecksAggregatedListCall.MaxResults
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

	resp := service.HealthChecks.AggregatedList(project).Filter(filterString).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *compute.HealthChecksAggregatedList) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, item := range page.Items {
			for _, healthCheck := range item.HealthChecks {
				d.StreamListItem(ctx, healthCheck)

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

//// TRANSFORM FUNCTIONS

func gcpComputeHealthCheckAka(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*compute.HealthCheck)
	region := getLastPathElement(types.SafeString(data.Region))
	project := strings.Split(data.SelfLink, "/")[6]

	akas := []string{"gcp://compute.googleapis.com/projects/" + project + "/regions/" + region + "/healthChecks/" + data.Name}

	if region == "" {
		akas = []string{"gcp://compute.googleapis.com/projects/" + project + "/global/healthChecks/" + data.Name}
	}

	return akas, nil
}

func gcpComputeHealthCheckLocation(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*compute.HealthCheck)
	param := d.Param.(string)

	regionName := getLastPathElement(types.SafeString(data.Region))
	project := strings.Split(data.SelfLink, "/")[6]

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
