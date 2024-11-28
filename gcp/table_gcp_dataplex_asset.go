package gcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"

	"google.golang.org/api/dataplex/v1"
)

func tableGcpDataplexAsset(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_dataplex_asset",
		Description: "GCP Dataplex Asset",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getDataplexAsset,
		},
		List: &plugin.ListConfig{
			Hydrate:       listDataplexAssets,
			KeyColumns: plugin.KeyColumnSlice{
				// The 'zone_name' is required to query this table. Since Steampipe does not yet support parent hydrate chaining, the Lake is used as the parent of Zone, which means Zone cannot be used as the parent of Asset.
				{Name: "zone_name", Require: plugin.Required, Operators: []string{"="}, CacheMatch: query_cache.CacheMatchExact},
				{Name: "display_name", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "state", Require: plugin.Optional, Operators: []string{"="}},
			},
			ShouldIgnoreError: isIgnorableError([]string{"404"}),
		},
		// Build matrix region is not required, because we must have to pass the zone_name or name to get the assets.
		Columns: []*plugin.Column{
			{
				Name:        "display_name",
				Description: "User friendly display name.",
				Type:        proto.ColumnType_STRING,
				Transform: transform.FromField("DisplayName"),
			},
			{
				Name:        "name",
				Description: "The relative resource name of the asset.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lake_name",
				Description: "The relative resource name of the lake.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(dataplexLakeNameForAsset),
			},
			{
				Name:        "zone_name",
				Description: "The relative resource name of the zone.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(dataplexZoneNameForAsset),
			},
			{
				Name:        "state",
				Description: "Current state of the asset.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "The time when the asset was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "update_time",
				Description: "The time when the asset was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "Description of the asset.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "uid",
				Description: "System generated globally unique ID for the asset.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "discovery_spec",
				Description: "Specification of the discovery feature applied to data referenced by this asset. When this spec is left unset, the asset will use the spec set on the parent zone.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "resource_spec",
				Description: "Specification of the resource that is referenced by this asset.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "discovery_status",
				Description: "Status of the discovery feature applied to data referenced by this asset.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "resource_status",
				Description: "Status of the resource referenced by this asset.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "security_status",
				Description: "Status of the security policy applied to resource referenced by this asset.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     dataplexAssetSelfLink,
				Transform:   transform.FromValue(),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Labels"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     gcpDataplexAssetTurbotData,
				Transform:   transform.FromField("Akas"),
			},

			// Standard GCP columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Hydrate:     gcpDataplexAssetTurbotData,
				Transform:   transform.FromField("Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Hydrate:     gcpDataplexAssetTurbotData,
				Transform:   transform.FromField("Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listDataplexAssets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	zoneName := d.EqualsQualString("zone_name")
	if zoneName == "" {
		return nil, nil
	}

	// Create Service Connection
	service, err := DataplexService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_dataplex_asset.listDataplexAssets", "connection_error", err)
		return nil, err
	}

	var filters []string
	if d.EqualsQualString("display_name") != "" {
		filters = append(filters, fmt.Sprintf("displayName = \"%s\"", d.EqualsQualString("display_name")))
	}

	if d.EqualsQualString("state") != "" {
		filters = append(filters, fmt.Sprintf("state = \"%s\"", d.EqualsQualString("state")))
	}

	filterString := ""
	if len(filters) > 0 {
		filterString = strings.Join(filters, "AND")
	}

	// Max limit is set as per documentation
	pageSize := types.Int64(1000)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *pageSize {
			pageSize = limit
		}
	}

	resp := service.Projects.Locations.Lakes.Zones.Assets.List(zoneName).PageSize(*pageSize).Filter(filterString)
	if err := resp.Pages(ctx, func(page *dataplex.GoogleCloudDataplexV1ListAssetsResponse) error {
		for _, asset := range page.Assets {
			d.StreamListItem(ctx, asset)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				page.NextPageToken = ""
				return nil
			}
		}
		return nil
	}); err != nil {
		plugin.Logger(ctx).Error("gcp_dataplex_asset.listDataplexAssets", "api_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getDataplexAsset(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	name := d.EqualsQualString("name")

	if len(name) < 1 {
		return nil, nil
	}

	// Create Service Connection
	service, err := DataplexService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_dataplex_asset.getDataplexAsset", "connection_error", err)
		return nil, err
	}

	resp, err := service.Projects.Locations.Lakes.Zones.Assets.Get(name).Do()
	if err != nil {
		plugin.Logger(ctx).Error("gcp_dataplex_asset.getDataplexAsset", "api_error", err)
		return nil, err
	}

	return resp, nil
}

func gcpDataplexAssetTurbotData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	asset := h.Item.(*dataplex.GoogleCloudDataplexV1Asset)
	splitName := strings.Split(asset.Name, "/")

	// Get project details
	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}

	turbotData := map[string]interface{}{
		"Project":  projectId,
		"Location": splitName[3],
		"Akas":     []string{"gcp://dataplex.googleapis.com/" + asset.Name},
	}

	return turbotData, nil
}

func dataplexAssetSelfLink(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(*dataplex.GoogleCloudDataplexV1Asset)

	selfLink := "https://dataplex.googleapis.com/v1/" + data.Name

	return selfLink, nil
}

//// TRANSFORM FUNCTION

func dataplexZoneNameForAsset(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*dataplex.GoogleCloudDataplexV1Asset)
	lakeName := strings.Split(data.Name, "/assets")[0]

	return lakeName, nil
}

func dataplexLakeNameForAsset(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*dataplex.GoogleCloudDataplexV1Asset)
	lakeName := strings.Split(data.Name, "/zones")[0]

	return lakeName, nil
}
