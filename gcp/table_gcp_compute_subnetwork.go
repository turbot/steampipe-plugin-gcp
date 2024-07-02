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

func tableGcpComputeSubnetwork(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_subnetwork",
		Description: "GCP Compute Subnetwork",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeSubnetwork,
			Tags:       map[string]string{"service": "compute", "action": "subnetworks.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeSubnetworks,
			KeyColumns: plugin.KeyColumnSlice{
				// String columns
				{Name: "state", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "private_ipv6_google_access", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "purpose", Require: plugin.Optional, Operators: []string{"<>", "="}},

				// Boolean columns
				{Name: "enable_flow_logs", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "private_ip_google_access", Require: plugin.Optional, Operators: []string{"<>", "="}},
			},
			Tags: map[string]string{"service": "compute", "action": "subnetworks.list"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getComputeSubnetworkIamPolicy,
				Tags: map[string]string{"service": "compute", "action": "subnetworks.getIamPolicy"},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "Name of the resource. Provided by the client when the resource is created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique identifier for the resource. This identifier is defined by the server.",
				Type:        proto.ColumnType_INT,
			},
			// For simplified view of network, without including the full path url of the network
			{
				Name:        "network_name",
				Description: "The name of the network to which this subnetwork belongs.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Network").Transform(lastPathElement),
			},
			{
				Name:        "state",
				Description: "Specifies the current state of the subnetwork.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "Type of the resource. Always compute#subnetwork for Subnetwork resources.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A user-specified, human-readable description of the address.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_timestamp",
				Description: "The creation timestamp of the resource.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "enable_flow_logs",
				Description: "Specifies whether to enable flow logging for this subnetwork, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "fingerprint",
				Description: "An unique system generated string, to reduce conflicts when multiple users change any property of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "log_config_aggregation_interval",
				Description: "Can only be specified if VPC flow logging for this subnetwork is enabled. Toggles the aggregation interval for collecting flow logs. Increasing the interval time will reduce the amount of generated flow logs for long lasting connections.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LogConfig.AggregationInterval"),
			},
			{
				Name:        "log_config_enable",
				Description: "Specifies whether to enable flow logging for this subnetwork, or not.",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
				Transform:   transform.FromField("LogConfig.Enable"),
			},
			{
				Name:        "log_config_filter_expr",
				Description: "Can only be specified if VPC flow logs for this subnetwork is enabled. Export filter used to define which VPC flow logs should be logged.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LogConfig.FilterExpr"),
			},
			{
				Name:        "log_config_flow_sampling",
				Description: "Can only be specified if VPC flow logging for this subnetwork is enabled. The value of the field must be in [0, 1]. Set the sampling rate of VPC flow logs within the subnetwork where 1.0 means all collected logs are reported and 0.0 means no logs are reported. Default is 0.5, which means half of all collected logs are reported.",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("LogConfig.FlowSampling"),
			},
			{
				Name:        "log_config_metadata",
				Description: "Configures whether all, none or a subset of metadata fields should be added to the reported VPC flow logs.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LogConfig.Metadata"),
			},
			{
				Name:        "gateway_address",
				Description: "The gateway address for default routes to reach destination addresses outside this subnetwork.",
				Type:        proto.ColumnType_IPADDR,
			},
			{
				Name:        "ip_cidr_range",
				Description: "The range of internal addresses that are owned by this subnetwork.",
				Type:        proto.ColumnType_CIDR,
			},
			{
				Name:        "ipv6_cidr_range",
				Description: "The range of internal IPv6 addresses that are owned by this subnetwork.",
				Type:        proto.ColumnType_CIDR,
			},
			{
				Name:        "network",
				Description: "The URL of the network to which this subnetwork belongs.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "private_ip_google_access",
				Description: "Specifies whether the VMs in this subnet can access Google services without assigned external IP addresses.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "private_ipv6_google_access",
				Description: "The private IPv6 google access type for the VMs in this subnet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "purpose",
				Description: "The purpose of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "region",
				Description: "The URL of the region where the Subnetwork resides.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "role",
				Description: "Specifies the role of the subnetwork.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "log_config_metadata_fields",
				Description: "Can only be specified if VPC flow logs for this subnetwork is enabled and 'metadata' was set to CUSTOM_METADATA.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("LogConfig.MetadataFields"),
			},
			{
				Name:        "secondary_ip_ranges",
				Description: "An array of configurations for secondary IP ranges for VM instances contained in this subnetwork.",
				Type:        proto.ColumnType_JSON,
			},

			{
				Name:        "iam_policy",
				Description: "An Identity and Access Management (IAM) policy, which specifies access controls for Google Cloud resources. A `Policy` is a collection of `bindings`. A `binding` binds one or more `members` to a single `role`. Members can be user accounts, service accounts, Google groups, and domains (such as G Suite). A `role` is a named list of permissions; each `role` can be an IAM predefined role or a user-created custom role. For some types of Google Cloud resources, a `binding` can also specify a `condition`, which is a logical expression that allows access to a resource only if the expression evaluates to `true`.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getComputeSubnetworkIamPolicy,
				Transform:   transform.FromValue(),
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
				Transform:   transform.FromP(gcpComputeSubnetworkTurbotData, "Akas"),
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
				Transform:   transform.FromP(gcpComputeSubnetworkTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeSubnetworks(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeSubnetworks")
	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	filterQuals := []filterQualMap{
		{"state", "state", "string"},
		{"private_ipv6_google_access", "privateIpv6GoogleAccess", "string"},
		{"purpose", "purpose", "string"},
		{"enable_flow_logs", "enableFlowLogs", "boolean"},
		{"private_ip_google_access", "privateIpGoogleAccess", "boolean"},
	}

	filters := buildQueryFilterFromQuals(filterQuals, d.Quals)
	filterString := ""
	if len(filters) > 0 {
		filterString = strings.Join(filters, " ")
	}

	// Max limit is set as per documentation
	// https://pkg.go.dev/google.golang.org/api@v0.48.0/compute/v1?utm_source=gopls#SubnetworksAggregatedListCall.MaxResults
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

	resp := service.Subnetworks.AggregatedList(project).Filter(filterString).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *compute.SubnetworkAggregatedList) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, item := range page.Items {
			for _, subnetwork := range item.Subnetworks {
				d.StreamListItem(ctx, subnetwork)

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

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getComputeSubnetwork(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getComputeSubnetwork")

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

	var subnetwork compute.Subnetwork
	name := d.EqualsQuals["name"].GetStringValue()

	resp := service.Subnetworks.AggregatedList(project).Filter("name=" + name)
	if err := resp.Pages(ctx, func(page *compute.SubnetworkAggregatedList) error {
		for _, item := range page.Items {
			for _, i := range item.Subnetworks {
				subnetwork = *i
			}
		}
		return nil
	},
	); err != nil {
		return nil, err
	}

	// If the specified resource is not present, API does not return any not found errors
	if len(subnetwork.Name) < 1 {
		return nil, nil
	}

	return &subnetwork, nil
}

func getComputeSubnetworkIamPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	subnetwork := h.Item.(*compute.Subnetwork)

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	var resp *compute.Policy
	project := strings.Split(subnetwork.SelfLink, "/")[6]
	regionName := getLastPathElement(types.SafeString(subnetwork.Region))

	resp, err = service.Subnetworks.GetIamPolicy(project, regionName, subnetwork.Name).Do()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

//// TRANSFORM FUNCTIONS

func gcpComputeSubnetworkTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	subnetwork := d.HydrateItem.(*compute.Subnetwork)
	param := d.Param.(string)

	region := getLastPathElement(types.SafeString(subnetwork.Region))
	project := strings.Split(subnetwork.SelfLink, "/")[6]

	turbotData := map[string]interface{}{
		"Project": project,
		"Akas":    []string{"gcp://compute.googleapis.com/projects/" + project + "/regions/" + region + "/subnetworks/" + subnetwork.Name},
	}

	return turbotData[param], nil
}
