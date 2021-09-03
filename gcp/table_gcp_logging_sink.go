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
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getGcpLoggingSink,
		},
		List: &plugin.ListConfig{
			Hydrate:           listGcpLoggingSinks,
			ShouldIgnoreError: isIgnorableError([]string{"403"}),
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
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CreateTime").NullIfZero(),
			},
			{
				Name:        "include_children",
				Description: "Specifies whether a particular log entry from the children is exported depends on the sink's filter expression",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getSinkSelfLink,
				Transform:   transform.FromValue(),
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
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("UpdateTime").NullIfZero(),
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
				Hydrate:     sinkNameToAkas,
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
				Hydrate:     plugin.HydrateFunc(getProject).WithCache(),
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// FETCH FUNCTIONS

func listGcpLoggingSinks(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := LoggingService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details
	getProjectCached := plugin.HydrateFunc(getProject).WithCache()
	projectId, err := getProjectCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

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
	plugin.Logger(ctx).Trace("getGcpLoggingSink")

	// Create Service Connection
	service, err := LoggingService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details
	getProjectCached := plugin.HydrateFunc(getProject).WithCache()
	projectId, err := getProjectCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	name := d.KeyColumnQuals["name"].GetStringValue()

	op, err := service.Projects.Sinks.Get("projects/" + project + "/sinks/" + name).Do()
	if err != nil {
		plugin.Logger(ctx).Debug("getGcpLoggingSink__", "ERROR", err)
		return nil, err
	}

	// If the name has been passed as empty string, API does not returns any error
	if len(op.Name) < 1 {
		return nil, nil
	}

	return op, nil
}

func sinkNameToAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	sink := h.Item.(*logging.LogSink)

	// Get project details
	getProjectCached := plugin.HydrateFunc(getProject).WithCache()
	projectId, err := getProjectCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	akas := []string{"gcp://logging.googleapis.com/projects/" + project + "/sinks/" + sink.Name}
	return akas, nil
}

func getSinkSelfLink(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	sink := h.Item.(*logging.LogSink)

	// Get project details
	getProjectCached := plugin.HydrateFunc(getProject).WithCache()
	projectId, err := getProjectCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	selfLink := "https://www.googleapis.com/logging/v2/projects/" + project + "/sinks/" + sink.Name
	return selfLink, nil
}
