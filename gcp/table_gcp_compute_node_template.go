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

func tableGcpComputeNodeTemplate(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_node_template",
		Description: "GCP Compute Node Template",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeNodeTemplate,
			Tags:       map[string]string{"service": "compute", "action": "nodeTemplates.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeNodeTemplates,
			KeyColumns: plugin.KeyColumnSlice{
				// String columns
				{Name: "cpu_overcommit_type", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "node_type", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "status", Require: plugin.Optional, Operators: []string{"<>", "="}},
			},
			Tags: map[string]string{"service": "compute", "action": "nodeTemplates.list"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getComputeNodeGroupIamPolicy,
				Tags: map[string]string{"service": "compute", "action": "nodeTemplates.getIamPolicy"},
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
				Name:        "status",
				Description: "Specifies the status of the node template.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A user-specified, human-readable description of the node template.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_timestamp",
				Description: "The creation timestamp of the resource.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "cpu_overcommit_type",
				Description: "Specifies the CPU overcommit.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "node_type",
				Description: "Specifies the type of the nodes to use for node groups, that are created from this template.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "node_type_flexibility_cpus",
				Description: "The URL of the network in which to reserve the address.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("NodeTypeFlexibility.Cpus"),
			},
			{
				Name:        "node_type_flexibility_local_ssd",
				Description: "Specifies the networking tier used for configuring this address.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("NodeTypeFlexibility.LocalSsd"),
			},
			{
				Name:        "node_type_flexibility_memory",
				Description: "Specifies the prefix length if the resource represents an IP range.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("NodeTypeFlexibility.Memory"),
			},
			{
				Name:        "region",
				Description: "The name of the region where the node template resides.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "The server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "server_binding_type",
				Description: "Specifies the binding properties for the physical server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerBinding.Type"),
			},
			{
				Name:        "status_message",
				Description: "A human-readable explanation of the resource status.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "node_affinity_labels",
				Description: "A list of labels to use for node affinity, which will be used in instance scheduling.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "iam_policy",
				Description: "An Identity and Access Management (IAM) policy, which specifies access controls for Google Cloud resources. A `Policy` is a collection of `bindings`. A `binding` binds one or more `members` to a single `role`. Members can be user accounts, service accounts, Google groups, and domains (such as G Suite). A `role` is a named list of permissions; each `role` can be an IAM predefined role or a user-created custom role. For some types of Google Cloud resources, a `binding` can also specify a `condition`, which is a logical expression that allows access to a resource only if the expression evaluates to `true`.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getComputeNodeTemplateIamPolicy,
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
				Transform:   transform.FromP(gcpComputeNodeTemplateTurbotData, "Akas"),
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
				Transform:   transform.FromP(gcpComputeNodeTemplateTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeNodeTemplates(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeNodeTemplates")

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	filterQuals := []filterQualMap{
		{"cpu_overcommit_type", "cpuOvercommitType", "string"},
		{"node_type", "nodeType", "string"},
		{"status", "status", "string"},
	}

	filters := buildQueryFilterFromQuals(filterQuals, d.Quals)
	filterString := ""
	if len(filters) > 0 {
		filterString = strings.Join(filters, " ")
	}

	// Max limit is set as per documentation
	// https://pkg.go.dev/google.golang.org/api@v0.48.0/compute/v1?utm_source=gopls#NodeTemplatesAggregatedListCall.MaxResults
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

	resp := service.NodeTemplates.AggregatedList(project).Filter(filterString).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *compute.NodeTemplateAggregatedList) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, item := range page.Items {
			for _, nodeTemplate := range item.NodeTemplates {
				d.StreamListItem(ctx, nodeTemplate)

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

func getComputeNodeTemplate(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

	var nodeTemplate compute.NodeTemplate
	name := d.EqualsQuals["name"].GetStringValue()

	resp := service.NodeTemplates.AggregatedList(project).Filter("name=" + name)
	if err := resp.Pages(
		ctx,
		func(page *compute.NodeTemplateAggregatedList) error {
			for _, item := range page.Items {
				for _, i := range item.NodeTemplates {
					nodeTemplate = *i
				}
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	// If the specified resource is not present, API does not return any not found errors
	if len(nodeTemplate.Name) < 1 {
		return nil, nil
	}

	return &nodeTemplate, nil
}

func getComputeNodeTemplateIamPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	nodeTemplate := h.Item.(*compute.NodeTemplate)
	regionName := getLastPathElement(types.SafeString(nodeTemplate.Region))
	project := strings.Split(nodeTemplate.SelfLink, "/")[6]

	req, err := service.NodeTemplates.GetIamPolicy(project, regionName, nodeTemplate.Name).Do()
	if err != nil {
		return nil, err
	}

	return req, nil
}

//// TRANSFORM FUNCTIONS

func gcpComputeNodeTemplateTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	nodeTemplate := d.HydrateItem.(*compute.NodeTemplate)
	param := d.Param.(string)

	project := strings.Split(nodeTemplate.SelfLink, "/")[6]
	region := getLastPathElement(types.SafeString(nodeTemplate.Region))

	turbotData := map[string]interface{}{
		"Project": project,
		"Akas":    []string{"gcp://compute.googleapis.com/projects/" + project + "/regions/" + region + "/nodeTemplates/" + nodeTemplate.Name},
	}

	return turbotData[param], nil
}
