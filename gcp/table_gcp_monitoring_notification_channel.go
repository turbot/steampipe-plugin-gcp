package gcp

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

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
		},
		List: &plugin.ListConfig{
			Hydrate: listGcpMonitoringNotificationChannels,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The full REST resource name for this channel.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(notificationChannelNameToTurbotData, "Name"),
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
				Transform:   transform.FromConstant(projectName),
			},
		},
	}
}

//// FETCH FUNCTIONS

func listGcpMonitoringNotificationChannels(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	service, err := monitoring.NewService(ctx)
	if err != nil {
		return nil, err
	}

	project := projectName
	resp := service.Projects.NotificationChannels.List("projects/" + project)
	if err := resp.Pages(
		ctx,
		func(page *monitoring.ListNotificationChannelsResponse) error {
			for _, notificationChannel := range page.NotificationChannels {
				d.StreamListItem(ctx, notificationChannel)
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
	logger := plugin.Logger(ctx)
	logger.Trace("getGcpMonitoringNotificationChannel")

	project := projectName
	name := d.KeyColumnQuals["name"].GetStringValue()

	service, err := monitoring.NewService(ctx)
	if err != nil {
		return nil, err
	}

	op, err := service.Projects.NotificationChannels.Get("projects/" + project + "/notificationChannels/" + name).Do()
	if err != nil {
		logger.Debug("getGcpMonitoringNotificationChannel__", "ERROR", err)
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
		"Name":  splittedTitle[len(splittedTitle)-1],
		"Title": title,
		"Akas":  []string{"gcp://monitoring.googleapis.com/" + notificationChannel.Name},
	}
	return turbotData[param], nil
}
