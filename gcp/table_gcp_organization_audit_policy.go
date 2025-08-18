package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"google.golang.org/api/cloudresourcemanager/v1"
)

//// TABLE DEFINITION

func tableGcpOrganizationAuditPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_organization_audit_policy",
		Description: "GCP Organization Audit Policy",
		List: &plugin.ListConfig{
			Hydrate:       listGcpOrganizationAuditPolicies,
			ParentHydrate: listGCPOrganizations,
			Tags:          map[string]string{"service": "resourcemanager", "action": "organizations.getIamPolicy"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "organization_id",
				Description: "The unique identifier for the organization.",
				Type:        proto.ColumnType_STRING,
			},
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
				Hydrate:     organizationServiceNameToAkas,
				Transform:   transform.FromValue(),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromConstant("global"),
			},
		},
	}
}

//// LIST FUNCTION

func listGcpOrganizationAuditPolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := CloudResourceManagerService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get the organization from the parent hydrate
	organization := h.Item.(*cloudresourcemanager.Organization)
	organizationId := getLastPathElement(organization.Name)

	resp, err := service.Organizations.GetIamPolicy("organizations/"+organizationId, &cloudresourcemanager.GetIamPolicyRequest{}).Context(ctx).Do()
	
	// apply rate limiting
	d.WaitForListRateLimit(ctx)
	
	if err != nil {
		plugin.Logger(ctx).Error("listGcpOrganizationAuditPolicies", "organization_id", organizationId, "error", err)
		return nil, err
	}

	for _, auditConfig := range resp.AuditConfigs {

		// Add organization_id to the audit config for reference
		auditConfigWithOrg := &OrganizationAuditConfig{
			OrganizationId:  organizationId,
			Service:         auditConfig.Service,
			AuditLogConfigs: auditConfig.AuditLogConfigs,
		}
		d.StreamListItem(ctx, auditConfigWithOrg)

		// Check if context has been cancelled or if the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			break
		}
	}

	return nil, nil
}

func organizationServiceNameToAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	auditConfig := h.Item.(*OrganizationAuditConfig)

	akas := []string{"gcp://cloudresourcemanager.googleapis.com/organizations/" + auditConfig.OrganizationId + "/services/" + auditConfig.Service}
	return akas, nil
}

// OrganizationAuditConfig is a custom struct to include organization_id with audit config
type OrganizationAuditConfig struct {
	OrganizationId  string                                 `json:"organization_id"`
	Service         string                                 `json:"service"`
	AuditLogConfigs []*cloudresourcemanager.AuditLogConfig `json:"audit_log_configs"`
}
