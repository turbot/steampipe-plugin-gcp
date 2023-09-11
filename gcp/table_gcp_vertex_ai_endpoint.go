package gcp

import (
	"context"

	aiplatformpb "cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/iterator"
	aiplatform "google.golang.org/genproto/googleapis/cloud/aiplatform/v1beta1"
)

//// TABLE DEFINITION

func tableGcpVertexAIEndpoint(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_vertex_ai_endpoint",
		Description: "GCP Vertex AI Endpoint",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name"}),
			Hydrate:    getVertexAIEndpoint,
		},
		List: &plugin.ListConfig{
			Hydrate: listVertexAIEndpoints,
			// KeyColumns: plugin.KeyColumnSlice{
			// 	{Name: "location", Require: plugin.Optional},
			// },
		},
		GetMatrixItemFunc: BuildAIPlatformLocationList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The resource name of the Endpoint.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(lastPathElement),
			},
			{
				Name:        "display_name",
				Description: "The display name of the Endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The description of the Endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "deployed_models",
				Description: "The models deployed in this Endpoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "traffic_split",
				Description: "A map from a DeployedModel's ID to the percentage of this Endpoint's traffic that should be forwarded to that DeployedModel..",
				Type:        proto.ColumnType_JSON,
			},
			// {
			// 	Name:        "location_id",
			// 	Description: "The zone where the instance will be provisioned. If not provided, the service will choose a zone from the specified region for the instance.",
			// 	Type:        proto.ColumnType_STRING,
			// },
			{
				Name:        "etag",
				Description: "Used to perform consistent read-modify-write updates. If not set, a blind 'overwrite' update happens.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "labels",
				Description: "The labels with user-defined metadata to organize your Endpoints.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "create_time",
				Description: "Timestamp when this Endpoint was created.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "update_time",
				Description: "Timestamp when this Endpoint was last updated.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "encryption_spec",
				Description: "Customer-managed encryption key spec for an Endpoint. If set, this Endpoint and all sub-resources of this Endpoint will be secured by this key.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "network",
				Description: "The full name of the Google Compute Engine network to which the Endpoint should be peered.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "enable_private_service_connect",
				Description: "If true, expose the Endpoint via private service connect.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "model_deployment_monitoring_job",
				Description: "The time the instance was created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "predict_request_response_logging_config",
				Description: "Configures the request-response logging for online prediction.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard steampipe columns
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
				Transform:   transform.FromP(gcpVertexAITurbotData, "Akas"),
			},

			// Standard gcp columns
			// {
			// 	Name:        "location",
			// 	Description: ColumnDescriptionLocation,
			// 	Type:        proto.ColumnType_STRING,
			// 	Transform:   transform.FromP(gcpVertexAITurbotData, "location"),
			// },
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

//// LIST FUNCTIONS

func listVertexAIEndpoints(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Create Service Connection
	service, err := AIplatformService(ctx, d)
	if err != nil {
		logger.Error("gcp_vertex_ai_endpoint.listVertexAIEndpoints", "connection_error", err)
		return nil, err
	}

	// location := d.EqualsQualString("location")
	matrixLocation := d.EqualsQualString(matrixKeyAIPlatformLocation)
	// Since, when the service API is disabled, matrixLocation value will be nil
	// if location != "" && location != matrixLocation {
	// 	return nil, nil
	// }

	getProjectCached := plugin.HydrateFunc(getProject).WithCache()
	projectId, err := getProjectCached(ctx, d, h)
	if err != nil {
		logger.Error("gcp_vertex_ai_endpoint.listVertexAIEndpoints", "cache_error", err)
		return nil, err
	}
	project := projectId.(string)

	parent := "projects/" + project + "/locations/" + matrixLocation
	req := &aiplatformpb.ListEndpointsRequest{
		Parent: parent,
	}

	it := service.ListEndpoints(ctx, req)
	for {
		resp, err := it.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			logger.Error("gcp_vertex_ai_endpoint.listVertexAIEndpoints", "api_error", err)
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

func getVertexAIEndpoint(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	endpointName := d.EqualsQualString("display_name")

	// Create Service Connection
	service, err := AIplatformService(ctx, d)
	if err != nil {
		logger.Error("gcp_vertex_ai_endpoint.getVertexAIEndpoint", "connection_error", err)
		return nil, err
	}

	location := d.EqualsQualString("location")
	matrixLocation := d.EqualsQualString(matrixKeyAIPlatformLocation)
	// Since, when the service API is disabled, matrixLocation value will be nil
	if location != "" && location != matrixLocation {
		return nil, nil
	}

	getProjectCached := plugin.HydrateFunc(getProject).WithCache()
	projectId, err := getProjectCached(ctx, d, h)
	if err != nil {
		logger.Error("gcp_vertex_ai_endpoint.getVertexAIEndpoint", "cache_error", err)
		return nil, err
	}
	project := projectId.(string)

	name := "projects/" + project + "/locations/" + matrixLocation + "/endpoints/" + endpointName

	req := &aiplatformpb.GetEndpointRequest{
		Name: name,
	}

	op, err := service.GetEndpoint(ctx, req)
	if err != nil {
		logger.Error("gcp_vertex_ai_endpoint.getVertexAIEndpoint", "api_error", err)
		return nil, err
	}

	return op, nil
}

/// TRANSFORM FUNCTIONS

// func gcpVertexAITurbotData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 	endpoint := h.Item.(*aiplatformpb.Endpoint)

// 	project := strings.Split(endpoint.Name, "/")[1]
// 	var location string
// 	matrixLocation := d.EqualsQualString(matrixKeyAIPlatformLocation)
// 	// Since, when the service API is disabled, matrixLocation value will be nil
// 	if matrixLocation != "" {
// 		location = matrixLocation
// 	}

// 	turbotData := map[string]interface{}{
// 		"Project":  project,
// 		"Location": location,
// 		"Akas":     []string{"gcp://aiplatform.googleapis.com/projects/" + project + "/regions/" + location + "/clusters/" + endpoint.Name},
// 	}

// 	return turbotData, nil
// }

func gcpVertexAITurbotData(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	endpoint := d.HydrateItem.(*aiplatform.Endpoint)
	param := d.Param.(string)
	akas := []string{"gcp://aiplatform.googleapis.com/" + endpoint.DisplayName}
	data := make(map[string]interface{}, 0)
	data["akas"] = akas
	return data[param], nil
}

// func gcpRedisInstanceCreateTime(_ context.Context, d *transform.TransformData) (interface{}, error) {
// 	instanceCreateTime := d.HydrateItem.(*redispb.Instance).CreateTime
// 	if instanceCreateTime == nil {
// 		return nil, nil
// 	}
// 	return instanceCreateTime.AsTime(), nil
// }

