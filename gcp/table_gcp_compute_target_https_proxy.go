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

func tableGcpComputeTargetHttpsProxy(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_target_https_proxy",
		Description: "GCP Compute Target Https Proxy",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeTargetHttpsProxy,
			Tags:       map[string]string{"service": "compute", "action": "targetHttpProxies.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeTargetHttpsProxies,
			KeyColumns: plugin.KeyColumnSlice{
				// Boolean columns
				{Name: "proxy_bind", Require: plugin.Optional, Operators: []string{"<>", "="}},
			},
			Tags: map[string]string{"service": "compute", "action": "targetHttpProxies.list"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "Name of the resource. Provided by the client when the resource is created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "A server-defined unique identifier for the resource.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "creation_timestamp",
				Description: "Specifies the time when the resource is created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "A user defined description for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "authorization_policy",
				Description: "Specifies an URL referring to a networksecurity.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "Type of the resource. Always compute#targetHttpsProxy for target HTTPS proxies.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "proxy_bind",
				Description: "This field only applies when the forwarding rule that references this target proxy has a loadBalancingScheme set to INTERNAL_SELF_MANAGED.",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
			},
			{
				Name:        "quic_override",
				Description: "Specifies the QUIC override policy for this TargetHttpsProxy resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "region",
				Description: "An URL of the region where the regional TargetHttpsProxy resides.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "server_tls_policy",
				Description: "An URL referring to a networksecurity.ServerTlsPolicy resource that describes how the proxy should authenticate inbound traffic.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ssl_policy",
				Description: "An URL of SslPolicy resource that will be associated with the TargetHttpsProxy resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "url_map",
				Description: "A fully-qualified or valid partial URL to the UrlMap resource that defines the mapping from URL to the BackendService.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ssl_certificates",
				Description: "A list of URLs to SslCertificate resources that are used to authenticate connections between users and the load balancer.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "location_type",
				Description: "Location type where the target https proxy resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpComputeTargetHttpsProxyLocation, "Type"),
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
				Transform:   transform.From(gcpComputeTargetHttpsProxyAka),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpComputeTargetHttpsProxyLocation, "Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpComputeTargetHttpsProxyLocation, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeTargetHttpsProxies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeTargetHttpsProxies")

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	filterQuals := []filterQualMap{
		{"proxy_bind", "proxyBind", "boolean"},
	}

	filters := buildQueryFilterFromQuals(filterQuals, d.Quals)
	filterString := ""
	if len(filters) > 0 {
		filterString = strings.Join(filters, " ")
	}

	// Max limit is set as per documentation
	// https://pkg.go.dev/google.golang.org/api@v0.48.0/compute/v1?utm_source=gopls#TargetHttpsProxiesAggregatedListCall.MaxResults
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

	resp := service.TargetHttpsProxies.AggregatedList(project).Filter(filterString).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *compute.TargetHttpsProxyAggregatedList) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, item := range page.Items {
			for _, targetHttpsProxy := range item.TargetHttpsProxies {
				d.StreamListItem(ctx, targetHttpsProxy)

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

func getComputeTargetHttpsProxy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getComputeTargetHttpsProxy")

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
	name := d.EqualsQuals["name"].GetStringValue()

	var targetHttpsProxy compute.TargetHttpsProxy
	resp := service.TargetHttpsProxies.AggregatedList(project).Filter("name=" + name)
	if err := resp.Pages(ctx, func(page *compute.TargetHttpsProxyAggregatedList) error {
		for _, item := range page.Items {
			for _, i := range item.TargetHttpsProxies {
				targetHttpsProxy = *i
			}
		}
		return nil
	},
	); err != nil {
		return nil, err
	}

	if len(targetHttpsProxy.Name) < 1 {
		return nil, nil
	}

	return &targetHttpsProxy, nil
}

//// TRANSFORM FUNCTIONS

func gcpComputeTargetHttpsProxyAka(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*compute.TargetHttpsProxy)
	region := getLastPathElement(types.SafeString(data.Region))
	project := strings.Split(data.SelfLink, "/")[6]

	akas := []string{"gcp://compute.googleapis.com/projects/" + project + "/regions/" + region + "/targetHttpsProxies/" + data.Name}

	if region == "" {
		akas = []string{"gcp://compute.googleapis.com/projects/" + project + "/global/targetHttpsProxies/" + data.Name}
	}

	return akas, nil
}

func gcpComputeTargetHttpsProxyLocation(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*compute.TargetHttpsProxy)
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
