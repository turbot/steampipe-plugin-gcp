package gcp

import (
	"context"
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
		Description: "Google Cloud Armor Security Policy",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getGcpComputeSecurityPolicy,
		},
		List: &plugin.ListConfig{
			Hydrate: listGcpComputeSecurityPolicies,
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the security policy."},
			{Name: "id", Type: proto.ColumnType_INT, Description: "The unique identifier for the resource."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "An optional description of this resource."},
			{Name: "self_link", Type: proto.ColumnType_STRING, Description: "The fully-qualified URL linking back to the resource."},
			{Name: "fingerprint", Type: proto.ColumnType_STRING, Description: "Fingerprint of this resource. A hash of the contents stored in this object."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The type of the security policy."},
			{Name: "rules", Type: proto.ColumnType_JSON, Description: "The set of rules that belong to this policy."},
			{Name: "adaptive_protection_config", Type: proto.ColumnType_JSON, Description: "Configuration for adaptive protection."},
			{Name: "advanced_options_config", Type: proto.ColumnType_JSON, Description: "Advanced options configuration."},
			{Name: "recaptcha_options_config", Type: proto.ColumnType_JSON, Description: "reCAPTCHA options configuration."},
			{Name: "labels", Type: proto.ColumnType_JSON, Description: "Labels to apply to this security policy."},
			{Name: "label_fingerprint", Type: proto.ColumnType_STRING, Description: "Fingerprint of the labels."},
			{Name: "project", Type: proto.ColumnType_STRING, Transform: transform.FromP(securityPolicyTurbotData, "Project"), Description: "The GCP Project ID."},
			{Name: "location", Type: proto.ColumnType_STRING, Transform: transform.FromConstant("global"), Description: "The location for the resource. Always 'global' for security policies."},
			{Name: "akas", Type: proto.ColumnType_JSON, Transform: transform.FromP(securityPolicyTurbotData, "Akas"), Description: "Array of globally unique identifier strings (also known as 'also known as' or 'akas')."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "Title of the resource."},
		},
	}
}

//// LIST FUNCTION

func listGcpComputeSecurityPolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listGcpComputeSecurityPolicies")

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

	resp := service.SecurityPolicies.List(project).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *compute.SecurityPolicyList) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)
		for _, policy := range page.Items {
			d.StreamListItem(ctx, policy)
			if d.RowsRemaining(ctx) == 0 {
				return nil
			}
		}
		return nil
	}); err != nil {
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
		"Akas":    []string{"gcp://compute.googleapis.com" + data.SelfLink},
	}
	return turbotData[d.Param.(string)], nil
}
