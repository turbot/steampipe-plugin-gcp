package gcp

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"google.golang.org/api/compute/v1"
)

//// TABLE DEFINITION

func tableGcpComputeTargetSslProxy(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_target_ssl_proxy",
		Description: "GCP Compute Target SSL Proxy",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeTargetSslProxy,
		},
		List: &plugin.ListConfig{
			Hydrate:           listComputeTargetSslProxies,
			ShouldIgnoreError: isIgnorableError([]string{"403"}),
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
				Name:        "kind",
				Description: "The type of the resource. Always compute#targetSslProxy for target SSL proxies.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A user-specified, human-readable description of the target ssl proxy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_timestamp",
				Description: "The creation timestamp of the resource.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "proxy_header",
				Description: "Specifies the type of proxy header to append before sending data to the backend, either NONE or PROXY_V1. The default is NONE.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "The server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service",
				Description: "Specifies the url of the backend service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ssl_policy",
				Description: "The URL of the SslPolicy resource that will be associated with the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ssl_certificates",
				Description: "A list of urls to SslCertificate resources that are used to authenticate connections to Backends.",
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
				Transform:   transform.FromP(computeTargetSslProxyTurbotData, "Akas"),
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
				Transform:   transform.FromP(computeTargetSslProxyTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeTargetSslProxies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeTargetSslProxies")

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details
	getProjectCached := plugin.HydrateFunc(getProject).WithCache()
	projectId, err := getProjectCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	resp := service.TargetSslProxies.List(project)
	if err := resp.Pages(ctx, func(page *compute.TargetSslProxyList) error {
		for _, targetSslProxy := range page.Items {
			d.StreamListItem(ctx, targetSslProxy)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getComputeTargetSslProxy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getComputeTargetSslProxy")

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details
	getProjectCached := plugin.HydrateFunc(getProject).WithCache()
	projectId, err := getProjectCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)
	name := d.KeyColumnQuals["name"].GetStringValue()

	// Error: json: invalid use of ,string struct tag, trying to unmarshal "projects/<project_name>/global/targetSslProxies/" into uint64
	if len(name) < 1 {
		return nil, nil
	}

	resp, err := service.TargetSslProxies.Get(project, name).Do()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

//// TRANSFORM FUNCTIONS

func computeTargetSslProxyTurbotData(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*compute.TargetSslProxy)
	param := d.Param.(string)

	project := strings.Split(data.SelfLink, "/")[6]

	turbotData := map[string]interface{}{
		"Project": project,
		"Akas":    []string{"gcp://compute.googleapis.com/projects/" + project + "/global/targetSslProxies/" + data.Name},
	}

	return turbotData[param], nil
}
