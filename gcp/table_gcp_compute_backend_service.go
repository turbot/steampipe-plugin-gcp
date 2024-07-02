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

func tableGcpComputeBackendService(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_backend_service",
		Description: "GCP Compute Backend Service",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeBackendService,
			Tags:       map[string]string{"service": "compute", "action": "backendServices.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeBackendServices,
			KeyColumns: plugin.KeyColumnSlice{
				// String columns
				{Name: "load_balancing_scheme", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "port_name", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "session_affinity", Require: plugin.Optional, Operators: []string{"<>", "="}},

				// Boolean columns
				{Name: "enable_cdn", Require: plugin.Optional, Operators: []string{"<>", "="}},
			},
			Tags: map[string]string{"service": "compute", "action": "backendServices.list"},
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
				Name:        "enable_cdn",
				Description: "Specifies whether the Cloud CDN is enabled for the backend service, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("EnableCDN"),
			},
			{
				Name:        "load_balancing_scheme",
				Description: "Specifies the type of the load balancer.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A user-specified, human-readable description of the backend service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_timestamp",
				Description: "The creation timestamp of the resource.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "affinity_cookie_ttl_sec",
				Description: "Specifies the lifetime of the cookies in seconds. Only applicable if the loadBalancingScheme is EXTERNAL, INTERNAL_SELF_MANAGED, or INTERNAL_MANAGED, the protocol is HTTP or HTTPS, and the sessionAffinity is GENERATED_COOKIE, or HTTP_COOKIE.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "connection_draining_timeout_sec",
				Description: "Specifies the amount of time in seconds to allow existing connections to persist while on unhealthy backend VMs. Only applicable if the protocol is not UDP. The valid range is [0, 3600].",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ConnectionDraining.DrainingTimeoutSec"),
			},
			{
				Name:        "fingerprint",
				Description: "An unique system generated string, to reduce conflicts when multiple users change any property of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "locality_lb_policy",
				Description: "Specifies the load balancing algorithm used within the scope of the locality.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "log_config_enable",
				Description: "Specifies whether to enable logging for the load balancer traffic served by this backend service, or not.",
				Type:        proto.ColumnType_BOOL,
				// Default:     false,
				Transform: transform.FromField("LogConfig.Enable"),
			},
			{
				Name:        "log_config_sample_rate",
				Description: "Specifies the sampling rate of requests to the load balancer where 1.0 means all logged requests are reported and 0.0 means no logged requests are reported. The default value is 1.0.",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("LogConfig.SampleRate"),
			},
			{
				Name:        "network",
				Description: "The URL of the network to which this backend service belongs.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "port",
				Description: "Specifies the TCP port to connect on the backend. The default value is 80.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "port_name",
				Description: "A named port on a backend instance group representing the port for communication to the backend VMs in that group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "protocol",
				Description: "Specifies the protocol that the BackendService uses to communicate with backends.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "region",
				Description: "The URL of the region where the regional backend service resides. This field is not applicable to global backend services.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "security_policy",
				Description: "The resource URL for the security policy associated with this backend service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "The server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "session_affinity",
				Description: "Specifies the type of session affinity to use. The default is NONE. Session affinity is not applicable if the protocol is UDP.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "signed_url_cache_max_age_sec",
				Description: "Specifies the maximum number of seconds the response to a signed URL request will be considered fresh.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("CdnPolicy.SignedUrlCacheMaxAgeSec"),
			},
			{
				Name:        "timeout_sec",
				Description: "Specifies the backend service timeout.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "backends",
				Description: "An list of backends that serve this BackendService.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "cdn_policy_cache_key_policy",
				Description: "Specifies the CacheKeyPolicy for this CdnPolicy.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CdnPolicy.CacheKeyPolicy"),
			},
			{
				Name:        "circuit_breakers",
				Description: "Settings controlling the volume of connections to a backend service.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "failover_policy",
				Description: "Applicable only to Failover for Internal TCP/UDP Load Balancing.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "health_checks",
				Description: "A list of URLs to the healthChecks, httpHealthChecks (legacy), or httpsHealthChecks (legacy) resource for health checking this backend service.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "iap",
				Description: "Specifies the configurations for Identity-Aware Proxy on this resource.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "security_settings",
				Description: "Specifies the security policy that applies to this backend service.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "signed_url_key_names",
				Description: "A list of names of the keys for signing request URLs.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CdnPolicy.SignedUrlKeyNames"),
			},
			{
				Name:        "location_type",
				Description: "Location type where the backend service resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpComputeBackendServiceLocation, "Type"),
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
				Transform:   transform.From(gcpComputeBackendServiceAka),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpComputeBackendServiceLocation, "Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpComputeBackendServiceLocation, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeBackendServices(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeBackendServices")

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	filterQuals := []filterQualMap{
		{"load_balancing_scheme", "loadBalancingScheme", "string"},
		{"port_name", "portName", "string"},
		{"session_affinity", "sessionAffinity", "string"},
		{"enable_cdn", "enableCdn", "boolean"},
	}

	filters := buildQueryFilterFromQuals(filterQuals, d.Quals)
	filterString := ""
	if len(filters) > 0 {
		filterString = strings.Join(filters, " ")
	}

	// Max limit is set as per documentation
	// https://pkg.go.dev/google.golang.org/api@v0.48.0/compute/v1?utm_source=gopls#BackendServicesAggregatedListCall.MaxResults
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

	resp := service.BackendServices.AggregatedList(project).Filter(filterString).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *compute.BackendServiceAggregatedList) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, item := range page.Items {
			for _, backendService := range item.BackendServices {
				d.StreamListItem(ctx, backendService)

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

func getComputeBackendService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

	var backendService compute.BackendService
	name := d.EqualsQuals["name"].GetStringValue()

	resp := service.BackendServices.AggregatedList(project).Filter("name=" + name)
	if err := resp.Pages(
		ctx,
		func(page *compute.BackendServiceAggregatedList) error {
			for _, item := range page.Items {
				for _, i := range item.BackendServices {
					backendService = *i
				}
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	// If the specified resource is not present, API does not return any not found errors
	if len(backendService.Name) < 1 {
		return nil, nil
	}

	return &backendService, nil
}

//// TRANSFORM FUNCTIONS

func gcpComputeBackendServiceAka(_ context.Context, d *transform.TransformData) (interface{}, error) {
	backendService := d.HydrateItem.(*compute.BackendService)
	region := getLastPathElement(types.SafeString(backendService.Region))
	project := strings.Split(backendService.SelfLink, "/")[6]

	akas := []string{"gcp://compute.googleapis.com/projects/" + project + "/regions/" + region + "/backendServices/" + backendService.Name}

	if region == "" {
		akas = []string{"gcp://compute.googleapis.com/projects/" + project + "/global/backendServices/" + backendService.Name}
	}

	return akas, nil
}

func gcpComputeBackendServiceLocation(_ context.Context, d *transform.TransformData) (interface{}, error) {
	backendService := d.HydrateItem.(*compute.BackendService)
	param := d.Param.(string)
	regionName := getLastPathElement(types.SafeString(backendService.Region))
	project := strings.Split(backendService.SelfLink, "/")[6]

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
