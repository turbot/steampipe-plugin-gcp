package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"google.golang.org/api/cloudresourcemanager/v1"
)

//// TABLE DEFINITION

func tableGcpAuditPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_audit_policy",
		Description: "GCP Audit Policy",
		List: &plugin.ListConfig{
			Hydrate: listGcpAuditPolicies,
			Tags:    map[string]string{"service": "resourcemanager", "action": "projects.getIamPolicy"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "service",
				Description: "Specifies a service that will be enabled for audit logging",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "audit_log_configs",
				Description: "The configuration for logging of each type of permission",
				Type:        proto.ColumnType_JSON,
			},

			// standard steampipe columns
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     serviceNameToAkas,
				Transform:   transform.FromValue(),
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
				Hydrate:     getProject,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listGcpAuditPolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
	// apply rate limiting
	d.WaitForListRateLimit(ctx)
	if err != nil {
		return nil, err
	}

	for _, auditConfig := range resp.AuditConfigs {
		d.StreamListItem(ctx, auditConfig)
	}

	return nil, nil
}

func serviceNameToAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	auditConfig := h.Item.(*cloudresourcemanager.AuditConfig)

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	akas := []string{"gcp://cloudresourcemanager.googleapis.com/projects/" + project + "/services/" + auditConfig.Service}
	return akas, nil
}
