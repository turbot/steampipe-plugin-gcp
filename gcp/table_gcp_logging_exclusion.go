package gcp

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"google.golang.org/api/logging/v2"
)

//// TABLE DEFINITION

func tableGcpLoggingExclusion(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_logging_exclusion",
		Description: "GCP Logging Exclusion",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getGcpLoggingExclusion,
			Tags:       map[string]string{"service": "logging", "action": "exclusions.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listGcpLoggingExclusions,
			Tags:    map[string]string{"service": "logging", "action": "exclusions.list"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The client-assigned identifier, unique within the project",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "disabled",
				Description: "Specifies whether the exclusion is disabled, or not. If disabled it does not exclude any log entries.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "filter",
				Description: "An advanced logs filter that matches the log entries to be excluded",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A user-specified, human-readable description of the exclusion",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "The creation timestamp of the exclusion",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "update_time",
				Description: "The last update timestamp of the exclusion",
				Type:        proto.ColumnType_TIMESTAMP,
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
				Hydrate:     exclusionNameToAkas,
				Transform:   transform.FromValue(),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromConstant("global"),
			},
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

//// FETCH FUNCTIONS

func listGcpLoggingExclusions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := LoggingService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Max limit isn't mentioned in the documentation
	// Default limit is set as 1000
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

	resp := service.Projects.Exclusions.List("projects/" + project).PageSize(*pageSize)
	if err := resp.Pages(
		ctx,
		func(page *logging.ListExclusionsResponse) error {
			// apply rate limiting
			d.WaitForListRateLimit(ctx)

			for _, exclusion := range page.Exclusions {
				d.StreamListItem(ctx, exclusion)

				// Check if context has been cancelled or if the limit has been hit (if specified)
				// if there is a limit, it will return the number of rows required to reach this limit
				if d.RowsRemaining(ctx) == 0 {
					page.NextPageToken = ""
					return nil
				}
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getGcpLoggingExclusion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getGcpLoggingExclusion")

	// Create Service Connection
	service, err := LoggingService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	name := d.EqualsQuals["name"].GetStringValue()

	op, err := service.Projects.Exclusions.Get("projects/" + project + "/exclusions/" + name).Do()
	if err != nil {
		plugin.Logger(ctx).Debug("getGcpLoggingExclusion__", "ERROR", err)
		return nil, err
	}

	// If the name has been passed as empty string, API does not returns any error
	if len(op.Name) < 1 {
		return nil, nil
	}

	return op, nil
}

func exclusionNameToAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	exclusion := h.Item.(*logging.LogExclusion)

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	akas := []string{"gcp://logging.googleapis.com/projects/" + project + "/exclusions/" + exclusion.Name}
	return akas, nil
}
