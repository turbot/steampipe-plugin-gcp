package gcp

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"google.golang.org/api/cloudresourcemanager/v1"
)

//// TABLE DEFINITION

func tableGcpProjectOrganizationPolicy(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_project_organization_policy",
		Description: "GCP Project Organization Policy",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getProjectOrganizationPolicy,
			Tags:       map[string]string{"service": "resourcemanager", "action": "projects.getIamPolicy"},
		},
		List: &plugin.ListConfig{
			Hydrate: listProjectOrganizationPolicies,
			Tags:    map[string]string{"service": "resourcemanager", "action": "projects.list"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The name of the Constraint the Policy is configuring.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Constraint").Transform(lastPathElement),
			},
			{
				Name:        "update_time",
				Description: "The time stamp the Policy was previously updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "version",
				Description: "Version of the Policy. Default version is 0.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "etag",
				Description: "An opaque tag indicating the current version of the Policy, used for concurrency control.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "list_policy",
				Description: "List of values either allowed or disallowed.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "boolean_policy",
				Description: "For boolean Constraints, whether to enforce the Constraint or not.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "restore_default",
				Description: "Restores the default behavior of the constraint; independent of Constraint type.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Constraint").Transform(lastPathElement),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getOrganizationPolicyAkas,
				Transform:   transform.FromValue(),
			},

			// GCP standard columns
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

//// LIST FUNCTION

func listProjectOrganizationPolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

	// Max limit isn't mentioned in the documentation
	// Default limit is set as 1000
	rb := &cloudresourcemanager.ListOrgPoliciesRequest{
		PageSize: *types.Int64(1000),
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < rb.PageSize {
			rb.PageSize = *limit
		}
	}

	resp := service.Projects.ListOrgPolicies("projects/"+project, rb)
	if err := resp.Pages(ctx, func(page *cloudresourcemanager.ListOrgPoliciesResponse) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, orgPolicy := range page.Policies {
			d.StreamListItem(ctx, orgPolicy)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				page.NextPageToken = ""
				return nil
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getProjectOrganizationPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getProjectOrganizationPolicy")

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

	id := d.EqualsQuals["id"].GetStringValue()
	rb := &cloudresourcemanager.GetOrgPolicyRequest{
		Constraint: "constraints/" + id,
	}

	req, err := service.Projects.GetOrgPolicy("projects/"+project, rb).Do()
	if err != nil {
		plugin.Logger(ctx).Debug("getProjectOrganizationPolicy", "ERROR", err)
		return nil, err
	}
	return req, nil
}

func getOrganizationPolicyAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	// Build resource aka
	akas := []string{"gcp://cloudresourcemanager.googleapis.com/projects/" + project}

	return akas, nil
}
