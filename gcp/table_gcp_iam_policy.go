package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"google.golang.org/api/cloudresourcemanager/v1"
)

//// TABLE DEFINITION

func tableGcpIAMPolicy(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_iam_policy",
		Description: "GCP IAM Policy",
		List: &plugin.ListConfig{
			Hydrate: listGcpIamPolicies,
		},
		Columns: []*plugin.Column{
			{
				Name:        "version",
				Description: "Version specifies the format of the policy. Valid values are `0`, `1`, and `3`.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "etag",
				Description: "Etag is used for optimistic concurrency control as a way to help prevent simultaneous updates of a policy from overwriting each other.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "bindings",
				Description: "A list of `members` to a `role`. Optionally, may specify a `condition` that determines how and when the `bindings` are applied. Each of the `bindings` must contain at least one member.",
				Type:        proto.ColumnType_JSON,
			},

			// standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(iamPolicyTurbotData, "Title"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(iamPolicyTurbotData, "Akas"),
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
				Transform:   transform.FromConstant(projectName),
			},
		},
	}
}

//// FETCH FUNCTIONS

func listGcpIamPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	project := projectName
	plugin.Logger(ctx).Trace("listGcpIamPolicies", "GCP_PROJECT: ", project)

	service, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		return nil, err
	}

	rb := &cloudresourcemanager.GetIamPolicyRequest{}
	resp, err := service.Projects.GetIamPolicy(project, rb).Context(ctx).Do()
	if err != nil {
		return nil, err
	}
	d.StreamListItem(ctx, resp)

	return nil, nil
}

//// TRANSFORM FUNCTION

func iamPolicyTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	project := projectName
	param := types.SafeString(d.Param)

	// Get the resource title
	title := strings.ToUpper(project) + " IAM Policy"

	// Get data for turbot defined properties
	akas := []string{"gcp://cloudresourcemanager.googleapis.com/projects/" + project + "/iamPolicy"}

	// Mapping all turbot defined properties
	turbotData := map[string]interface{}{
		"Akas":  akas,
		"Title": title,
	}

	return turbotData[param], nil
}
