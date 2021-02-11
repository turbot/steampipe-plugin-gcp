package gcp

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

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
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeNodeTemplates,
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
				Transform:   transform.From(computeNodeTemplateAka),
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

func listComputeNodeTemplates(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeNodeTemplates")
	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	project := projectName
	resp := service.NodeTemplates.AggregatedList(project)
	if err := resp.Pages(ctx, func(page *compute.NodeTemplateAggregatedList) error {
		for _, item := range page.Items {
			for _, nodeTemplate := range item.NodeTemplates {
				d.StreamListItem(ctx, nodeTemplate)
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
	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	var nodeTemplate compute.NodeTemplate
	name := d.KeyColumnQuals["name"].GetStringValue()
	project := projectName

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
	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	nodeTemplate := h.Item.(*compute.NodeTemplate)
	regionName := getLastPathElement(types.SafeString(nodeTemplate.Region))

	req, err := service.NodeTemplates.GetIamPolicy(activeProject(), regionName, nodeTemplate.Name).Do()
	if err != nil {
		return nil, err
	}

	return req, nil
}

//// TRANSFORM FUNCTIONS

func computeNodeTemplateAka(_ context.Context, d *transform.TransformData) (interface{}, error) {
	nodeTemplate := d.HydrateItem.(*compute.NodeTemplate)
	regionName := getLastPathElement(types.SafeString(nodeTemplate.Region))

	akas := []string{"gcp://compute.googleapis.com/projects/" + activeProject() + "/regions/" + regionName + "/nodeTemplates/" + nodeTemplate.Name}

	return akas, nil
}
