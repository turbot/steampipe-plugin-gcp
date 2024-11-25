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

func tableGcpOrganizationProject(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_organization_project",
		Description: "GCP Organization Project",
		List: &plugin.ListConfig{
			Hydrate: listGCPOrganizationProjects,
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
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(projectSelfLink),
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
			{
				Name:        "access_approval_settings",
				Description: "The access approval settings associated with this project.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getProjectAccessApprovalSettings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "ancestors",
				Description: "The ancestors of the project in the resource hierarchy, from bottom to top.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getProjectAncestors,
				Transform:   transform.FromValue(),
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

func listGCPOrganizationProjects(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := CloudResourceManagerService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_organization_project.listGCPOrganizationProjects", "service_err", err)
		return nil, err
	}

	// Max limit is not documented
	pageSize := types.Int64(500)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *pageSize {
			pageSize = limit
		}
	}

	// List projects
	resp := service.Projects.List().PageSize(*pageSize)
	if err := resp.Pages(ctx, func(page *cloudresourcemanager.ListProjectsResponse) error {
		for _, project := range page.Projects {
			d.StreamListItem(ctx, project)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				page.NextPageToken = ""
				return nil
			}
		}
		return nil
	}); err != nil {
		plugin.Logger(ctx).Error("gcp_organization_project.listGCPOrganizationProjects", "api_err", err)
		return nil, err
	}

	return nil, nil
}
