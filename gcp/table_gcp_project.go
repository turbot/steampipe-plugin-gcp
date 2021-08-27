package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	"google.golang.org/api/cloudresourcemanager/v1"
)

//// TABLE DEFINITION

func tableGcpProject(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_project",
		Description: "GCP Project",
		List: &plugin.ListConfig{
			Hydrate:           listGCPProjects,
			ShouldIgnoreError: isIgnorableError([]string{"403"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the project.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "project_id",
				Description: "An unique, user-assigned ID of the Project.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "project_number",
				Description: "The number uniquely identifying the project.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "lifecycle_state",
				Description: "Specifies the project lifecycle state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "Creation time of the project.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "parent",
				Description: "An optional reference to a parent Resource.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "labels",
				Description: "A list of labels attached to this project.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Labels"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getProjectAka,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listGCPProjects(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listGCPProjects")

	// Create Service Connection
	service, err := CloudResourceManagerService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get curent project
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}
	project := projectData.Project

	resp, err := service.Projects.List().Filter("name=" + project).Do()
	for _, project := range resp.Projects {
		d.StreamListItem(ctx, project)
	}
	if err != nil {
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getProjectAka(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getProjectAka")

	// Get project details
	project := h.Item.(*cloudresourcemanager.Project)

	// Build resource aka
	akas := []string{"gcp://cloudresourcemanager.googleapis.com/projects/" + project.ProjectId}

	return akas, nil
}
