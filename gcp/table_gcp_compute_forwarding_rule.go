package gcp

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

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
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeForwardingRules,
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
				Type:        proto.ColumnType_IPADDR,
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
				Transform:   transform.From(forwardingRuleAka),
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
				Transform:   transform.FromConstant(projectName),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeForwardingRules(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	project := projectName
	resp := service.ForwardingRules.AggregatedList(project)
	if err := resp.Pages(ctx, func(page *compute.ForwardingRuleAggregatedList) error {
		for _, item := range page.Items {
			for _, forwardingRule := range item.ForwardingRules {
				d.StreamListItem(ctx, forwardingRule)
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
	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	var forwardingRule compute.ForwardingRule
	name := d.KeyColumnQuals["name"].GetStringValue()
	project := projectName

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

func forwardingRuleAka(_ context.Context, d *transform.TransformData) (interface{}, error) {
	forwardingRule := d.HydrateItem.(*compute.ForwardingRule)
	regionName := getLastPathElement(types.SafeString(forwardingRule.Region))

	akas := []string{"gcp://compute.googleapis.com/projects/" + activeProject() + "/regions/" + regionName + "/forwardingRules/" + forwardingRule.Name}

	return akas, nil
}
