package gcp

import (
	"context"

	resourcemanagerpb "cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"
	"cloud.google.com/go/resourcemanager/apiv3"
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
			Hydrate: listGcpTagBindings,
			KeyColumns: plugin.SingleColumn("parent"),
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
			// Additional fields as needed
			// Standard Steampipe columns
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
		},
	}
}

//// LIST FUNCTION

func listGcpTagBindings(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	parent := d.EqualsQuals["parent"].GetStringValue()
	client, err := resourcemanager.NewTagBindingsClient(ctx)
	if err != nil {
		return nil, err
	}

	// Construct the request
	req := &resourcemanagerpb.ListTagBindingsRequest{
		Parent: parent,
	}
	plugin.Logger(ctx).Warn("listGcpTagBindingsrequest", "req", req)
	// Iterate through the results
	it := client.ListTagBindings(ctx, req)
	for {
		tagBinding, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		plugin.Logger(ctx).Warn("insidelistGcpTagBindingsresponse", tagBinding)
		d.StreamListItem(ctx, tagBinding)

		// Context cancellation check
		if d.RowsRemaining(ctx) == 0 {
			break
		}
	}

	return nil, nil
}
