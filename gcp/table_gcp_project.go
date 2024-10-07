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

func tableGcpProject(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_project",
		Description: "GCP Project",
		List: &plugin.ListConfig{
			Hydrate: listGCPProjects,
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
			{
				Name:        "billing_information",
				Description: "The billing information of the project.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getProjectBillingInfo,
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

func listGCPProjects(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := CloudResourceManagerService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_project.listGCPProjects", "service_err", err)
		return nil, err
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)
	plugin.Logger(ctx).Debug("gcp_project.listGCPProjects", "project_id", project)

	resp, err := service.Projects.List().Filter("id=" + project).Do()
	if err != nil {
		plugin.Logger(ctx).Error("gcp_project.listGCPProjects", "api_err", err)
		return nil, err
	}

	for _, project := range resp.Projects {
		d.StreamListItem(ctx, project)
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getProjectAka(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get project details
	project := h.Item.(*cloudresourcemanager.Project)

	// Build resource aka
	akas := []string{"gcp://cloudresourcemanager.googleapis.com/projects/" + project.ProjectId}

	return akas, nil
}

func getProjectAccessApprovalSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := AccessApprovalService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_project.getProjectAccessApprovalSettings", "connection_error", err)
		return nil, err
	}

	// Get project details
	projectId := h.Item.(*cloudresourcemanager.Project).ProjectId

	resp, err := service.Projects.GetAccessApprovalSettings("projects/" + projectId + "/accessApprovalSettings").Do()
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("gcp_project.getProjectAccessApprovalSettings", "api_err", err)
		return nil, err
	}
	return resp, nil
}

func getProjectAncestors(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := CloudResourceManagerService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_project.getProjectAncestors", "connection_error", err)
		return nil, err
	}

	// Get project details
	projectId := h.Item.(*cloudresourcemanager.Project).ProjectId

	resp, err := service.Projects.GetAncestry(projectId, &cloudresourcemanager.GetAncestryRequest{}).Do()
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("gcp_project.getProjectAncestors", "api_err", err)
		return nil, err
	}
	return resp.Ancestor, nil
}

func getProjectBillingInfo(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := BillingService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_project.getProjectBillingInformation", "connection_error", err)
		return nil, err
	}

	// Get project details
	projectId := h.Item.(*cloudresourcemanager.Project).ProjectId

	resp, err := service.Projects.GetBillingInfo("projects/" + projectId).Do()
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, nil
		}
		plugin.Logger(ctx).Error("gcp_project.getProjectBillingInformation", "api_err", err)
		return nil, err
	}
	return resp, nil
}

func projectSelfLink(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*cloudresourcemanager.Project)
	selfLink := "https://cloudresourcemanager.googleapis.com/v1/projects/" + data.Name

	return selfLink, nil
}
