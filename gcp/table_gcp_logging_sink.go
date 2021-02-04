package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"google.golang.org/api/logging/v2"
)

//// TABLE DEFINITION

func tableGcpLoggingSink(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_logging_sink",
		Description: "GCP Logging Sink",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			Hydrate:           getGcpLoggingSink,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listGcpLoggingSinks,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The client-assigned sink identifier, unique within the project",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "destination",
				Description: "Specifies the destination, in which the logs will be exported",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "disabled",
				Description: "Specifies whether the sink is disabled, or not",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "filter",
				Description: "An advanced logs filter. The log entries which will match the filter, will be exported.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A user-specified, human-readable description of the sink",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "The creation timestamp of the sink",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "include_children",
				Description: "Specifies whether a particular log entry from the children is exported depends on the sink's filter expression",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "unique_writer_identity",
				Description: "An IAM identity—a service account or group—under which Logging writes the exported log entries to the sink's destination",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WriterIdentity"),
			},
			{
				Name:        "update_time",
				Description: "The last update timestamp of the sink",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "exclusions",
				Description: "A list of exclusion filters. Log entries that match any of the exclusion filters will not be exported.",
				Type:        proto.ColumnType_JSON,
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
				Transform:   transform.From(sinkNameToAkas),
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
				Transform:   transform.FromConstant(activeProject()),
			},
		},
	}
}

//// FETCH FUNCTIONS

func listGcpLoggingSinks(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	service, err := logging.NewService(ctx)
	if err != nil {
		return nil, err
	}

	project := activeProject()
	resp := service.Projects.Sinks.List("projects/" + project)
	if err := resp.Pages(
		ctx,
		func(page *logging.ListSinksResponse) error {
			for _, sink := range page.Sinks {
				d.StreamListItem(ctx, sink)
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getGcpLoggingSink(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getGcpLoggingSink")

	service, err := logging.NewService(ctx)
	if err != nil {
		return nil, err
	}

	project := activeProject()
	name := d.KeyColumnQuals["name"].GetStringValue()

	op, err := service.Projects.Sinks.Get("projects/" + project + "/sinks/" + name).Do()
	if err != nil {
		logger.Debug("getGcpLoggingSink__", "ERROR", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func sinkNameToAkas(_ context.Context, d *transform.TransformData) (interface{}, error) {
	sink := d.HydrateItem.(*logging.LogSink)
	project := activeProject()

	// Get data for turbot defined properties
	akas := []string{"gcp://logging.googleapis.com/projects/" + project + "/sinks/" + sink.Name}

	return akas, nil
}
