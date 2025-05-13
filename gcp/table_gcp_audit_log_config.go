package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"google.golang.org/api/cloudresourcemanager/v1"
)

type auditLogConfigRow struct {
	Service         string   `json:"service"`
	LogType         string   `json:"log_type"`
	ExemptedMembers []string `json:"exempted_members"`
}

//// TABLE DEFINITION

func tableGcpAuditLogConfig(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_audit_log_config",
		Description: "GCP Audit Log Config provides information about audit log configurations for Google Cloud services.",
		List: &plugin.ListConfig{
			Hydrate: listGcpAuditLogConfigs,
			Tags:    map[string]string{"service": "cloudresourcemanager", "action": "getIamPolicy"},
		},
		Columns: []*plugin.Column{
			// Key columns
			{
				Name:        "service",
				Description: "Specifies a service that will be enabled for audit logging.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "log_type",
				Description: "The log type that this config enables. Possible values: ADMIN_READ, DATA_WRITE, DATA_READ.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "exempted_members",
				Description: "Specifies the identities that do not cause logging for this type of permission.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Service"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     auditLogConfigToAkas,
				Transform:   transform.FromValue(),
			},

			// Standard gcp columns
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

//// LIST FUNCTION

func listGcpAuditLogConfigs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := CloudResourceManagerService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details
	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	resp, err := service.Projects.GetIamPolicy(project, &cloudresourcemanager.GetIamPolicyRequest{}).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	for _, auditConfig := range resp.AuditConfigs {
		for _, logConfig := range auditConfig.AuditLogConfigs {
			d.StreamListItem(ctx, auditLogConfigRow{
				Service:         auditConfig.Service,
				LogType:         logConfig.LogType,
				ExemptedMembers: logConfig.ExemptedMembers,
			})
		}
	}

	return nil, nil
}

//// TRANSFORM FUNCTIONS

func auditLogConfigToAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	config := h.Item.(auditLogConfigRow)

	// Get project details
	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	akas := []string{"gcp://cloudresourcemanager.googleapis.com/projects/" + project + "/services/" + config.Service + "/logType/" + config.LogType}
	return akas, nil
}
