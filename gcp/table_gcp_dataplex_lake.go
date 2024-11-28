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

func tableGcpDataplexLake(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_dataplex_lake",
		Description: "GCP Dataplex Lake",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getDataplexLake,
		},
		List: &plugin.ListConfig{
			Hydrate: listDataplexLakes,
			KeyColumns: plugin.KeyColumnSlice{
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
				Description: "The relative resource name of the lake.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "Current state of the lake.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "The time when the lake was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "update_time",
				Description: "The time when the lake was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "Description of the lake.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "uid",
				Description: "System generated globally unique ID for the lake.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_account",
				Description: "Service account associated with this lake. This service account must be authorized to access or operate on resources managed by the lake.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "asset_status",
				Description: "Aggregated status of the underlying assets of the lake.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "metastore",
				Description: "Settings to manage lake and Dataproc Metastore service instance association.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "metastore_status",
				Description: "Metastore status of the lake.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     dataplexLakeSelfLink,
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
				Hydrate:     gcpDataplexLakeTurbotData,
				Transform:   transform.FromField("Akas"),
			},

			// Standard GCP columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Hydrate:     gcpDataplexLakeTurbotData,
				Transform:   transform.FromField("Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Hydrate:     gcpDataplexLakeTurbotData,
				Transform:   transform.FromField("Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listDataplexLakes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var location string
	matrixLocation := d.EqualsQualString(matrixKeyLocation)
	// Since, when the service API is disabled, matrixLocation value will be nil
	if matrixLocation != "" {
		location = matrixLocation
	}

	// We should not make the APi call for the regions "global", "eu" and "us"".
	// Error 400: Malformed name: 'projects/parker-aaa/locations/global/lakes'
	//  Error 400: Malformed name: 'projects/parker-aaa/locations/us/lakes'
	if location == "global" || location == "us" || location == "eu" {
		return nil, nil
	}

	// Create Service Connection
	service, err := DataplexService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_dataplex_lake.listDataplexLakes", "connection_error", err)
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

	// Get project details
	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	parent := "projects/" + project + "/locations/" + location

	resp := service.Projects.Locations.Lakes.List(parent).PageSize(*pageSize).Filter(filterString)
	if err := resp.Pages(ctx, func(page *dataplex.GoogleCloudDataplexV1ListLakesResponse) error {
		for _, lake := range page.Lakes {
			d.StreamListItem(ctx, lake)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				page.NextPageToken = ""
				return nil
			}
		}
		return nil
	}); err != nil {
		plugin.Logger(ctx).Error("gcp_dataplex_lake.listDataplexLakes", "api_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getDataplexLake(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
	// Error 400: Malformed name: 'projects/parker-aaa/locations/global/lakes'
	// Error 400: Malformed name: 'projects/parker-aaa/locations/us/lakes'
	if location == "global" || location == "us" || location == "eu" {
		return nil, nil
	}

	// Create Service Connection
	service, err := DataplexService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_dataplex_lake.getDataplexLake", "connection_error", err)
		return nil, err
	}

	resp, err := service.Projects.Locations.Lakes.Get(name).Do()
	if err != nil {
		plugin.Logger(ctx).Error("gcp_dataplex_lake.getDataplexLake", "api_error", err)
		return nil, err
	}

	return resp, nil
}

func gcpDataplexLakeTurbotData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	lake := h.Item.(*dataplex.GoogleCloudDataplexV1Lake)

	// Get project details
	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}

	var location string
	matrixLocation := d.EqualsQualString(matrixKeyLocation)
	// Since, when the service API is disabled, matrixLocation value will be nil
	if matrixLocation != "" {
		location = matrixLocation
	}

	turbotData := map[string]interface{}{
		"Project":  projectId,
		"Location": location,
		"Akas":     []string{"gcp://dataplex.googleapis.com/" + lake.Name},
	}

	return turbotData, nil
}

func dataplexLakeSelfLink(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(*dataplex.GoogleCloudDataplexV1Lake)

	selfLink := "https://dataplex.googleapis.com/v1/" + data.Name

	return selfLink, nil
}
