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

func tableGcpComputeSslPolicy(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_ssl_policy",
		Description: "GCP Compute SSL Policy",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeSslPolicy,
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeSslPolicies,
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
				Description: "The type of the resource. Always compute#sslPolicyfor SSL policies.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A user-specified, human-readable description of the SSL policy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_timestamp",
				Description: "The creation timestamp of the resource.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "fingerprint",
				Description: "A hash of the contents stored in this object. An up-to-date fingerprint must be provided in order to update the SslPolicy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "min_tls_version",
				Description: "The minimum version of SSL protocol that can be used by the clients to establish a connection with the load balancer. Valid values are TLS_1_0, TLS_1_1 and TLS_1_2.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "profile",
				Description: "Profile specifies the set of SSL features that can be used by the load balancer when negotiating SSL with clients.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "The server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "custom_features",
				Description: "A list of features enabled when the selected profile is CUSTOM.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getComputeSslPolicy,
			},
			{
				Name:        "enabled_features",
				Description: "A list of features enabled in the SSL policy.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getComputeSslPolicy,
			},
			{
				Name:        "warnings",
				Description: "A list of warning messages, if any potential misconfigurations are detected for this SSL policy.",
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
				Transform:   transform.FromP(computeSslPolicyTurbotData, "Akas"),
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
				Transform:   transform.FromP(computeSslPolicyTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeSslPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeSslPolicies")

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}
	project := projectData.Project

	resp := service.SslPolicies.List(project)
	if err := resp.Pages(ctx, func(page *compute.SslPoliciesList) error {
		for _, sslPolicy := range page.Items {
			d.StreamListItem(ctx, sslPolicy)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getComputeSslPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getComputeSslPolicy")

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}
	project := projectData.Project

	var name string
	if h.Item != nil {
		name = h.Item.(*compute.SslPolicy).Name
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
	}

	resp, err := service.SslPolicies.Get(project, name).Do()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

//// TRANSFORM FUNCTIONS

func computeSslPolicyTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*compute.SslPolicy)
	param := d.Param.(string)

	project := strings.Split(data.SelfLink, "/")[6]

	turbotData := map[string]interface{}{
		"Project": project,
		"Akas":    []string{"gcp://compute.googleapis.com/projects/" + project + "/global/sslPolicies/" + data.Name},
	}

	return turbotData[param], nil
}
