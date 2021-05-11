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

func tableGcpOrganizationPolicy(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_organization_policy",
		Description: "GCP Organization Policy",
		List: &plugin.ListConfig{
			Hydrate: listGcpOrganizationPolicies,
		},
		Columns: []*plugin.Column{
			{
				Name:        "constraint",
				Description: "The name of the Constraint the Policy is configuring, for example, constraints/serviceuser.services.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "etag",
				Description: "An opaque tag indicating the current version of the Policy, used for concurrency control.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "updateTime",
				Description: "The time stamp the Policy was previously updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "version",
				Description: "Version of the Policy. Default version is 0.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "listPolicy",
				Description: "List of values either allowed or disallowed.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "booleanPolicy",
				Description: "For boolean Constraints, whether to enforce the Constraint or not.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "restoreDefault",
				Description: "Restores the default behavior of the constraint; independent of Constraint type.",
				Type:        proto.ColumnType_JSON,
			},

			// standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Hydrate:     listGcpOrganizationPolicies,
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getOrganizationPolicyTurbotData,
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

func listGcpOrganizationPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := CloudResourceManagerService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}
	project := projectData.Project
	plugin.Logger(ctx).Trace("listGcpOrganizationPolicies", "GCP_PROJECT: ", project)

	rb := &cloudresourcemanager.ListOrgPoliciesRequest{}
	resp, err := service.Projects.ListOrgPolicies(project, rb).Context(ctx).Do()
	if err != nil {
		return nil, err
	}
	d.StreamListItem(ctx, resp)

	return nil, nil
}

func getOrganizationPolicyTurbotData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get project details
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}
	project := projectData.Project

	// Get the resource title
	title := strings.ToUpper(project) + " Org Policy"

	// Build resource aka
	akas := []string{"gcp://cloudresourcemanager.googleapis.com/projects/" + project + "/OrgPolicy"}

	// Mapping all turbot defined properties
	turbotData := map[string]interface{}{
		"Akas":  akas,
		"Title": title,
	}

	return turbotData, nil
}
