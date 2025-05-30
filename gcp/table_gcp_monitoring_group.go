package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/monitoring/v3"
)

func tableGcpMonitoringGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_monitoring_group",
		Description: "GCP Monitoring Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getMonitoringGroup,
			Tags:       map[string]string{"service": "monitoring", "action": "groups.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listMonitoringGroup,
			Tags:    map[string]string{"service": "monitoring", "action": "groups.list"},
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

func listMonitoringGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := MonitoringService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Max limit isn't mentioned in the documentation
	// Default limit is set as 1000
	pageSize := types.Int64(1000)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *pageSize {
			pageSize = limit
		}
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	resp := service.Projects.Groups.List("projects/" + project).PageSize(*pageSize)

	if err := resp.Pages(ctx, func(page *monitoring.ListGroupsResponse) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, group := range page.Group {
			d.StreamListItem(ctx, group)

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

func getMonitoringGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getMonitoringGroup")

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)
	name := d.EqualsQuals["name"].GetStringValue()

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
