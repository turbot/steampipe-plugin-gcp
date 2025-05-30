package gcp

import (
	"context"
	"strings"

	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"google.golang.org/api/iterator"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGcpVertexAINotebookRuntimeTemplate(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_vertex_ai_notebook_runtime_template",
		Description: "GCP Vertex AI Notebook Runtime Template",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			Hydrate:           getAIPlatformNotebookRuntimeTemplate,
			ShouldIgnoreError: isIgnorableError([]string{"Unauthenticated", "Unimplemented", "InvalidArgument"}),
			Tags:              map[string]string{"service": "aiplatform", "action": "notebookRuntimeTemplates.get"},
		},
		List: &plugin.ListConfig{
			Hydrate:           listAIPlatformNotebookRuntimeTemplates,
			ShouldIgnoreError: isIgnorableError([]string{"Unauthenticated", "Unimplemented", "InvalidArgument"}),
			Tags:              map[string]string{"service": "aiplatform", "action": "notebookRuntimeTemplates.list"},
		},
		GetMatrixItemFunc: BuildVertexAILocationListByClientType("Notebook"),
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The resource name of the Notebook Runtime Template.",
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "display_name",
				Type:        proto.ColumnType_STRING,
				Description: "The display name of the Notebook Runtime Template.",
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "The description of the Notebook Runtime Template.",
			},
			{
				Name:        "notebook_runtime_type",
				Type:        proto.ColumnType_STRING,
				Description: "The type of the notebook runtime template.",
				Transform:   transform.From(getNotebookRuntimeTemplateType),
			},
			{
				Name:        "is_default",
				Type:        proto.ColumnType_BOOL,
				Description: "Specifies whether this is the default template.",
			},
			{
				Name:        "service_account",
				Type:        proto.ColumnType_STRING,
				Description: "The service account that the runtime workload runs as.",
			},
			{
				Name:        "etag",
				Type:        proto.ColumnType_STRING,
				Description: "Used to perform consistent read-modify-write updates.",
			},
			{
				Name:        "create_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CreateTime").Transform(convertTimestamppbAsTime),
				Description: "Timestamp when this Notebook Runtime Template was created.",
			},
			{
				Name:        "update_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("UpdateTime").Transform(convertTimestamppbAsTime),
				Description: "Timestamp when this Notebook Runtime Template was last updated.",
			},
			{
				Name:        "machine_spec",
				Type:        proto.ColumnType_JSON,
				Description: "The specification of a single machine for the template.",
			},
			{
				Name:        "data_persistent_disk_spec",
				Type:        proto.ColumnType_JSON,
				Description: "The specification of persistent disk attached to the runtime as data disk storage.",
			},
			{
				Name:        "network_spec",
				Type:        proto.ColumnType_JSON,
				Description: "The specification of the network for the Notebook Runtime Template.",
			},
			{
				Name:        "idle_shutdown_config",
				Type:        proto.ColumnType_JSON,
				Description: "The idle shutdown configuration of the Notebook Runtime Template.",
			},
			{
				Name:        "euc_config",
				Type:        proto.ColumnType_JSON,
				Description: "The EUC (End User Computing) configuration of the Notebook Runtime Template.",
			},
			{
				Name:        "shielded_vm_config",
				Type:        proto.ColumnType_JSON,
				Description: "The configuration for Shielded VM for the runtime template.",
			},
			{
				Name:        "network_tags",
				Type:        proto.ColumnType_JSON,
				Description: "The Compute Engine network tags to add to the runtime.",
			},

			// Standard steampipe columns
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
				Description: ColumnDescriptionTitle,
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Labels"),
				Description: ColumnDescriptionTags,
			},
			{
				Name:        "akas",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(gcpNotebookRuntimeTemplate, "Akas"),
				Description: ColumnDescriptionAkas,
			},

			// Standard gcp columns
			{
				Name:        "location",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpNotebookRuntimeTemplate, "Location"),
				Description: ColumnDescriptionLocation,
			},
			{
				Name:        "project",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getProject,
				Transform:   transform.FromValue(),
				Description: ColumnDescriptionProject,
			},
		},
	}
}

func listAIPlatformNotebookRuntimeTemplates(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	region := d.EqualsQualString("location")
	var location string
	matrixLocation := d.EqualsQualString(matrixKeyLocation)
	// Since, when the service API is disabled, matrixLocation value will be nil
	if matrixLocation != "" {
		location = matrixLocation
	}

	// Minimize API call as per given location
	if region != "" && region != location {
		logger.Warn("gcp_vertex_ai_notebook_runtime_template.listAIPlatformNotebookRuntimeTemplates", "location", region, "matrixLocation", location)
		return nil, nil
	}

	// Get project details
	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	// Create Service Connection
	service, err := AIService(ctx, d, "Notebook")
	if err != nil {
		logger.Error("gcp_vertex_ai_notebook_runtime_template.listAIPlatformNotebookRuntimeTemplates", "service_error", err)
		return nil, err
	}

	// Max limit isn't mentioned in the documentation
	// Default limit is set as 1000
	pageSize := int32(1000)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(pageSize) {
			pageSize = int32(*limit)
		}
	}

	req := &aiplatformpb.ListNotebookRuntimeTemplatesRequest{
		Parent:   "projects/" + project + "/locations/" + location,
		PageSize: pageSize,
	}

	// Call the API
	it := service.Notebook.ListNotebookRuntimeTemplates(ctx, req)
	for {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Error("gcp_vertex_ai_notebook_runtime_template.listAIPlatformNotebookRuntimeTemplates", "api_error", err)
			return nil, err
		}

		d.StreamListItem(ctx, resp)

		// Check if context has been cancelled or if the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			break
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAIPlatformNotebookRuntimeTemplate(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	matrixLocation := d.EqualsQualString(matrixKeyLocation)

	name := d.EqualsQualString("name")
	splitName := strings.Split(name, "/")

	// Validate - name should not be blank and restrict the API call for other locations
	if len(name) > 3 && splitName[3] != matrixLocation {
		return nil, nil
	}

	service, err := AIService(ctx, d, "Notebook")
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "NotFound") {
			return nil, nil
		}
		logger.Error("gcp_vertex_ai_notebook_runtime_template.getAIPlatformNotebookRuntimeTemplate", "service_error", err)
		return nil, err
	}

	// Assuming the 'name' column contains the full resource name of the Model
	// e.g., projects/projectID/locations/locationID/models/modelID
	req := &aiplatformpb.GetNotebookRuntimeTemplateRequest{
		Name: name,
	}

	// Call the API
	result, err := service.Notebook.GetNotebookRuntimeTemplate(ctx, req)
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "NotFound") {
			return nil, nil
		}
		logger.Error("gcp_vertex_ai_notebook_runtime_template.getAIPlatformNotebookRuntimeTemplate", "api_error", err)
		return nil, err
	}

	return result, nil
}

/// TRANSFORM FUNCTIONS

func gcpNotebookRuntimeTemplate(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	param := d.Param.(string)
	AIData := d.HydrateItem.(*aiplatformpb.NotebookRuntimeTemplate)
	akas := []string{"gcp://aiplatform.googleapis.com/" + AIData.Name}

	turbotData := map[string]interface{}{
		"Location": strings.Split(AIData.Name, "/")[3],
		"Akas":     akas,
	}
	return turbotData[param], nil
}

func getNotebookRuntimeTemplateType(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	AIData := d.HydrateItem.(*aiplatformpb.NotebookRuntimeTemplate)

	return aiplatformpb.NotebookRuntimeType_name[int32(AIData.NotebookRuntimeType)], nil
}
