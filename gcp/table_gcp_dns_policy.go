package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	"google.golang.org/api/dns/v1"
)

//// TABLE DEFINITION

func tableDnsPolicy(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_dns_policy",
		Description: "GCP DNS Policy",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getDnsPolicy,
		},
		List: &plugin.ListConfig{
			Hydrate:           listDnsPolicies,
			ShouldIgnoreError: isIgnorableError([]string{"403"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "An user assigned name for this policy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "An unique identifier for the resource, defined by the server.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "enable_logging",
				Description: "Controls whether logging is enabled for the networks bound to this policy. Defaults to no logging if not set.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "description",
				Description: "A mutable string of at most 1024 characters associated with this resource for the user's convenience. Has no effect on the policy's function.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "enable_inbound_forwarding",
				Description: "Allows networks bound to this policy to receive DNS queries sent by VMs or applications over VPN connections. When enabled, a virtual IP address will be allocated from each of the sub-networks that are bound to this policy.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "kind",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "networks",
				Description: "A list of network names specifying networks to which this policy is applied.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "target_name_servers",
				Description: "Sets an alternative name server for the associated networks. When specified, all DNS queries are forwarded to a name server that you choose. Names such as .internal are not available when an alternative name server is specified.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AlternativeNameServerConfig.TargetNameServers"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDnsPolicyAka,
				Transform:   transform.FromValue(),
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
				Hydrate:     plugin.HydrateFunc(getProject).WithCache(),
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listDnsPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listDnsPolicies")

	// Create Service Connection
	service, err := DnsService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}
	project := projectData.Project

	res := service.Policies.List(project)
	if err := res.Pages(ctx, func(page *dns.PoliciesListResponse) error {
		for _, policy := range page.Policies {
			d.StreamListItem(ctx, policy)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getDnsPolicy(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getDnsPolicy")

	// Create Service Connection
	service, err := DnsService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}
	project := projectData.Project
	name := d.KeyColumnQuals["name"].GetStringValue()

	resp, err := service.Policies.Get(project, name).Do()
	if err != nil {
		return nil, err
	}

	if len(resp.Name) < 1 {
		return nil, nil
	}

	return resp, nil
}

func getDnsPolicyAka(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(*dns.Policy)

	// Get project details
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}
	project := projectData.Project

	akas := []string{"gcp://dns.googleapis.com/projects/" + project + "/policies/" + data.Name}

	return akas, nil
}
