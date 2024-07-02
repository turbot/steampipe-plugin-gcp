package gcp

import (
	"context"
	"strings"

	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableGcpVertexAIEndpoint(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_vertex_ai_endpoint",
		Description: "GCP Vertex AI Endpoint",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			Hydrate:           getAIPlatformEndpoint,
			ShouldIgnoreError: isIgnorableError([]string{"Unimplemented"}),
			Tags:              map[string]string{"service": "aiplatform", "action": "endpoints.get"},
		},
		List: &plugin.ListConfig{
			Hydrate:           listAIPlatformEndpoints,
			ShouldIgnoreError: isIgnorableError([]string{"Unimplemented"}),
			Tags:              map[string]string{"service": "aiplatform", "action": "endpoints.list"},
		},
		GetMatrixItemFunc: BuildVertexAILocationListByClientType("Endpoint"),
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(lastPathElement),
				Description: "The resource name of the Endpoint.",
			},
			{
				Name:        "create_time",
				Description: "Timestamp when this Endpoint was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CreateTime").Transform(convertTimestamppbAsTime),
			},
			{
				Name:        "display_name",
				Type:        proto.ColumnType_STRING,
				Description: "Human-readable display name of this key that you can modify.",
			},
			{
				Name:        "description",
				Description: "The description of the Endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "enable_private_service_connect",
				Description: "If true, expose the Endpoint via private service connect. Only one of the fields, network or enable_private_service_connect, can be set.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "etag",
				Description: "Used to perform consistent read-modify-write updates. If not set, a blind 'overwrite' update happens.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "model_deployment_monitoring_job",
				Description: "Resource name of the Model Monitoring job associated with this Endpoint if monitoring is enabled by JobService.CreateModelDeploymentMonitoringJob.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "network",
				Description: "The full name of the Google Compute Engine network (https://cloud.google.com//compute/docs/networks-and-firewalls#networks) to which the Endpoint should be peered.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "update_time",
				Description: "Timestamp when this Endpoint was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("UpdateTime").Transform(convertTimestamppbAsTime),
			},

			// JSON columns
			{
				Name:        "deployed_models",
				Description: "The models deployed in this Endpoint. To add or remove DeployedModels use EndpointService.DeployModel and EndpointService.UndeployModel respectively.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "encryption_spec",
				Description: "Customer-managed encryption key spec for an Endpoint. If set, this Endpoint and all sub-resources of this Endpoint will be secured by this key.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "predict_request_response_logging_config",
				Description: "Configures the request-response logging for online prediction.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "traffic_split",
				Description: "A map from a DeployedModel's ID to the percentage of this Endpoint's traffic that should be forwarded to that DeployedModel.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "labels",
				Description: "The labels with user-defined metadata to organize your Endpoints.",
				Type:        proto.ColumnType_JSON,
			},

			// standard steampipe columns
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
				Transform:   transform.FromP(gcpAIPlatformTurbotData, "Akas"),
			},
			// Standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpAIPlatformTurbotData, "Location"),
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

//// LIST FUNCTION

func listAIPlatformEndpoints(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
		return nil, nil
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		logger.Error("gcp_vertex_ai_endpoint.listAIPlatformEndpoints", "cache_error", err)
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
	service, err := AIService(ctx, d, "Endpoint")
	if err != nil {
		logger.Error("gcp_vertex_ai_endpoint.listAIPlatformEndpoints", "service_error", err)
		return nil, err
	}

	input := &aiplatformpb.ListEndpointsRequest{
		Parent:   "projects/" + project + "/locations/" + location,
		PageSize: int32(*pageSize),
	}

	it := service.Endpoint.ListEndpoints(ctx, input)

	for {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		resp, err := it.Next()
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				return nil, nil
			}
			if err == iterator.Done {
				break
			}
			logger.Error("gcp_vertex_ai_endpoint.listAIPlatformEndpoints", "api_error", err)
			return nil, err
		}

		d.StreamListItem(ctx, resp)

		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAIPlatformEndpoint(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	matrixLocation := d.EqualsQualString(matrixKeyLocation)

	// Get project details

	projectId, err := getProject(ctx, d, h)
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

	// Create Service Connection
	service, err := AIService(ctx, d, "Endpoint")
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "NotFound") {
			return nil, nil
		}
		logger.Error("gcp_vertex_ai_endpoint.getAIPlatformEndpoint", "service_error", err)
		return nil, err
	}
	input := &aiplatformpb.GetEndpointRequest{
		Name: "projects/" + project + "/locations/" + matrixLocation + "/endpoints/" + name,
	}
	op, err := service.Endpoint.GetEndpoint(ctx, input)
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "NotFound") {
			return nil, nil
		}
		logger.Error("gcp_vertex_ai_endpoint.getAIPlatformEndpoint", "api_error", err)
		return nil, err
	}
	return op, nil
}

/// TRANSFORM FUNCTIONS

func gcpAIPlatformTurbotData(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	param := d.Param.(string)
	AIData := d.HydrateItem.(*aiplatformpb.Endpoint)
	akas := []string{"gcp://aiplatform.googleapis.com/" + AIData.Name}

	turbotData := map[string]interface{}{
		"Location": strings.Split(AIData.Name, "/")[3],
		"Akas":     akas,
	}
	return turbotData[param], nil
}

func convertTimestamppbAsTime(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	v := d.Value
	if v != nil {
		timeValue := v.(*timestamppb.Timestamp)
		return timeValue.AsTime(), nil
	}
	return nil, nil
}
