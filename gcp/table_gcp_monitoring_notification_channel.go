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

//// TABLE DEFINITION

func tableGcpMonitoringNotificationChannel(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_monitoring_notification_channel",
		Description: "GCP Monitoring Notification Channel",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getGcpMonitoringNotificationChannel,
			Tags:       map[string]string{"service": "monitoring", "action": "notificationChannels.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listGcpMonitoringNotificationChannels,
			KeyColumns: plugin.KeyColumnSlice{
				// String columns
				{Name: "type", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "display_name", Require: plugin.Optional, Operators: []string{"<>", "="}},
			},
			Tags: map[string]string{"service": "monitoring", "action": "notificationChannels.list"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The full REST resource name for this channel.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(lastPathElement),
			},
			{
				Name:        "display_name",
				Description: "A human-readable name for the notification channel.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "enabled",
				Description: "Specifies whether the notifications are forwarded to the described channel, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(notificationChannelSelfLink),
			},
			{
				Name:        "type",
				Description: "Specifies the type of the notification channel.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A user-specified, human-readable description of the notification channel.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "verification_status",
				Description: "Specifies whether this channel has been verified, or not.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "user_labels",
				Description: "A list of user-supplied key/value data that does not need to conform to the corresponding NotificationChannelDescriptor's schema unlike the labels field.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "labels",
				Description: "A set of labels attached with the notification channel.",
				Type:        proto.ColumnType_JSON,
			},

			// standard steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Labels"),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(notificationChannelNameToTurbotData, "Title"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(notificationChannelNameToTurbotData, "Akas"),
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
				Transform:   transform.FromP(notificationChannelNameToTurbotData, "Project"),
			},
		},
	}
}

//// FETCH FUNCTIONS

func listGcpMonitoringNotificationChannels(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := MonitoringService(ctx, d)
	if err != nil {
		return nil, err
	}

	filterQuals := []filterQualMap{
		{"type", "type", "string"},
		{"display_name", "displayName", "string"},
	}

	filters := buildQueryFilterFromQuals(filterQuals, d.Quals)
	filterString := ""
	if len(filters) > 0 {
		filterString = strings.Join(filters, " ")
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

	resp := service.Projects.NotificationChannels.List("projects/" + project).Filter(filterString).PageSize(*pageSize)
	if err := resp.Pages(
		ctx,
		func(page *monitoring.ListNotificationChannelsResponse) error {
			// apply rate limiting
			d.WaitForListRateLimit(ctx)

			for _, notificationChannel := range page.NotificationChannels {
				d.StreamListItem(ctx, notificationChannel)

				// Check if context has been cancelled or if the limit has been hit (if specified)
				// if there is a limit, it will return the number of rows required to reach this limit
				if d.RowsRemaining(ctx) == 0 {
					page.NextPageToken = ""
					return nil
				}
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getGcpMonitoringNotificationChannel(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getGcpMonitoringNotificationChannel")

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

	op, err := service.Projects.NotificationChannels.Get("projects/" + project + "/notificationChannels/" + name).Do()
	if err != nil {
		plugin.Logger(ctx).Debug("getGcpMonitoringNotificationChannel__", "ERROR", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func notificationChannelNameToTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	notificationChannel := d.HydrateItem.(*monitoring.NotificationChannel)
	param := d.Param.(string)

	// get the resource title
	splittedTitle := strings.Split(notificationChannel.Name, "/")

	var title string
	if notificationChannel.DisplayName != "" {
		title = notificationChannel.DisplayName
	} else {
		title = splittedTitle[len(splittedTitle)-1]
	}

	turbotData := map[string]interface{}{
		"Project": splittedTitle[1],
		"Title":   title,
		"Akas":    []string{"gcp://monitoring.googleapis.com/" + notificationChannel.Name},
	}
	return turbotData[param], nil
}

func notificationChannelSelfLink(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*monitoring.NotificationChannel)
	selfLink := "https://monitoring.googleapis.com/v3/" + data.Name

	return selfLink, nil
}
