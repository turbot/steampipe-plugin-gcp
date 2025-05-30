package gcp

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"google.golang.org/api/cloudresourcemanager/v1"
)

//// TABLE DEFINITION

func tableGcpIAMPolicy(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_iam_policy",
		Description: "GCP IAM Policy",
		List: &plugin.ListConfig{
			Hydrate: listGcpIamPolicies,
			Tags:    map[string]string{"service": "resourcemanager", "action": "projects.getIamPolicy"},
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
				Hydrate:     getIamPolicyTurbotData,
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIamPolicyTurbotData,
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

func listGcpIamPolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := CloudResourceManagerService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)
	plugin.Logger(ctx).Trace("listGcpIamPolicies", "GCP_PROJECT: ", project)

	rb := &cloudresourcemanager.GetIamPolicyRequest{}
	resp, err := service.Projects.GetIamPolicy(project, rb).Context(ctx).Do()
	// apply rate limiting
	d.WaitForListRateLimit(ctx)

	if err != nil {
		return nil, err
	}
	d.StreamListItem(ctx, resp)

	return nil, nil
}

func getIamPolicyTurbotData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	// Get the resource title
	title := strings.ToUpper(project) + " IAM Policy"

	// Build resource aka
	akas := []string{"gcp://cloudresourcemanager.googleapis.com/projects/" + project + "/iamPolicy"}

	// Mapping all turbot defined properties
	turbotData := map[string]interface{}{
		"Akas":  akas,
		"Title": title,
	}

	return turbotData, nil
}
