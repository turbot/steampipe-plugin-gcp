package gcp

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	"google.golang.org/api/monitoring/v3"
)

func tableGcpMonitoringGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_monitoring_group",
		Description: "GCP Monitoring Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getMonitoringGroup,
		},
		List: &plugin.ListConfig{
			Hydrate:           listMonitoringGroup,
			ShouldIgnoreError: isIgnorableError([]string{"403"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of this group",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(lastPathElement),
			},
			{
				Name:        "display_name",
				Description: "A user-assigned name for this group, used only for display purposes.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "filter",
				Description: "The filter used to determine which monitored resources belong to this group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_cluster",
				Description: "If true, the members of this group are considered to be a cluster. The system can perform additional analysis on groups that are clusters.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "parent_name",
				Description: "The name of the group's parent, if it has one.",
				Type:        proto.ColumnType_STRING,
			},

			// standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(groupInfoToTurbotData, "Title"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(groupInfoToTurbotData, "Akas"),
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
				Transform:   transform.FromP(groupInfoToTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listMonitoringGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := MonitoringService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}
	project := projectData.Project

	resp := service.Projects.Groups.List("projects/" + project)

	if err := resp.Pages(ctx, func(page *monitoring.ListGroupsResponse) error {
		for _, group := range page.Group {
			d.StreamListItem(ctx, group)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getMonitoringGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getMonitoringGroup")

	// Get project details
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}
	project := projectData.Project
	name := d.KeyColumnQuals["name"].GetStringValue()

	// Create Service Connection
	service, err := MonitoringService(ctx, d)
	if err != nil {
		return nil, err
	}

	req, err := service.Projects.Groups.Get("projects/" + project + "/groups/" + name).Do()
	if err != nil {
		return nil, err
	}

	// If the name has been passed as empty string, API does not returns any error
	if len(req.Name) < 1 {
		return nil, nil
	}

	return req, nil
}

//// TRANSFORM FUNCTIONS

func groupInfoToTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	group := d.HydrateItem.(*monitoring.Group)
	param := d.Param.(string)

	// get the resource title
	splittedTitle := strings.Split(group.Name, "/")

	var title string
	if group.DisplayName != "" {
		title = group.DisplayName
	} else {
		title = splittedTitle[len(splittedTitle)-1]
	}

	turbotData := map[string]interface{}{
		"Project": splittedTitle[1],
		"Title":   title,
		"Akas":    []string{"gcp://monitoring.googleapis.com/" + group.Name},
	}

	return turbotData[param], nil
}
