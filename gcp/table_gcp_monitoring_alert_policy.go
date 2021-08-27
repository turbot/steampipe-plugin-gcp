package gcp

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	"google.golang.org/api/monitoring/v3"
)

func tableGcpMonitoringAlert(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_monitoring_alert_policy",
		Description: "GCP Monitoring Alert Policy",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getMonitoringAlertPolicy,
		},
		List: &plugin.ListConfig{
			Hydrate:           listMonitoringAlertPolicies,
			ShouldIgnoreError: isIgnorableError([]string{"403"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "display_name",
				Description: "A short name or phrase used to identify the policy in dashboards, notifications and incidents.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "The resource name for this policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(lastPathElement),
			},
			{
				Name:        "enabled",
				Description: "Indicates whether the policy is enabled, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "combiner",
				Description: "How to combine the results of multiple conditions to determine if an incident should be opened.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_record",
				Description: "A read-only record of the creation of the alerting policy.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "mutation_record",
				Description: "A read-only record of the most recent change to the alerting policy.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "conditions",
				Description: "A list of conditions for the policy.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "documentation",
				Description: "Documentation that is included with notifications and incidents related to this policy.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "notification_channels",
				Description: "Identifies the notification channels to which notifications should be sent when incidents are opened or closed or when new violations occur on an already opened incident.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "user_labels",
				Description: "User-supplied key/value data to be used for organizing and identifying the AlertPolicy objects.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "validity",
				Description: "Read-only description of how the alert policy is invalid.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("UserLabels"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(monitoringAlertPolicyTurbotData, "Akas"),
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
				Transform:   transform.FromP(monitoringAlertPolicyTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listMonitoringAlertPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
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

	resp := service.Projects.AlertPolicies.List("projects/" + project)
	if err := resp.Pages(ctx, func(page *monitoring.ListAlertPoliciesResponse) error {
		for _, alertPolicy := range page.AlertPolicies {
			d.StreamListItem(ctx, alertPolicy)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getMonitoringAlertPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getMonitoringAlertPolicy")

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

	req, err := service.Projects.AlertPolicies.Get("projects/" + project + "/alertPolicies/" + name).Do()
	if err != nil {
		return nil, err
	}

	// If the name has been passed as empty string, API does not return any error
	if len(req.Name) < 1 {
		return nil, nil
	}

	return req, nil
}

//// TRANSFORM FUNCTIONS

func monitoringAlertPolicyTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*monitoring.AlertPolicy)
	param := d.Param.(string)

	splittedData := strings.Split(data.Name, "/")

	turbotData := map[string]interface{}{
		"Project": splittedData[1],
		"Akas":    []string{"gcp://monitoring.googleapis.com/" + data.Name},
	}

	return turbotData[param], nil
}
