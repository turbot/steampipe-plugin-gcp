package gcp

import (
	"context"
	"os"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"google.golang.org/api/cloudresourcemanager/v1"
)

//// TABLE DEFINITION

func tableGcpAuditPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_audit_policy",
		Description: "GCP Audit Policy",
		List: &plugin.ListConfig{
			Hydrate: listGcpAuditPolicies,
		},
		Columns: gcpColumns([]*plugin.Column{
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
			{
				Name:        "akas",
				Description: "Array of globally unique identifier strings (also known as) for the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(serviceNameToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listGcpAuditPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	service, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		return nil, err
	}

	project := os.Getenv("GCP_PROJECT")
	resp, err := service.Projects.GetIamPolicy(project, &cloudresourcemanager.GetIamPolicyRequest{}).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	for _, auditConfig := range resp.AuditConfigs {
		d.StreamListItem(ctx, auditConfig)
	}

	return nil, nil
}

//// TRANSFORM FUNCTIONS

func serviceNameToAkas(_ context.Context, d *transform.TransformData) (interface{}, error) {
	auditConfig := d.HydrateItem.(*cloudresourcemanager.AuditConfig)
	project := os.Getenv("GCP_PROJECT")

	akas := []string{"gcp://cloudresourcemanager.googleapis.com/projects/" + project + "/services/" + auditConfig.Service}

	return akas, nil
}
