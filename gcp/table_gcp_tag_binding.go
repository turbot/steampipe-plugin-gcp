package gcp

import (
	"context"

	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	resourcemanagerpb "cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/iterator"
)

func tableGcpTagBinding(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_tag_binding",
		Description: "GCP Tag Binding",
		List: &plugin.ListConfig{
			Hydrate:           listGcpTagBindings,
			KeyColumns:        plugin.SingleColumn("parent"),
			ShouldIgnoreError: isIgnorableError([]string{"InvalidArgument"}),
			Tags:              map[string]string{"service": "resourcemanager", "action": "hierarchyNodes.listTagBindings"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the TagBinding. This is a string of the form: `tagBindings/{full-resource-name}/{tag-value-name}`.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "parent",
				Description: "The full resource name of the resource the TagValue is bound to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tag_value",
				Description: "The TagValue of the TagBinding. Must be of the form `tagValues/456`.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tag_value_namespaced_name",
				Description: "The namespaced name for the TagValue of the TagBinding.",
				Type:        proto.ColumnType_STRING,
			},
			// Standard Steampipe columns
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			//Standard GCP columns
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Hydrate:     getProject,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listGcpTagBindings(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	parent := d.EqualsQuals["parent"].GetStringValue()

	if parent == "" {
		return nil, nil
	}

	pageSize := types.Int64(300)

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *pageSize {
			pageSize = limit
		}
	}

	// Create Service Connection
	client, err := resourcemanager.NewTagBindingsClient(ctx)
	if err != nil {
		return nil, err
	}

	// Construct the request
	req := &resourcemanagerpb.ListTagBindingsRequest{
		Parent:   parent,
		PageSize: int32(*pageSize),
	}

	// Iterate through the results
	it := client.ListTagBindings(ctx, req)
	for {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		tagBinding, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		d.StreamListItem(ctx, tagBinding)

		// Context cancellation check
		if d.RowsRemaining(ctx) == 0 {
			break
		}
	}

	return nil, nil
}
