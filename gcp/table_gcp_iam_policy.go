package gcp

import (
	"context"
	"strings"

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
				Hydrate:     getGcpIamPolicyTurbotData,
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGcpIamPolicyTurbotData,
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
				Transform:   transform.FromConstant(activeProject()),
			},
		},
	}
}

//// FETCH FUNCTIONS

func listGcpIamPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	project := activeProject()
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

//// HYDRATE FUNCTIONS

func getGcpIamPolicyTurbotData(ctx context.Context, d *plugin.QueryData, p *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getGcpIamPolicyTurbotData")
	commonData, err := getCommonColumns(ctx, d, p)
	if err != nil {
		return err, nil
	}

	commonColumnData := commonData.(*gcpCommonColumnData)

	// Get the resource title
	title := strings.ToUpper(commonColumnData.Project) + " IAM Policy"

	// Get data for turbot defined properties
	akas := []string{"gcp://cloudresourcemanager.googleapis.com/projects/" + commonColumnData.Project + "/iamPolicy"}

	// Mapping all turbot defined properties
	turbotData := map[string]interface{}{
		"Akas":  akas,
		"Title": title,
	}

	return turbotData, nil
}
