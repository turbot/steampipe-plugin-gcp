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

func tableGcpDataplexTask(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_dataplex_task",
		Description: "GCP Dataplex Task",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getDataplexTask,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listDataplexLakes,
			Hydrate:       listDataplexTasks,
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
				Description: "The relative resource name of the task.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lake_name",
				Description: "The relative resource name of the lake.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(dataplexLakeNameForTask),
			},
			{
				Name:        "state",
				Description: "Current state of the task.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "The time when the task was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "update_time",
				Description: "The time when the task was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "Description of the task.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "uid",
				Description: "System generated globally unique ID for the task.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "trigger_spec",
				Description: "Spec related to how often and when a task should be triggered.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "execution_spec",
				Description: "Spec related to how a task is executed.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "execution_status",
				Description: "Status of the latest task executions.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "notebook",
				Description: "Config related to running scheduled Notebooks.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "spark",
				Description: "Config related to running custom Spark tasks.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     dataplexTaskSelfLink,
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
				Hydrate:     gcpDataplexTaskTurbotData,
				Transform:   transform.FromField("Akas"),
			},

			// Standard GCP columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Hydrate:     gcpDataplexTaskTurbotData,
				Transform:   transform.FromField("Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Hydrate:     gcpDataplexTaskTurbotData,
				Transform:   transform.FromField("Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listDataplexTasks(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	lake := h.Item.(*dataplex.GoogleCloudDataplexV1Lake)

	lakeName := d.EqualsQualString("lake_name")

	if lakeName != "" && lakeName != lake.Name{
		return nil, nil
	}

	// Create Service Connection
	service, err := DataplexService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_dataplex_task.listDataplexTasks", "connection_error", err)
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

	resp := service.Projects.Locations.Lakes.Tasks.List(lake.Name).PageSize(*pageSize).Filter(filterString)
	if err := resp.Pages(ctx, func(page *dataplex.GoogleCloudDataplexV1ListTasksResponse) error {
		for _, task := range page.Tasks {
			d.StreamListItem(ctx, task)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				page.NextPageToken = ""
				return nil
			}
		}
		return nil
	}); err != nil {
		plugin.Logger(ctx).Error("gcp_dataplex_task.listDataplexTasks", "api_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getDataplexTask(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
		plugin.Logger(ctx).Error("gcp_dataplex_task.getDataplexLake", "connection_error", err)
		return nil, err
	}

	resp, err := service.Projects.Locations.Lakes.Tasks.Get(name).Do()
	if err != nil {
		plugin.Logger(ctx).Error("gcp_dataplex_task.getDataplexLake", "api_error", err)
		return nil, err
	}

	return resp, nil
}

func gcpDataplexTaskTurbotData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	task := h.Item.(*dataplex.GoogleCloudDataplexV1Task)
	splitName := strings.Split(task.Name, "/")

	// Get project details
	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}

	turbotData := map[string]interface{}{
		"Project":  projectId,
		"Location": splitName[3],
		"Akas":     []string{"gcp://dataplex.googleapis.com/" + task.Name},
	}

	return turbotData, nil
}

func dataplexTaskSelfLink(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(*dataplex.GoogleCloudDataplexV1Task)

	selfLink := "https://dataplex.googleapis.com/v1/" + data.Name

	return selfLink, nil
}

//// TRANSFORM FUNCTION

func dataplexLakeNameForTask(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*dataplex.GoogleCloudDataplexV1Task)
	lakeName := strings.Split(data.Name, "/tasks")[0]

	return lakeName, nil
}
