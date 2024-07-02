package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	compute "google.golang.org/api/compute/v0.beta"
)

//// TABLE DEFINITION

func tableGcpComputeForwardingRule(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_forwarding_rule",
		Description: "GCP Compute Forwarding Rule",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeForwardingRule,
			Tags:       map[string]string{"service": "compute", "action": "forwardingRules.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeForwardingRules,
			KeyColumns: plugin.KeyColumnSlice{
				// String columns
				{Name: "ip_protocol", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "load_balancing_scheme", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "network_tier", Require: plugin.Optional, Operators: []string{"<>", "="}},

				// Boolean columns
				{Name: "allow_global_access", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "all_ports", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "is_mirroring_collector", Require: plugin.Optional, Operators: []string{"<>", "="}},
			},
			Tags: map[string]string{"service": "compute", "action": "forwardingRules.list"},
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
				Name:        "region",
				Description: "The URL of the region where the regional forwarding rule resides.",
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
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Labels"),
			},
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
				Transform:   transform.FromP(forwardingRuleSelfLinkToTurbotData, "Akas"),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Region").Transform(lastPathElement),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(forwardingRuleSelfLinkToTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeForwardingRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create Service Connection
	service, err := ComputeBetaService(ctx, d)
	if err != nil {
		return nil, err
	}

	filterQuals := []filterQualMap{
		{"ip_protocol", "ipProtocol", "string"},
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
	// https://pkg.go.dev/google.golang.org/api@v0.48.0/compute/v0.beta?utm_source=gopls#ForwardingRulesAggregatedListCall.MaxResults
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

	resp := service.ForwardingRules.AggregatedList(project).Filter(filterString).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *compute.ForwardingRuleAggregatedList) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, item := range page.Items {
			for _, forwardingRule := range item.ForwardingRules {
				d.StreamListItem(ctx, forwardingRule)

				// Check if context has been cancelled or if the limit has been hit (if specified)
				// if there is a limit, it will return the number of rows required to reach this limit
				if d.RowsRemaining(ctx) == 0 {
					page.NextPageToken = ""
					return nil
				}
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getComputeForwardingRule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := ComputeBetaService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	var forwardingRule compute.ForwardingRule
	name := d.EqualsQuals["name"].GetStringValue()

	resp := service.ForwardingRules.AggregatedList(project).Filter("name=" + name)
	if err := resp.Pages(
		ctx,
		func(page *compute.ForwardingRuleAggregatedList) error {
			for _, item := range page.Items {
				for _, i := range item.ForwardingRules {
					forwardingRule = *i
				}
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	// If the specified resource is not present, API does not return any not found errors
	if len(forwardingRule.Name) < 1 {
		return nil, nil
	}

	return &forwardingRule, nil
}

//// TRANSFORM FUNCTIONS

func forwardingRuleSelfLinkToTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	forwardingRule := d.HydrateItem.(*compute.ForwardingRule)
	param := d.Param.(string)

	project := strings.Split(forwardingRule.SelfLink, "/")[6]
	region := getLastPathElement(types.SafeString(forwardingRule.Region))

	turbotData := map[string]interface{}{
		"Project": project,
		"Akas":    []string{"gcp://compute.googleapis.com/projects/" + project + "/regions/" + region + "/forwardingRules/" + forwardingRule.Name},
	}

	return turbotData[param], nil
}
