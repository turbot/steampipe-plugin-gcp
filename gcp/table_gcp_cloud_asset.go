package gcp

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"

	"google.golang.org/api/cloudasset/v1"
)

//// TABLE DEFINITION

func tableGcpCloudAsset(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_cloud_asset",
		Description: "GCP Cloud Asset",
		List: &plugin.ListConfig{
			Hydrate: listCloudAssets,
			Tags:    map[string]string{"service": "cloudasset", "action": "assets.listResource"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The full name of the asset.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "asset_type",
				Description: "The type of the asset.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "update_time",
				Description: "The last update timestamp of an asset.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "ancestors",
				Description: "The ancestry path of an asset in Google Cloud resource hierarchy.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},

			// Standard GCP columns
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

func listCloudAssets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return listCloudAssetsByContentType(ctx, d, h, "", "gcp_cloud_asset.listCloudAssets", func(ctx context.Context, d *plugin.QueryData, asset *cloudasset.Asset) {
		d.StreamListItem(ctx, asset)
	})
}

// Shared list implementation for all gcp_cloud_asset* tables. The Cloud Asset
// API returns different asset content depending on the requested content type,
// and accepts at most one content type per call. An empty contentType requests
// only basic asset information (name, asset_type, ancestors, update_time).
func listCloudAssetsByContentType(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, contentType string, logCtx string, stream func(ctx context.Context, d *plugin.QueryData, asset *cloudasset.Asset)) (interface{}, error) {

	// Create Service Connection
	service, err := CloudAssetService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error(logCtx, "service_error", err)
		return nil, err
	}

	// Max limit is set as per documentation
	pageSize := types.Int64(1000)
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

	input := "projects/" + project

	resp := service.Assets.List(input).PageSize(*pageSize)
	if contentType != "" {
		resp = resp.ContentType(contentType)
	}

	if err := resp.Pages(ctx, func(page *cloudasset.ListAssetsResponse) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, item := range page.Assets {
			stream(ctx, d, item)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				page.NextPageToken = ""
				return nil
			}
		}
		return nil
	}); err != nil {
		plugin.Logger(ctx).Error(logCtx, "api_error", err)
		return nil, err
	}

	return nil, nil
}
