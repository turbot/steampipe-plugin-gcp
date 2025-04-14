package gcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"google.golang.org/api/dataplex/v1"
)

func tableGcpDataplexZone(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_dataplex_zone",
		Description: "GCP Dataplex Zone",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getDataplexZone,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listDataplexLakes,
			Hydrate:       listDataplexZones,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "lake_name", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "display_name", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "state", Require: plugin.Optional, Operators: []string{"="}},
			},
		},
		GetMatrixItemFunc: BuildDataplexLocationList,
		Columns: []*plugin.Column{
			{
				Name:        "display_name",
				Description: "User friendly display name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "The relative resource name of the zone.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lake_name",
				Description: "The relative resource name of the lake.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(dataplexLakeName),
			},
			{
				Name:        "state",
				Description: "Current state of the zone.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "The time when the zone was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "update_time",
				Description: "The time when the zone was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "Description of the zone.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "uid",
				Description: "System generated globally unique ID for the zone.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of the zone.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "asset_status",
				Description: "Aggregated status of the underlying assets of the zone.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "resource_spec",
				Description: "Specification of the resources that are referenced by the assets within this zone.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "discovery_spec",
				Description: "Specification of the discovery feature applied to data in this zone.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     dataplexZoneSelfLink,
				Transform:   transform.FromValue(),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(lastPathElement),
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
				Hydrate:     gcpDataplexZoneTurbotData,
				Transform:   transform.FromField("Akas"),
			},

			// Standard GCP columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Hydrate:     gcpDataplexZoneTurbotData,
				Transform:   transform.FromField("Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Hydrate:     gcpDataplexZoneTurbotData,
				Transform:   transform.FromField("Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listDataplexZones(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	lake := h.Item.(*dataplex.GoogleCloudDataplexV1Lake)

	lakeName := d.EqualsQualString("lake_name")

	if lakeName != "" && lakeName != lake.Name{
		return nil, nil
	}

	// Create Service Connection
	service, err := DataplexService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_dataplex_zone.listDataplexZones", "connection_error", err)
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

	resp := service.Projects.Locations.Lakes.Zones.List(lake.Name).PageSize(*pageSize).Filter(filterString)
	if err := resp.Pages(ctx, func(page *dataplex.GoogleCloudDataplexV1ListZonesResponse) error {
		for _, zone := range page.Zones {
			d.StreamListItem(ctx, zone)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				page.NextPageToken = ""
				return nil
			}
		}
		return nil
	}); err != nil {
		plugin.Logger(ctx).Error("gcp_dataplex_zone.listDataplexZones", "api_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getDataplexZone(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var location string
	matrixLocation := d.EqualsQualString(matrixKeyLocation)
	// Since, when the service API is disabled, matrixLocation value will be nil
	if matrixLocation != "" {
		location = matrixLocation
	}

	name := d.EqualsQualString("name")

	if len(name) < 1 {
		return nil, nil
	}

	// Check for handle Array index out of range error for any wrong input in query parameter.
	// We should not make the API call for other regions.
	splitName := strings.Split(name, "/")
	if len(splitName) > 3 && splitName[3] != location {
		return nil, nil
	}

	// We should not make the API call for the regions "global", "eu" and "us"".
	if location == "global" || location == "us" || location == "eu" {
		return nil, nil
	}

	// Create Service Connection
	service, err := DataplexService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_dataplex_zone.getDataplexLake", "connection_error", err)
		return nil, err
	}

	resp, err := service.Projects.Locations.Lakes.Zones.Get(name).Do()
	if err != nil {
		plugin.Logger(ctx).Error("gcp_dataplex_zone.getDataplexLake", "api_error", err)
		return nil, err
	}

	return resp, nil
}

func gcpDataplexZoneTurbotData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	zone := h.Item.(*dataplex.GoogleCloudDataplexV1Zone)
	splitName := strings.Split(zone.Name, "/")

	// Get project details
	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}

	turbotData := map[string]interface{}{
		"Project":  projectId,
		"Location": splitName[3],
		"Akas":     []string{"gcp://dataplex.googleapis.com/" + zone.Name},
	}

	return turbotData, nil
}

func dataplexZoneSelfLink(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(*dataplex.GoogleCloudDataplexV1Zone)

	selfLink := "https://dataplex.googleapis.com/v1/" + data.Name

	return selfLink, nil
}

//// TRANSFORM FUNCTION

func dataplexLakeName(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*dataplex.GoogleCloudDataplexV1Zone)
	lakeName := strings.Split(data.Name, "/zones")[0]

	return lakeName, nil
}
