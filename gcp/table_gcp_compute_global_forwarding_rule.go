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

func tableGcpComputeGlobalForwardingRule(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_global_forwarding_rule",
		Description: "GCP Compute Global Forwarding Rule",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeGlobalForwardingRule,
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeGlobalForwardingRules,
			KeyColumns: plugin.KeyColumnSlice{
				// String columns
				{Name: "ip_protocol", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "ip_version", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "load_balancing_scheme", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "network_tier", Require: plugin.Optional, Operators: []string{"<>", "="}},

				// Boolean columns
				{Name: "allow_global_access", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "all_ports", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "is_mirroring_collector", Require: plugin.Optional, Operators: []string{"<>", "="}},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "A friendly name that identifies the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique identifier for the resource.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "description",
				Description: "A user-specified, human-readable description of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ip_address",
				Description: "Specifies the IP address that this forwarding rule serves.",
				Type:        proto.ColumnType_INET,
				Transform:   transform.FromField("IPAddress"),
			},
			{
				Name:        "allow_global_access",
				Description: "Specifies whether clients can access ILB from all regions, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "kind",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "The server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_timestamp",
				Description: "The creation timestamp of the resource.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "all_ports",
				Description: "Specify this field to allow packets addressed to any ports will be forwarded to the backends configured with this forwarding rule.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "backend_service",
				Description: "Specifies the BackendService resource to receive the matched traffic.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "fingerprint",
				Description: "a hash of the contents stored in this object and used for optimistic locking.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ip_protocol",
				Description: "The IP protocol to which this rule applies.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("IPProtocol"),
			},
			{
				Name:        "ip_version",
				Description: "The IP Version that will be used by this forwarding rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_mirroring_collector",
				Description: "Indicates whether or not this load balancer can be used as a collector for packet mirroring.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "load_balancing_scheme",
				Description: "Specifies the forwarding rule type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "network",
				Description: "Specifies the network that the load balanced IP should belong to for this Forwarding Rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "network_tier",
				Description: "Specifies tthe networking tier used for configuring this load balancer.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "port_range",
				Description: "Specifies the port range. Packets addressed to ports in the specified range will be forwarded to target or backendService.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_label",
				Description: "A prefix to the service name for this Forwarding Rule. If specified, the prefix is the first label of the fully qualified service name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_name",
				Description: "The internal fully qualified service name for this Forwarding Rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "subnetwork",
				Description: "Specifies the subnetwork that the load balanced IP should belong to for this Forwarding Rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "target",
				Description: "The URL of the target resource to receive the matched traffic.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "metadata_filters",
				Description: "Opaque filter criteria used by Loadbalancer to restrict routing configuration to a limited set of xDS compliant clients.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "ports",
				Description: "A list of ports can be configured.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "labels",
				Description: "A list of labels attached to this resource.",
				Type:        proto.ColumnType_JSON,
			},

			// standard steampipe columns
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
				Transform:   transform.FromP(globalForwardingRuleSelfLinkToTurbotData, "Akas"),
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
				Transform:   transform.FromP(globalForwardingRuleSelfLinkToTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeGlobalForwardingRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	filterQuals := []filterQualMap{
		{"ip_protocol", "ipProtocol", "string"},
		{"ip_version", "ipVersion", "string"},
		{"load_balancing_scheme", "loadBalancingScheme", "string"},
		{"network_tier", "networkTier", "string"},
		{"allow_global_access", "allowGlobalAccess", "boolean"},
		{"all_ports", "allPorts", "boolean"},
		{"is_mirroring_collector", "isMirroringCollector", "boolean"},
	}

	filters := buildQueryFilterFromQuals(filterQuals, d.Quals)
	filterString := ""
	if len(filters) > 0 {
		filterString = strings.Join(filters, " ")
	}

	// Max limit is set as per documentation
	// https://pkg.go.dev/google.golang.org/api@v0.48.0/compute/v1?utm_source=gopls#GlobalForwardingRulesListCall.MaxResults
	pageSize := types.Int64(500)
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

	resp := service.GlobalForwardingRules.List(project).Filter(filterString).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *compute.ForwardingRuleList) error {
		for _, globalForwardingRule := range page.Items {
			d.StreamListItem(ctx, globalForwardingRule)

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

	return nil, err
}

//// HYDRATE FUNCTIONS

func getComputeGlobalForwardingRule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

	name := d.EqualsQuals["name"].GetStringValue()

	req, err := service.GlobalForwardingRules.Get(project, name).Do()
	if err != nil {
		return nil, err
	}

	return req, nil
}

//// TRANSFORM FUNCTIONS

func globalForwardingRuleSelfLinkToTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	globalForwardingRule := d.HydrateItem.(*compute.ForwardingRule)
	param := d.Param.(string)

	project := strings.Split(globalForwardingRule.SelfLink, "/")[6]

	turbotData := map[string]interface{}{
		"Project": project,
		"Akas":    []string{"gcp://compute.googleapis.com/projects/" + project + "/global/forwardingRules/" + globalForwardingRule.Name},
	}

	return turbotData[param], nil
}
