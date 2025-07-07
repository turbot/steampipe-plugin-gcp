package gcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/compute/v1"
)

//// TABLE DEFINITION

func tableGcpComputeSecurityPolicy(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_security_policy",
		Description: "GCP Armor Security Policy",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getGcpComputeSecurityPolicy,
			Tags:       map[string]string{"service": "compute", "action": "securityPolicies.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listGcpComputeSecurityPolicies,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "id", Require: plugin.Optional},
				{Name: "type", Require: plugin.Optional},
				{Name: "description", Require: plugin.Optional},
				{Name: "self_link", Require: plugin.Optional},
				{Name: "filter", Require: plugin.Optional, CacheMatch: "exact"},
			},
			Tags: map[string]string{"service": "compute", "action": "securityPolicies.list"},
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the security policy."},
			{Name: "id", Type: proto.ColumnType_INT, Description: "The unique identifier for the resource."},
			{Name: "creation_timestamp", Type: proto.ColumnType_TIMESTAMP, Description: "Creation timestamp in RFC3339 text format."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "An optional description of this resource."},
			{Name: "self_link", Type: proto.ColumnType_STRING, Description: "Server-defined URL for the resource."},
			{Name: "filter", Type: proto.ColumnType_STRING, Transform: transform.FromQual("filter"), Description: "The filter pattern for the search."},
			{Name: "fingerprint", Type: proto.ColumnType_STRING, Description: "Specifies a fingerprint for this resource, which is essentially a hash of the metadata's contents and used for optimistic locking."},
			{Name: "label_fingerprint", Type: proto.ColumnType_STRING, Description: "A fingerprint for the labels being applied to this security policy, which is essentially a hash of the labels set used for optimistic locking."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The type indicates the intended use of the security policy."},
			{Name: "rules", Type: proto.ColumnType_JSON, Description: "A list of rules that belong to this policy."},
			{Name: "adaptive_protection_config", Type: proto.ColumnType_JSON, Description: "Configuration settings for adaptive protection, which automatically detects and mitigates layer 7 DDoS attacks by analyzing traffic patterns and blocking suspicious requests."},
			{Name: "advanced_options_config", Type: proto.ColumnType_JSON, Description: "Configuration for advanced security options, such as JSON parsing and custom header actions applied to requests evaluated by this security policy."},
			{Name: "ddos_Protection_config", Type: proto.ColumnType_JSON, Description: "Configuration that enables or disables DDoS protection for the security policy, specifying the protection level against layer 7 volumetric attacks."},
			{Name: "recaptcha_options_config", Type: proto.ColumnType_JSON, Description: "Configuration settings for Google reCAPTCHA Enterprise integration, including site keys and actions to take on failed reCAPTCHA challenges."},
			{Name: "user_defined_fields", Type: proto.ColumnType_JSON, Description: "Definitions of user-defined fields for CLOUD_ARMOR_NETWORK policies."},
			{Name: "force_send_fields", Type: proto.ColumnType_JSON, Description: "ForceSendFields is a list of field names (e.g. AdaptiveProtectionConfig) to unconditionally include in API requests."},
			{Name: "labels", Type: proto.ColumnType_JSON, Description: "Labels for this resource."},

			// standard steampipe columns
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: ColumnDescriptionTitle},
			{Name: "akas", Type: proto.ColumnType_JSON, Transform: transform.FromP(securityPolicyTurbotData, "Akas"), Description: ColumnDescriptionAkas},
			{Name: "tags", Type: proto.ColumnType_JSON, Transform: transform.FromField("Labels"), Description: ColumnDescriptionTags},

			// standard gcp columns
			{Name: "project", Type: proto.ColumnType_STRING, Transform: transform.FromP(securityPolicyTurbotData, "Project"), Description: ColumnDescriptionProject},
			{Name: "location", Type: proto.ColumnType_STRING, Transform: transform.FromConstant("global"), Description: ColumnDescriptionLocation},
		},
	}
}

//// LIST FUNCTION

func listGcpComputeSecurityPolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details
	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	// Set page size to 500 or query limit if lower
	pageSize := types.Int64(500)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *pageSize {
			pageSize = limit
		}
	}

	filter := ""

	if d.EqualsQualString("filter") != "" {
		filter = d.EqualsQualString("filter")
	} else {
		filter = buildComputeSecurityPolicyFilterParam(d.Quals)
	}

	resp := service.SecurityPolicies.List(project).MaxResults(*pageSize).Filter(filter)
	if err := resp.Pages(ctx, func(page *compute.SecurityPolicyList) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)
		for _, policy := range page.Items {
			d.StreamListItem(ctx, policy)
			//
			if d.RowsRemaining(ctx) == 0 {
				return nil
			}
		}
		return nil
	}); err != nil {
		plugin.Logger(ctx).Error("gcp_compute_security_policy.listGcpComputeSecurityPolicies", "api_error", err)
		return nil, err
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getGcpComputeSecurityPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details
	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)
	name := d.EqualsQualString("name")

	// Return nil, if no input provided
	if name == "" {
		return nil, nil
	}
	policy, err := service.SecurityPolicies.Get(project, name).Do()
	if err != nil {
		return nil, err
	}
	return policy, nil
}

//// UTILITY FUNCTION

func buildComputeSecurityPolicyFilterParam(equalQuals plugin.KeyColumnQualMap) string {
	filter := ""

	filterQuals := []filterQualMap{
		{"type", "type", "string"},
		{"description", "description", "string"},
		{"id", "id", "int"},
		{"self_link", "selfLink", "string"},
	}

	for _, filterQualItem := range filterQuals {
		filterQual := equalQuals[filterQualItem.ColumnName]
		if filterQual == nil {
			continue
		}

		// Check only if filter qual map matches with optional column name
		if filterQual.Name == filterQualItem.ColumnName {
			if filterQual.Quals == nil {
				continue
			}
		}

		for _, qual := range filterQual.Quals {
			if qual.Value != nil {
				value := qual.Value
				switch filterQualItem.Type {
				case "string":
					if filter == "" {
						filter = filterQualItem.PropertyPath + " = \"" + value.GetStringValue() + "\""
					} else {
						filter = filter + " AND " + filterQualItem.PropertyPath + " = \"" + value.GetStringValue() + "\""
					}
				case "int":
					intVal := value.GetInt64Value()
					if filter == "" {
						filter = filterQualItem.PropertyPath + " = " + fmt.Sprintf("%d", intVal)
					} else {
						filter = filter + " AND " + filterQualItem.PropertyPath + " = " + fmt.Sprintf("%d", intVal)
					}
				}
			}
		}
	}
	return filter
}

////  TRANSFORM FUNCTION

func securityPolicyTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*compute.SecurityPolicy)
	project := ""
	if data.SelfLink != "" {
		parts := strings.Split(data.SelfLink, "/")
		if len(parts) > 6 {
			project = parts[6]
		}
	}
	turbotData := map[string]interface{}{
		"Project": project,
		"Akas":    []string{"gcp://compute.googleapis.com/projects/" + project + "/global/securityPolicies/" + data.Name},
	}
	return turbotData[d.Param.(string)], nil
}
