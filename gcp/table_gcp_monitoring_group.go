package gcp

import (
	"context"
	"os"
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
			KeyColumns:        plugin.SingleColumn("name"),
			Hydrate:           getMonitoringGroup,
			ShouldIgnoreError: isNotFoundError([]string{"400"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listMonitoringGroup,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of this group",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(groupInfoToTurbotData, "Name"),
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
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(groupInfoToTurbotData, "Title"),
			},
			{
				Name:        "akas",
				Description: "Array of globally unique identifier strings (also known as) for the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(groupInfoToTurbotData, "Akas"),
			},
			{
				Name:        "project",
				Description: "The Google Project in which the resource is located",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(groupInfoToTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listMonitoringGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	project := os.Getenv("GCP_PROJECT")

	service, err := monitoring.NewService(ctx)
	if err != nil {
		return nil, err
	}

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
	project := os.Getenv("GCP_PROJECT")
	name := d.KeyColumnQuals["name"].GetStringValue()

	service, err := monitoring.NewService(ctx)
	if err != nil {
		return nil, err
	}

	req, err := service.Projects.Groups.Get("projects/" + project + "/groups/" + name).Do()
	if err != nil {
		return nil, err
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
		"Name":    splittedTitle[len(splittedTitle)-1],
		"Title":   title,
		"Akas":    []string{"gcp://monitoring.googleapis.com/" + group.Name},
	}

	return turbotData[param], nil
}
