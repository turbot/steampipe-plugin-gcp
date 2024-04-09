package gcp

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
			Tags:       map[string]string{"service": "dns", "action": "policies.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listDnsPolicies,
			Tags:    map[string]string{"service": "dns", "action": "policies.list"},
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
				Hydrate:     getProject,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listDnsPolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listDnsPolicies")

	// Create Service Connection
	service, err := DnsService(ctx, d)
	if err != nil {
		return nil, err
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

	res := service.Policies.List(project).MaxResults(*pageSize)
	if err := res.Pages(ctx, func(page *dns.PoliciesListResponse) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, policy := range page.Policies {
			d.StreamListItem(ctx, policy)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				page.NextPageToken = ""
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

func getDnsPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getDnsPolicy")

	// Create Service Connection
	service, err := DnsService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)
	name := d.EqualsQuals["name"].GetStringValue()

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

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	akas := []string{"gcp://dns.googleapis.com/projects/" + project + "/policies/" + data.Name}

	return akas, nil
}
