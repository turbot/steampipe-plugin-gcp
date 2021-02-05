package gcp

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"google.golang.org/api/compute/v1"
)

//// TABLE DEFINITION

func tableGcpComputeURLMap(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_url_map",
		Description: "GCP Compute URL Map",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeURLMap,
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeURLMaps,
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
				Name:        "description",
				Description: "A user-specified, human-readable description of the URL map.",
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
				Name:        "default_service",
				Description: "The full or partial URL of the defaultService resource to which traffic is directed if none of the hostRules match.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "fingerprint",
				Description: "An unique system generated string, to reduce conflicts when multiple users change any property of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "default_url_host_redirect",
				Description: "The host that will be used in the redirect response instead of the one that was supplied in the request. The value must be between 1 and 255 characters.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DefaultUrlRedirect.HostRedirect"),
			},
			{
				Name:        "default_url_https_redirect",
				Description: "Specifies whether the URL scheme in the redirected request is set to https, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("DefaultUrlRedirect.HttpsRedirect"),
			},
			{
				Name:        "default_url_path_redirect",
				Description: "The path that will be used in the redirect response instead of the one that was supplied in the request.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DefaultUrlRedirect.PathRedirect"),
			},
			{
				Name:        "default_url_prefix_redirect",
				Description: "The prefix that replaces the prefixMatch specified in the HttpRouteRuleMatch, retaining the remaining portion of the URL before redirecting the request.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DefaultUrlRedirect.PrefixRedirect"),
			},
			{
				Name:        "default_url_redirect_response_code",
				Description: "Specifies the HTTP Status code to use for this RedirectAction.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DefaultUrlRedirect.RedirectResponseCode"),
			},
			{
				Name:        "default_url_strip_query",
				Description: "Specifies whether any accompanying query portion of the original URL is removed prior to redirecting the request, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("DefaultUrlRedirect.StripQuery"),
			},
			{
				Name:        "region",
				Description: "The URL of the region where the regional backend service resides. This field is not applicable to global backend services.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "The server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "default_route_action",
				Description: "DefaultRouteAction takes effect when none of the hostRules match. The load balancer performs advanced routing actions like URL rewrites, header transformations, etc. prior to forwarding the request to the selected backend. If defaultRouteAction specifies any weightedBackendServices, defaultService must not be set.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "host_rules",
				Description: "The list of HostRules to use against the URL.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "path_matchers",
				Description: "The list of named PathMatchers to use against the URL.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "request_headers_to_add",
				Description: "A list of headers to add to a matching request prior to forwarding the request to the backendService.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("HeaderAction.RequestHeadersToAdd"),
			},
			{
				Name:        "request_headers_to_remove",
				Description: "A list of header names for headers that need to be removed from the request prior to forwarding the request to the backendService.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("HeaderAction.RequestHeadersToRemove"),
			},
			{
				Name:        "response_headers_to_add",
				Description: "A list of headers to add the response prior to sending the response back to the client.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("HeaderAction.ResponseHeadersToAdd"),
			},
			{
				Name:        "response_headers_to_remove",
				Description: "A list of header names for headers that need to be removed from the response prior to sending the response back to the client.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("HeaderAction.ResponseHeadersToRemove"),
			},
			{
				Name:        "tests",
				Description: "The list of expected URL mapping tests. Request to update this UrlMap will succeed only if all of the test cases pass.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "location_type",
				Description: "Location type where the url map resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpComputeURLMapLocation, "Type"),
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
				Transform:   transform.From(gcpComputeURLMapAka),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpComputeURLMapLocation, "Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromConstant(activeProject()),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeURLMaps(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeURLMaps")
	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	project := activeProject()
	resp := service.UrlMaps.AggregatedList(project)
	if err := resp.Pages(ctx, func(page *compute.UrlMapsAggregatedList) error {
		for _, item := range page.Items {
			for _, urlMap := range item.UrlMaps {
				d.StreamListItem(ctx, urlMap)
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getComputeURLMap(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	var urlMap compute.UrlMap
	name := d.KeyColumnQuals["name"].GetStringValue()
	project := activeProject()

	resp := service.UrlMaps.AggregatedList(project).Filter("name=" + name)
	if err := resp.Pages(
		ctx,
		func(page *compute.UrlMapsAggregatedList) error {
			for _, item := range page.Items {
				for _, i := range item.UrlMaps {
					urlMap = *i
				}
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	// If the specified resource is not present, API does not return any not found errors
	if len(urlMap.Name) < 1 {
		return nil, nil
	}

	return &urlMap, nil
}

//// TRANSFORM FUNCTIONS

func gcpComputeURLMapAka(_ context.Context, d *transform.TransformData) (interface{}, error) {
	urlMap := d.HydrateItem.(*compute.UrlMap)
	regionName := getLastPathElement(types.SafeString(urlMap.Region))

	akas := []string{"gcp://compute.googleapis.com/projects/" + activeProject() + "/regions/" + regionName + "/urlMaps/" + urlMap.Name}

	if regionName == "" {
		akas = []string{"gcp://compute.googleapis.com/projects/" + activeProject() + "/global/urlMaps/" + urlMap.Name}
	}

	return akas, nil
}

func gcpComputeURLMapLocation(_ context.Context, d *transform.TransformData) (interface{}, error) {
	urlMap := d.HydrateItem.(*compute.UrlMap)
	param := d.Param.(string)
	regionName := getLastPathElement(types.SafeString(urlMap.Region))

	locationData := map[string]string{
		"Type":     "REGIONAL",
		"Location": regionName,
	}

	if regionName == "" {
		locationData["Type"] = "GLOBAL"
		locationData["Location"] = "global"
	}

	return locationData[param], nil
}