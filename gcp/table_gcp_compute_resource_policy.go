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

func tableGcpComputeResourcePolicy(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_resource_policy",
		Description: "GCP Compute Resource Policy",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeResourcePolicy,
			Tags:       map[string]string{"service": "compute", "action": "resourcePolicies.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeResourcePolicies,
			KeyColumns: plugin.KeyColumnSlice{
				// String columns
				{Name: "status", Require: plugin.Optional, Operators: []string{"<>", "="}},
			},
			Tags: map[string]string{"service": "compute", "action": "resourcePolicies.list"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getComputeResourcePolicyIamPolicy,
				Tags: map[string]string{"service": "compute", "action": "resourcePolicies.getIamPolicy"},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource, provided by the client when initially creating the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique identifier for the resource.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Id"),
			},
			{
				Name:        "status",
				Description: "The status of resource policy creation. Possible values are: 'CREATING', 'DELETING', 'INVALID', and 'READY'.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "A server-defined fully-qualified URL for this resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_timestamp",
				Description: "The date and time, when the policy was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "An user-defined, human-readable description for this resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "Type of the resource. Always compute#resource_policies for resource policies.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "group_placement_policy",
				Description: "Resource policy for instances for placement configuration.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "instance_schedule_policy",
				Description: "Resource policy for scheduling instance operations.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "resource_status",
				Description: "The system status of the resource policy.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "snapshot_schedule_policy",
				Description: "Resource policy for persistent disks for creating snapshots.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "iam_policy",
				Description: "An Identity and Access Management (IAM) policy, which specifies access controls for Google Cloud resources. A `Policy` is a collection of `bindings`. A `binding` binds one or more `members` to a single `role`. Members can be user accounts, service accounts, Google groups, and domains (such as G Suite). A `role` is a named list of permissions; each `role` can be an IAM predefined role or a user-created custom role. For some types of Google Cloud resources, a `binding` can also specify a `condition`, which is a logical expression that allows access to a resource only if the expression evaluates to `true`.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getComputeResourcePolicyIamPolicy,
				Transform:   transform.FromValue(),
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
				Transform:   transform.From(gcpComputeResourcePolicyAkas),
			},

			// GCP standard columns
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
				Transform:   transform.FromP(gcpComputeResourcePolicyAkas, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeResourcePolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeResourcePolicies")

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	filterString := ""
	if d.EqualsQuals["status"] != nil {
		filterString = "status=" + d.EqualsQuals["status"].GetStringValue()
	}

	// Max limit is set as per documentation
	// https://pkg.go.dev/google.golang.org/api@v0.48.0/compute/v1?utm_source=gopls#ResourcePoliciesAggregatedListCall.MaxResults
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
	var resp *compute.ResourcePoliciesAggregatedListCall
	if filterString == "" {
		resp = service.ResourcePolicies.AggregatedList(project).MaxResults(*pageSize)
	} else {
		resp = service.ResourcePolicies.AggregatedList(project).Filter(filterString).MaxResults(*pageSize)
	}
	if err := resp.Pages(
		ctx,
		func(page *compute.ResourcePolicyAggregatedList) error {
			// apply rate limiting
			d.WaitForListRateLimit(ctx)

			for _, item := range page.Items {
				for _, policy := range item.ResourcePolicies {
					d.StreamListItem(ctx, policy)

					// Check if context has been cancelled or if the limit has been hit (if specified)
					// if there is a limit, it will return the number of rows required to reach this limit
					if d.RowsRemaining(ctx) == 0 {
						page.NextPageToken = ""
						return nil
					}
				}
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getComputeResourcePolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getComputeResourcePolicy")

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

	var resourcePolicy compute.ResourcePolicy
	name := d.EqualsQuals["name"].GetStringValue()

	// Return nil, if no input provided
	if name == "" {
		return nil, nil
	}

	resp := service.ResourcePolicies.AggregatedList(project).Filter("name=" + name)
	if err := resp.Pages(
		ctx,
		func(page *compute.ResourcePolicyAggregatedList) error {
			for _, item := range page.Items {
				for _, i := range item.ResourcePolicies {
					resourcePolicy = *i
				}
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	// If the specified resource is not present, API does not return any not found errors
	if len(resourcePolicy.Name) < 1 {
		return nil, nil
	}

	return &resourcePolicy, nil
}

func getComputeResourcePolicyIamPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(*compute.ResourcePolicy)

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	project := strings.Split(data.SelfLink, "/")[6]
	region := getLastPathElement(types.SafeString(data.Region))

	resp, err := service.ResourcePolicies.GetIamPolicy(project, region, data.Name).Do()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

//// TRANSFORM FUNCTIONS

func gcpComputeResourcePolicyAkas(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*compute.ResourcePolicy)

	akas := strings.ReplaceAll(data.SelfLink, "https://www.googleapis.com/compute/v1/", "gcp://compute.googleapis.com/")

	return []string{akas}, nil
}
