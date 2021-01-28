package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"google.golang.org/api/logging/v2"
)

//// TABLE DEFINITION

func tableGcpLoggingExclusion(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_logging_exclusion",
		Description: "GCP Logging Exclusion",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			Hydrate:           getGcpLoggingExclusion,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listGcpLoggingExclusions,
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
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "update_time",
				Description: "The last update timestamp of the exclusion",
				Type:        proto.ColumnType_STRING,
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
				Transform:   transform.From(exclusionNameToAkas),
			},

			// standard gcp columns
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromConstant(activeProject()),
			},
		},
	}
}

//// FETCH FUNCTIONS

func listGcpLoggingExclusions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	service, err := logging.NewService(ctx)
	if err != nil {
		return nil, err
	}

	project := activeProject()
	resp := service.Projects.Exclusions.List("projects/" + project)
	if err := resp.Pages(
		ctx,
		func(page *logging.ListExclusionsResponse) error {
			for _, exclusion := range page.Exclusions {
				d.StreamListItem(ctx, exclusion)
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
	logger := plugin.Logger(ctx)
	logger.Trace("getGcpLoggingExclusion")

	service, err := logging.NewService(ctx)
	if err != nil {
		return nil, err
	}

	project := activeProject()
	name := d.KeyColumnQuals["name"].GetStringValue()

	op, err := service.Projects.Exclusions.Get("projects/" + project + "/exclusions/" + name).Do()
	if err != nil {
		logger.Debug("getGcpLoggingExclusion__", "ERROR", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func exclusionNameToAkas(_ context.Context, d *transform.TransformData) (interface{}, error) {
	exclusion := d.HydrateItem.(*logging.LogExclusion)
	project := activeProject()

	// Get data for turbot defined properties
	akas := []string{"gcp://logging.googleapis.com/projects/" + project + "/exclusions/" + exclusion.Name}

	return akas, nil
}
