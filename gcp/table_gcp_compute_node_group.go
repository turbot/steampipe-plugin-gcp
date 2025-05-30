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

func tableGcpComputeNodeGroup(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_node_group",
		Description: "GCP Compute Node Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeNodeGroup,
			Tags:       map[string]string{"service": "compute", "action": "nodeGroups.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeNodeGroups,
			KeyColumns: plugin.KeyColumnSlice{
				// String columns
				{Name: "status", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "maintenance_policy", Require: plugin.Optional, Operators: []string{"<>", "="}},
			},
			Tags: map[string]string{"service": "compute", "action": "nodeGroups.list"},
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
				Name:        "status",
				Description: "Specifies the current state of the node group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "size",
				Description: "The total number of nodes in the node group.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "description",
				Description: "A user-specified, human-readable description of the node group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "autoscaling_policy_mode",
				Description: "Specifies the autoscaling mode of the node group. Set to one of: ON, OFF, or ONLY_SCALE_OUT.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AutoscalingPolicy.Mode"),
			},
			{
				Name:        "autoscaling_policy_max_nodes",
				Description: "The maximum number of nodes that the group should have. Must be set if autoscaling is enabled. Maximum value allowed is 100.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("AutoscalingPolicy.MaxNodes"),
			},
			{
				Name:        "autoscaling_policy_min_nodes",
				Description: "The minimum number of nodes that the group should have.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("AutoscalingPolicy.MinNodes"),
			},
			{
				Name:        "creation_timestamp",
				Description: "The creation timestamp of the resource.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "fingerprint",
				Description: "An unique system generated string, to reduce conflicts when multiple users change any property of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "maintenance_policy",
				Description: "Specifies how to handle instances when a node in the group undergoes maintenance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "node_template",
				Description: "The URL of the node template to create the node group from.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "The server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "zone",
				Description: "The name of the zone where the node group resides.",
				Type:        proto.ColumnType_STRING,
			},
			// zone_name is a simpler view of the zone, without the full path
			{
				Name:        "zone_name",
				Description: "The zone name in which the node group resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Zone").Transform(lastPathElement),
			},
			{
				Name:        "iam_policy",
				Description: "An Identity and Access Management (IAM) policy, which specifies access controls for Google Cloud resources. A `Policy` is a collection of `bindings`. A `binding` binds one or more `members` to a single `role`. Members can be user accounts, service accounts, Google groups, and domains (such as G Suite). A `role` is a named list of permissions; each `role` can be an IAM predefined role or a user-created custom role. For some types of Google Cloud resources, a `binding` can also specify a `condition`, which is a logical expression that allows access to a resource only if the expression evaluates to `true`.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getComputeNodeGroupIamPolicy,
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
				Transform:   transform.FromP(gcpComputeNodeGroupTurbotData, "Akas"),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Zone").Transform(lastPathElement),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpComputeNodeGroupTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeNodeGroups(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeNodeGroups")
	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	filterQuals := []filterQualMap{
		{"status", "status", "string"},
		{"maintenance_policy", "maintenancePolicy", "string"},
	}

	filters := buildQueryFilterFromQuals(filterQuals, d.Quals)
	filterString := ""
	if len(filters) > 0 {
		filterString = strings.Join(filters, " ")
	}

	// Max limit is set as per documentation
	// https://pkg.go.dev/google.golang.org/api@v0.48.0/compute/v1?utm_source=gopls#NodeGroupsAggregatedListCall.MaxResults
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

	resp := service.NodeGroups.AggregatedList(project).Filter(filterString).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *compute.NodeGroupAggregatedList) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, item := range page.Items {
			for _, nodeGroup := range item.NodeGroups {
				d.StreamListItem(ctx, nodeGroup)

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

func getComputeNodeGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

	var nodeGroup compute.NodeGroup
	name := d.EqualsQuals["name"].GetStringValue()

	resp := service.NodeGroups.AggregatedList(project).Filter("name=" + name)
	if err := resp.Pages(
		ctx,
		func(page *compute.NodeGroupAggregatedList) error {
			for _, item := range page.Items {
				for _, i := range item.NodeGroups {
					nodeGroup = *i
				}
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	// If the specified resource is not present, API does not return any not found errors
	if len(nodeGroup.Name) < 1 {
		return nil, nil
	}

	return &nodeGroup, nil
}

func getComputeNodeGroupIamPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	nodeGroup := h.Item.(*compute.NodeGroup)
	project := strings.Split(nodeGroup.SelfLink, "/")[6]
	zoneName := getLastPathElement(types.SafeString(nodeGroup.Zone))

	req, err := service.NodeGroups.GetIamPolicy(project, zoneName, nodeGroup.Name).Do()
	if err != nil {
		// Return nil, if the resource not present
		result := isIgnorableError([]string{"404"})
		if result != nil {
			return nil, nil
		}
		return nil, err
	}

	return req, nil
}

//// TRANSFORM FUNCTIONS

func gcpComputeNodeGroupTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	nodeGroup := d.HydrateItem.(*compute.NodeGroup)
	param := d.Param.(string)

	project := strings.Split(nodeGroup.SelfLink, "/")[6]
	zone := getLastPathElement(types.SafeString(nodeGroup.Zone))

	turbotData := map[string]interface{}{
		"Project": project,
		"Akas":    []string{"gcp://compute.googleapis.com/projects/" + project + "/zones/" + zone + "/nodeGroups/" + nodeGroup.Name},
	}

	return turbotData[param], nil
}
