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
	plugin.Logger(ctx).Error("inside tableGcpVertexAIModel")
	return &plugin.Table{
		Name:        "gcp_vertex_ai_model",
		Description: "GCP Vertex AI Model",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			Hydrate:           getAIPlatformModel,
			ShouldIgnoreError: isIgnorableError([]string{"Unimplemented", "InvalidArgument"}),
		},
		List: &plugin.ListConfig{
			Hydrate:           listAIPlatformModels,
			ShouldIgnoreError: isIgnorableError([]string{"Unimplemented", "InvalidArgument"}),
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
				Name:        "artifact_uri",
				Type:        proto.ColumnType_STRING,
				Description: "The path to the directory containing the Model artifact and any of its supporting files.",
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
			// Standard columns for integration with Steampipe
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: "The title of the model.",
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: "A map of tags for the resource.",
			},
		},
	}
}

func listAIPlatformModels(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Error("inside listAIPlatformModels")

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
	getProjectCached := plugin.HydrateFunc(getProject).WithCache()
	projectId, err := getProjectCached(ctx, d, h)
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
	logger.Error("listAIPlatformModels", "Model", it)
	for {
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
		logger.Warn("listAIPlatformModels.models", model.Name)
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
	logger.Error("inside getAIPlatformModels")

	matrixLocation := d.EqualsQualString(matrixKeyLocation)

	// Get project details
	getProjectCached := plugin.HydrateFunc(getProject).WithCache()
	projectId, err := getProjectCached(ctx, d, h)
	if err != nil {
		logger.Error("gcp_vertex_ai_endpoint.getAIPlatformEndpoint", "cache_error", err)
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
		logger.Error("gcp_vertex_ai_endpoint.getAIPlatformModel", "api_error", err)
		return nil, err
	}

	return result, nil
}
