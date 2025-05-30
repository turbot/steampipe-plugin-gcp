package gcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/cloudfunctions/v2"
)

func tableGcpCloudfunctionFunction(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_cloudfunctions_function",
		Description: "GCP Cloud Function",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "location"}),
			Hydrate:    getCloudFunction,
			Tags:       map[string]string{"service": "cloudfunctions", "action": "functions.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudFunctions,
			Tags:    map[string]string{"service": "cloudfunctions", "action": "functions.list"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the function.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel().Transform(lastPathElement),
			},
			{
				Name:        "status",
				Description: "Status of the function deployment (ACTIVE, OFFLINE, CLOUD_FUNCTION_STATUS_UNSPECIFIED,DEPLOY_IN_PROGRESS, DELETE_IN_PROGRESS, UNKNOWN).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("State"),
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(cloudFunctionSelfLink),
			},
			{
				Name:        "description",
				Description: "User-provided description of a function.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "runtime",
				Description: "The runtime in which to run the function.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BuildConfig.Runtime"),
			},

			// other columns
			{
				Name:        "available_memory_mb",
				Description: "The amount of memory in MB available for the function.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceConfig.AvailableMemory"),
			},
			{
				Name:        "build_environment_variables",
				Description: "Environment variables that shall be available during build time",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("BuildConfig.EnvironmentVariables"),
			},
			{
				Name:        "build_id",
				Description: "The Cloud Build ID of the latest successful deployment of the function.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BuildConfig.Build"),
			},
			{
				Name:        "entry_point",
				Description: "The name of the function (as defined in source code) that will be executed.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BuildConfig.EntryPoint"),
			},
			{
				Name:        "service_environment_variables",
				Description: "Environment variables that shall be available during function execution.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ServiceConfig.EnvironmentVariables"),
			},
			{
				Name:        "event_trigger",
				Description: "A source that fires events in response to a condition in another service.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "https_trigger",
				Description: "An HTTPS endpoint type of source that can be triggered via URL.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "iam_policy",
				Description: "The IAM policy for the function.", Transform: transform.FromValue(), Hydrate: getGcpCloudFunctionIamPolicy,
				Type: proto.ColumnType_JSON,
			},
			{
				Name:        "ingress_settings",
				Description: "The ingress settings for the function, controlling what traffic can reach it (INGRESS_SETTINGS_UNSPECIFIED, ALLOW_ALL, ALLOW_INTERNAL_ONLY, ALLOW_INTERNAL_AND_GCLB).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceConfig.IngressSettings"),
			},
			{
				Name:        "labels",
				Description: "Labels that apply to this function.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "max_instances",
				Description: "The limit on the maximum number of function instances that may coexist at a given time. In some cases, such as rapid traffic surges, Cloud Functions may, for a short period of time, create more instances than the specified max instances limit.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ServiceConfig.MaxInstanceCount"),
			},
			{
				Name:        "network",
				Description: "The VPC Network that this cloud function can connect to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_account_email",
				Description: "The email of the function's service account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceConfig.ServiceAccountEmail"),
			},
			{
				Name:        "build_source",
				Description: "The location of the function source code.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ServiceConfig.Source"),
			},
			{
				Name:        "source_archive_url",
				Description: "The Google Cloud Storage URL, starting with gs://, pointing to the zip archive which contains the function.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_repository",
				Description: "**Beta Feature** The source repository where a function is hosted.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_timeout",
				Description: "The function execution timeout. Execution is consideredfailed and can be terminated if the function is not completed at the end of the timeout period. Defaults to 60 seconds.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceConfig.TimeoutSeconds"),
			},
			{
				Name:        "update_time",
				Description: "The last update timestamp of the Cloud Function.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "vpc_connector",
				Description: "The VPC Network Connector that this cloud function can  connect to. This field is mutually exclusive with `network` field and will eventually replace it.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceConfig.VpcConnector"),
			},
			{
				Name:        "vpc_connector_egress_settings",
				Description: "The egress settings for the connector, controlling what traffic is diverted through it (VPC_CONNECTOR_EGRESS_SETTINGS_UNSPECIFIED, PRIVATE_RANGES_ONLY, ALL_TRAFFIC).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceConfig.VpcConnectorEgressSettings"),
			},
			{
				Name:        "kms_key_name",
				Description: "Resource name of a KMS crypto key (managed by the user) used to encrypt/decrypt function resources.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "url",
				Description: "The deployed URL for the function.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service",
				Description: "Name of the service associated with the function.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceConfig.Service"),
			},
			{
				Name:        "service_revision",
				Description: "The name of the service revision.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceConfig.Revision"),
			},
			{
				Name:        "service_all_traffic_on_latest_revision",
				Description: "Whether 100% of traffic is routed to the latest revision.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ServiceConfig.AllTrafficOnLatestRevision"),
			},
			{
				Name:        "service_available_cpu",
				Description: "The number of CPUs used in a single container instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceConfig.AvailableCpu"),
			},
			{
				Name:        "max_instance_request_concurrency",
				Description: "The maximum number of concurrent requests that each instance can receive.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ServiceConfig.MaxInstanceRequestConcurrency"),
			},
			{
				Name:        "min_instances",
				Description: "The minimum number of function instances that may coexist at a given time.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ServiceConfig.MinInstanceCount"),
			},
			{
				Name:        "satisfies_pzs",
				Description: "Reserved for future use.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "service_security_level",
				Description: "Configure whether the function only accepts HTTPS.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceConfig.SecurityLevel"),
			},
			{
				Name:        "service_secret_environment_variables",
				Description: "Secret environment variables configuration.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ServiceConfig.SecretEnvironmentVariables"),
			},
			{
				Name:        "service_secret_volumes",
				Description: "Secret volumes configuration.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ServiceConfig.SecretVolumes"),
			},
			{
				Name:        "state_messages",
				Description: "State Messages for this Cloud Function.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "build_config",
				Description: "Describes the Build step of the function that builds a container from the given source.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "service_config",
				Description: "Describes the Service being deployed. Currently deploys services to Cloud Run (fully managed).",
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
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Labels"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(gcpCloudFunctionTurbotData, "Akas"),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(locationFromFunctionName),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpCloudFunctionTurbotData, "Project"),
			},
		},
	}
}

//// HYDRATE FUNCTIONS

func listCloudFunctions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listCloudFunctions")

	// Create Service Connection
	service, err := CloudFunctionsService(ctx, d)
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

	data := "projects/" + project + "/locations/-" // '-' for all locations...

	resp := service.Projects.Locations.Functions.List(data).PageSize(*pageSize)
	if err := resp.Pages(
		ctx,
		func(page *cloudfunctions.ListFunctionsResponse) error {
			// apply rate limiting
			d.WaitForListRateLimit(ctx)

			for _, item := range page.Functions {
				d.StreamListItem(ctx, item)

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

	return nil, nil
}

func getCloudFunction(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("GetCloudFunction")

	// Create Service Connection
	service, err := CloudFunctionsService(ctx, d)
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
	location := d.EqualsQuals["location"].GetStringValue()

	// API https://cloud.google.com/functions/docs/reference/rest/v2/projects.locations.functions/get
	cloudFunction, err := service.Projects.Locations.Functions.Get("projects/" + project + "/locations/" + location + "/functions/" + name).Do()
	if err != nil {
		return nil, err
	}

	return cloudFunction, nil
}

func getGcpCloudFunctionIamPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getGcpCloudFunctionIamPolicy")

	// Create Service Connection
	service, err := CloudFunctionsService(ctx, d)
	if err != nil {
		return nil, err
	}

	function := h.Item.(*cloudfunctions.Function)

	resp, err := service.Projects.Locations.Functions.GetIamPolicy(function.Name).Do()
	if err != nil {
		return nil, err
	}

	if resp != nil {
		return resp, nil
	}

	return cloudfunctions.Policy{}, nil
}

//// TRANSFORM FUNCTIONS

func gcpCloudFunctionTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	function := d.HydrateItem.(*cloudfunctions.Function)
	param := d.Param.(string)

	project := strings.Split(function.Name, "/")[1]

	turbotData := map[string]interface{}{
		"Project": project,
		"Akas":    []string{"gcp://cloudfunctions.googleapis.com/" + function.Name},
	}

	return turbotData[param], nil
}

func locationFromFunctionName(_ context.Context, d *transform.TransformData) (interface{}, error) {
	functionName := types.SafeString(d.Value)
	parts := strings.Split(functionName, "/")
	if len(parts) != 6 {
		return nil, fmt.Errorf("transform locationFromFunctionName failed - unexpected name format: %s", functionName)
	}
	return parts[3], nil
}

func cloudFunctionSelfLink(_ context.Context, d *transform.TransformData) (interface{}, error) {
	cloudFunctionAttributeData := d.HydrateItem.(*cloudfunctions.Function)
	var selfLink string
	if cloudFunctionAttributeData.Environment == "GEN_1" {
		selfLink = "https://cloudfunctions.googleapis.com/v1/" + cloudFunctionAttributeData.Name
	} else {
		selfLink = "https://cloudfunctions.googleapis.com/v2/" + cloudFunctionAttributeData.Name
	}

	return selfLink, nil
}