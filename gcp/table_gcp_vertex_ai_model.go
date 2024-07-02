package gcp

import (
	"context"
	"strings"

	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/iterator"
)

func tableGcpVertexAIModel(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_vertex_ai_model",
		Description: "GCP Vertex AI Model",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			Hydrate:           getAIPlatformModel,
			ShouldIgnoreError: isIgnorableError([]string{"Unauthenticated", "Unimplemented", "InvalidArgument"}),
			Tags:              map[string]string{"service": "aiplatform", "action": "models.get"},
		},
		List: &plugin.ListConfig{
			Hydrate:           listAIPlatformModels,
			ShouldIgnoreError: isIgnorableError([]string{"Unauthenticated", "Unimplemented", "InvalidArgument"}),
			Tags:              map[string]string{"service": "aiplatform", "action": "models.list"},
		},
		GetMatrixItemFunc: BuildVertexAILocationListByClientType("Model"),
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The resource name of the Model.",
				Transform:   transform.FromField("Name").Transform(lastPathElement),
			},
			{
				Name:        "version_id",
				Type:        proto.ColumnType_STRING,
				Description: "The version ID of the model.",
			},
			{
				Name:        "version_create_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("VersionCreateTime").Transform(convertTimestamppbAsTime),
				Description: "Timestamp when this version was created.",
			},
			{
				Name:        "version_update_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("VersionUpdateTime").Transform(convertTimestamppbAsTime),
				Description: "Timestamp when this version was most recently updated.",
			},
			{
				Name:        "display_name",
				Type:        proto.ColumnType_STRING,
				Description: "The display name of the Model.",
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "The description of the Model.",
			},
			{
				Name:        "version_description",
				Type:        proto.ColumnType_STRING,
				Description: "The description of this version.",
			},
			{
				Name:        "metadata_schema_uri",
				Type:        proto.ColumnType_STRING,
				Description: "Points to a YAML file stored on Google Cloud Storage describing additional information about the model.",
			},
			{
				Name:        "training_pipeline",
				Type:        proto.ColumnType_STRING,
				Description: "The resource name of the TrainingPipeline that uploaded this model, if any.",
			},
			{
				Name:        "pipeline_job",
				Type:        proto.ColumnType_STRING,
				Description: "Populated if the model is produced by a pipeline job.",
			},
			{
				Name:        "artifact_uri",
				Type:        proto.ColumnType_STRING,
				Description: "The path to the directory containing the Model artifact and its supporting files.",
			},
			{
				Name:        "create_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CreateTime").Transform(convertTimestamppbAsTime),
				Description: "Timestamp when this Model was uploaded into Vertex AI.",
			},
			{
				Name:        "update_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("UpdateTime").Transform(convertTimestamppbAsTime),
				Description: "Timestamp when this Model was most recently updated.",
			},
			{
				Name:        "etag",
				Type:        proto.ColumnType_STRING,
				Description: "Used to perform consistent read-modify-write updates.",
			},
			{
				Name:        "metadata_artifact",
				Type:        proto.ColumnType_STRING,
				Description: "The resource name of the Artifact that was created in MetadataStore when creating the model.",
			},
			// JSON columns
			{
				Name:        "version_aliases",
				Type:        proto.ColumnType_JSON,
				Description: "User provided version aliases so that a model version can be referenced via alias.",
			},
			{
				Name:        "predict_schemata",
				Type:        proto.ColumnType_JSON,
				Description: "The schemata that describe formats of the model's predictions and explanations.",
			},
			{
				Name:        "metadata",
				Type:        proto.ColumnType_JSON,
				Description: "An additional information about the model; the schema of the metadata can be found in metadata_schema_uri, immutable.",
			},
			{
				Name:        "supported_export_formats",
				Type:        proto.ColumnType_JSON,
				Description: "The formats in which this model may be exported. If empty, this model is not available for export.",
			},
			{
				Name:        "container_spec",
				Type:        proto.ColumnType_JSON,
				Description: "The specification of the container that is to be used when deploying this model.",
			},
			{
				Name:        "supported_deployment_resources_types",
				Type:        proto.ColumnType_JSON,
				Description: "The configuration types this model supports for deployment.",
			},
			{
				Name:        "supported_input_storage_formats",
				Type:        proto.ColumnType_JSON,
				Description: "The formats this Model supports in BatchPredictionJob.input_config.",
			},
			{
				Name:        "supported_output_storage_formats",
				Type:        proto.ColumnType_JSON,
				Description: "The formats this Model supports in BatchPredictionJob.output_config.",
			},
			{
				Name:        "deployed_models",
				Type:        proto.ColumnType_JSON,
				Description: "The pointers to DeployedModels created from this Model.",
			},
			{
				Name:        "explanation_spec",
				Type:        proto.ColumnType_JSON,
				Description: "The default explanation specification for this Model.",
			},
			{
				Name:        "labels",
				Type:        proto.ColumnType_JSON,
				Description: "The labels with user-defined metadata to organize your Models.",
			},
			{
				Name:        "encryption_spec",
				Type:        proto.ColumnType_JSON,
				Description: "Customer-managed encryption key spec for a Model.",
			},
			{
				Name:        "model_source_info",
				Type:        proto.ColumnType_JSON,
				Description: "Source of a model.",
			},
			{
				Name:        "original_model_info",
				Type:        proto.ColumnType_JSON,
				Description: "If this model is a copy of another model, this contains info about the original.",
			},
			// Standard columns for integration with Steampipe
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
				Transform:   transform.FromP(gcpModelStandard, "Akas"),
				Description: ColumnDescriptionAkas,
			},
			// Standard gcp columns
			{
				Name:        "location",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpModelStandard, "Location"),
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

func listAIPlatformModels(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
		logger.Warn("gcp_vertex_ai_model.listAIPlatformModels", "location", region, "matrixLocation", location)
		return nil, nil
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		logger.Error("gcp_vertex_ai_model.listAIPlatformModels", "cache_error", err)
		return nil, err
	}

	project := projectId.(string)

	// Page size should be in range of [0, 100].
	pageSize := types.Int64(100)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *pageSize {
			pageSize = limit
		}
	}

	// Create Service Connection
	service, err := AIService(ctx, d, "Model")
	if err != nil {
		logger.Error("gcp_vertex_ai_model.listAIPlatformModels", "connection_error", err)
		return nil, err
	}

	req := &aiplatformpb.ListModelsRequest{
		Parent:   "projects/" + project + "/locations/" + location,
		PageSize: int32(*pageSize),
	}

	it := service.Model.ListModels(ctx, req)

	for {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		model, err := it.Next()
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				return nil, nil
			}
			if err == iterator.Done {
				break
			}
			logger.Error("gcp_vertex_ai_model.listAIPlatformModels", err)
			return nil, err
		}

		d.StreamListItem(ctx, model)

		if d.RowsRemaining(ctx) == 0 {
			break
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAIPlatformModel(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	matrixLocation := d.EqualsQualString(matrixKeyLocation)

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		logger.Error("gcp_vertex_ai_model.getAIPlatformModel", "cache_error", err)
		return nil, err
	}
	project := projectId.(string)

	name := d.EqualsQualString("name")

	// Validate - name should not be blank
	if name == "" {
		return nil, nil
	}

	service, err := AIService(ctx, d, "Model")
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "NotFound") {
			return nil, nil
		}
		logger.Error("gcp_vertex_ai_model.getAIPlatformModel", "service_error", err)
		return nil, err
	}

	// Assuming the 'name' column contains the full resource name of the Model
	// e.g., projects/projectID/locations/locationID/models/modelID
	req := &aiplatformpb.GetModelRequest{
		Name: "projects/" + project + "/locations/" + matrixLocation + "/models/" + name,
	}

	// Call the API
	result, err := service.Model.GetModel(ctx, req)
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "NotFound") {
			return nil, nil
		}
		logger.Error("gcp_vertex_ai_model.getAIPlatformModel", "api_error", err)
		return nil, err
	}

	return result, nil
}

/// TRANSFORM FUNCTIONS

func gcpModelStandard(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	param := d.Param.(string)
	AIData := d.HydrateItem.(*aiplatformpb.Model)
	akas := []string{"gcp://aiplatform.googleapis.com/" + AIData.Name}

	turbotData := map[string]interface{}{
		"Location": strings.Split(AIData.Name, "/")[3],
		"Akas":     akas,
	}
	return turbotData[param], nil
}
