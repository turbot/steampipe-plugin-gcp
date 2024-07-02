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

func tableGcpComputeFirewall(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_firewall",
		Description: "GCP Compute Firewall",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeFirewall,
			Tags:       map[string]string{"service": "compute", "action": "firewalls.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeFirewalls,
			KeyColumns: plugin.KeyColumnSlice{
				// String columns
				{Name: "direction", Require: plugin.Optional, Operators: []string{"<>", "="}},

				// Boolean columns
				{Name: "disabled", Require: plugin.Optional, Operators: []string{"<>", "="}},
			},
			Tags: map[string]string{"service": "compute", "action": "firewalls.list"},
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
				Transform:   transform.FromP(gcpComputeFirewallTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeFirewalls(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeFirewalls")
	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	filterQuals := []filterQualMap{
		{"direction", "direction", "string"},
		{"disabled", "disabled", "boolean"},
	}

	filters := buildQueryFilterFromQuals(filterQuals, d.Quals)
	filterString := ""
	if len(filters) > 0 {
		filterString = strings.Join(filters, " ")
	}

	// Max limit is set as per documentation
	// https://pkg.go.dev/google.golang.org/api@v0.48.0/compute/v1?utm_source=gopls#FirewallsListCall.MaxResults
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

	resp := service.Firewalls.List(project).Filter(filterString).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *compute.FirewallList) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, firewall := range page.Items {
			d.StreamListItem(ctx, firewall)

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

func getComputeFirewall(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

	project := strings.Split(firewall.SelfLink, "/")[6]

	turbotData := map[string]interface{}{
		"Action":  action,
		"Project": project,
		"Akas":    []string{"gcp://compute.googleapis.com/projects/" + project + "/global/firewalls/" + firewall.Name},
	}

	return turbotData[param], nil
}
