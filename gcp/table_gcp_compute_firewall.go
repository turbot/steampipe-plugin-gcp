package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"google.golang.org/api/compute/v1"
)

//// TABLE DEFINITION

func tableGcpComputeFirewall(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_firewall",
		Description: "GCP Compute Firewall",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeFirewall,
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeFirewalls,
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
				Name:        "direction",
				Description: "Direction of traffic to which this firewall applies.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "Specifies the type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "disabled",
				Description: "Indicates whether the firewall rule is disabled, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "description",
				Description: "A user-specified, human-readable description of the firewall.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_timestamp",
				Description: "The creation timestamp of the resource.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "action",
				Description: "Describes the type action specified by the rule.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpComputeFirewallTurbotData, "Action"),
			},
			{
				Name:        "log_config_enable",
				Description: "Specifies whether to enable logging for a particular firewall rule, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("LogConfig.Enable"),
			},
			{
				Name:        "log_config_metadata",
				Description: "Specifies whether to include or exclude metadata for firewall logs.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LogConfig.Metadata"),
			},
			{
				Name:        "network",
				Description: "The URL of the network resource for this firewall rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "priority",
				Description: "Specifies the priority for this rule. Relative priorities determine which rule takes effect if multiple rules apply. Lower values indicate higher priority.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "self_link",
				Description: "The server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "allowed",
				Description: "The list of ALLOW rules specified by this firewall.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "denied",
				Description: "The list of DENY rules specified by this firewall.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "destination_ranges",
				Description: "A list of CIDR ranges. The firewall rule applies only to traffic that has destination IP address in these ranges.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "source_ranges",
				Description: "A list of CIDR ranges. The firewall rule applies only to traffic originating from an instance with a service account in this list.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "source_service_accounts",
				Description: "A list of service account. The firewall rule applies only to traffic that has a source IP address in these ranges.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "source_tags",
				Description: "A list of tags. The firewall rule applies only to traffic with source IPs that match the primary network interfaces of VM instances that have the tag and are in the same VPC network.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "target_service_accounts",
				Description: "A list of service accounts indicating sets of instances located in the network that may make network connections as specified in Allowed",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "target_tags",
				Description: "A list of tags that controls which instances the firewall rule applies to.",
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
				Transform:   transform.FromP(gcpComputeFirewallTurbotData, "Akas"),
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
				Transform:   transform.FromConstant(activeProject()),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeFirewalls(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeFirewalls")
	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	project := activeProject()
	resp := service.Firewalls.List(project)
	if err := resp.Pages(ctx, func(page *compute.FirewallList) error {
		for _, firewall := range page.Items {
			d.StreamListItem(ctx, firewall)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getComputeFirewall(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	name := d.KeyColumnQuals["name"].GetStringValue()
	project := activeProject()

	// Error: pq: rpc error: code = Unknown desc = json: invalid use of ,string struct tag,
	// trying to unmarshal "projects/project/global/firewalls/" into uint64
	if len(name) < 1 {
		return nil, nil
	}

	req, err := service.Firewalls.Get(project, name).Do()
	if err != nil {
		return nil, err
	}

	return req, nil
}

//// TRANSFORM FUNCTIONS

func gcpComputeFirewallTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	firewall := d.HydrateItem.(*compute.Firewall)
	param := d.Param.(string)

	var action string
	if firewall.Allowed != nil {
		action = "Allow"
	}
	if firewall.Denied != nil {
		action = "Deny"
	}

	turbotData := map[string]interface{}{
		"Action": action,
		"Akas":   []string{"gcp://compute.googleapis.com/projects/" + activeProject() + "/global/firewalls/" + firewall.Name},
	}

	return turbotData[param], nil
}
