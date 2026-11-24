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

func tableGcpComputeTargetHttpProxy(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_target_http_proxy",
		Description: "GCP Compute Target Http Proxy",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeTargetHttpProxy,
			Tags:       map[string]string{"service": "compute", "action": "targetHttpProxies.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeTargetHttpProxies,
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
				Name:        "fingerprint",
				Description: "Fingerprint of this resource. A hash of the contents stored in this object, used for optimistic locking.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "http_keep_alive_timeout_sec",
				Description: "Specifies how long to keep a connection open, after completing a response, while there is no matching traffic (in seconds).",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "kind",
				Description: "Type of the resource. Always compute#targetHttpProxy for target HTTPS proxies.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "proxy_bind",
				Description: "This field only applies when the forwarding rule that references this target proxy has a loadBalancingScheme set to INTERNAL_SELF_MANAGED.",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
			},
			{
				Name:        "region",
				Description: "An URL of the region where the regional TargetHttpProxy resides. Not set for global target HTTP proxies.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "url_map",
				Description: "A fully-qualified or valid partial URL to the UrlMap resource that defines the mapping from URL to the BackendService.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "location_type",
				Description: "Location type where the target http proxy resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpComputeTargetHttpProxyLocation, "Type"),
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
				Transform:   transform.From(gcpComputeTargetHttpProxyAka),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpComputeTargetHttpProxyLocation, "Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpComputeTargetHttpProxyLocation, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeTargetHttpProxies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeTargetHttpProxies")

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
	// https://pkg.go.dev/google.golang.org/api/compute/v1#TargetHttpProxiesAggregatedListCall.MaxResults
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

	resp := service.TargetHttpProxies.AggregatedList(project).Filter(filterString).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *compute.TargetHttpProxyAggregatedList) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, item := range page.Items {
			for _, targetHttpProxy := range item.TargetHttpProxies {
				d.StreamListItem(ctx, targetHttpProxy)

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

func getComputeTargetHttpProxy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getComputeTargetHttpProxy")

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

	var targetHttpProxy compute.TargetHttpProxy
	resp := service.TargetHttpProxies.AggregatedList(project).Filter("name=" + name)
	if err := resp.Pages(ctx, func(page *compute.TargetHttpProxyAggregatedList) error {
		for _, item := range page.Items {
			for _, i := range item.TargetHttpProxies {
				targetHttpProxy = *i
			}
		}
		return nil
	},
	); err != nil {
		return nil, err
	}

	if len(targetHttpProxy.Name) < 1 {
		return nil, nil
	}

	return &targetHttpProxy, nil
}

//// TRANSFORM FUNCTIONS

func gcpComputeTargetHttpProxyAka(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*compute.TargetHttpProxy)
	region := getLastPathElement(types.SafeString(data.Region))
	project := strings.Split(data.SelfLink, "/")[6]

	akas := []string{"gcp://compute.googleapis.com/projects/" + project + "/regions/" + region + "/targetHttpProxies/" + data.Name}

	if region == "" {
		akas = []string{"gcp://compute.googleapis.com/projects/" + project + "/global/targetHttpProxies/" + data.Name}
	}

	return akas, nil
}

func gcpComputeTargetHttpProxyLocation(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*compute.TargetHttpProxy)
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
