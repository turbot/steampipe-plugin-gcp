package gcp

import (
	"context"
	"fmt"
	"strconv"
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
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudFunctions,
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
				Transform:   transform.From(getCloudFunctionStatus),
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
				Transform:   transform.FromP(getCloudFunctionBuildConfigData, "RunTime"),
			},

			// other columns
			{
				Name:        "available_memory_mb",
				Description: "The amount of memory in MB available for the function.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromP(getCloudFunctionServiceConfigData, "AvailableMemory"),
			},
			{
				Name:        "build_environment_variables",
				Description: "Environment variables that shall be available during build time",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "build_id",
				Description: "The Cloud Build ID of the latest successful deployment of the function.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getCloudFunctionBuildConfigData, "BuildId"),
			},
			{
				Name:        "entry_point",
				Description: "The name of the function (as defined in source code) that will be executed.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getCloudFunctionBuildConfigData, "EntryPoint"),
			},
			{
				Name:        "environment_variables",
				Description: "Environment variables that shall be available during function execution.",
				Type:        proto.ColumnType_JSON,
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
				Transform:   transform.FromP(getCloudFunctionServiceConfigData, "HttpsTriggers"),
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
				Transform:   transform.FromP(getCloudFunctionServiceConfigData, "IngressSettings"),
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
				Transform:   transform.FromP(getCloudFunctionServiceConfigData, "MaxInstances"),
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
				Transform:   transform.FromP(getCloudFunctionServiceConfigData, "ServiceAccountEmail"),
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
				Name:        "source_upload_url",
				Description: "The Google Cloud Storage signed URL used for source uploading, generated by google.cloud.functions.v1/v2.GenerateUploadUrl",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getCloudFunctionBuildConfigData, "SourceUploadUrl"),
			},
			{
				Name:        "timeout",
				Description: "The function execution timeout. Execution is consideredfailed and can be terminated if the function is not completed at the end of the timeout period. Defaults to 60 seconds.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getCloudFunctionServiceConfigData, "TimeOut"),
			},
			{
				Name:        "update_time",
				Description: "The last update timestamp of the Cloud Function.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "version_id",
				Description: "The version identifier of the Cloud Function. Each deployment attempt results in a new version of a function being created.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "vpc_connector",
				Description: "The VPC Network Connector that this cloud function can  connect to. This field is mutually exclusive with `network` field and will eventually replace it.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getCloudFunctionServiceConfigData, "VpcConnector"),
			},
			{
				Name:        "vpc_connector_egress_settings",
				Description: "The egress settings for the connector, controlling what traffic is diverted through it (VPC_CONNECTOR_EGRESS_SETTINGS_UNSPECIFIED, PRIVATE_RANGES_ONLY, ALL_TRAFFIC).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getCloudFunctionServiceConfigData, "VpcConnectorEgressSettings"),
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

func getCloudFunctionStatus(_ context.Context, d *transform.TransformData) (interface{}, error) {
	cloudFunctionAttributeData := d.HydrateItem.(*cloudfunctions.Function)

	status := cloudFunctionAttributeData.State
	return status, nil
}

func getCloudFunctionBuildConfigData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	cloudFunctionAttributeData := d.HydrateItem.(*cloudfunctions.Function)
	param := d.Param.(string)

	buildConfigs := cloudFunctionAttributeData.BuildConfig

	runTime := buildConfigs.Runtime

	attributesOfBuild := strings.Split(buildConfigs.Build, "/")
	buildId := attributesOfBuild[len(attributesOfBuild)-1]

	entryPoint := buildConfigs.EntryPoint

	sourceUploadUrlInitial := "https://storage.googleapis.com/"
	source := buildConfigs.Source
	storageSource := source.StorageSource
	sourceUploadUrl := sourceUploadUrlInitial + strings.TrimSpace(storageSource.Bucket) + strings.TrimSpace(storageSource.Object)

	buildConfigData := map[string]interface{}{
		"RunTime":         runTime,
		"BuildId":         buildId,
		"EntryPoint":      entryPoint,
		"SourceUploadUrl": sourceUploadUrl,
	}

	return buildConfigData[param], nil
}

func getCloudFunctionServiceConfigData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	cloudFunctionAttributeData := d.HydrateItem.(*cloudfunctions.Function)
	param := d.Param.(string)

	serviceConfigs := cloudFunctionAttributeData.ServiceConfig

	availableMemoryInString := strings.TrimSpace(serviceConfigs.AvailableMemory)
	if cloudFunctionAttributeData.Environment == "GEN_1" {
		availableMemoryInString = availableMemoryInString[0 : len(availableMemoryInString)-1]

	} else {
		availableMemoryInString = availableMemoryInString[0 : len(availableMemoryInString)-2]
	}
	availableMemoryInInt, err := strconv.Atoi(availableMemoryInString)
	if err != nil {
		availableMemoryInInt = 0
	}

	ingressSettings := serviceConfigs.IngressSettings

	httpsTriggers := make(map[string]string)
	httpsTriggers["Url"] = serviceConfigs.Uri
	securityLevel := serviceConfigs.SecurityLevel
	if len(securityLevel) == 0 {
		securityLevel = "SECURITY_LEVEL_UNSPECIFIED"
	}
	httpsTriggers["SecurityLevel"] = securityLevel

	maxInstances := serviceConfigs.MaxInstanceCount

	serviceAccountEmail := serviceConfigs.ServiceAccountEmail

	timeOut := serviceConfigs.TimeoutSeconds

	vpcConnector := serviceConfigs.VpcConnector
	vpcConnectorEgressSettings := serviceConfigs.VpcConnectorEgressSettings

	serviceConfigData := map[string]interface{}{
		"AvailableMemory":            availableMemoryInInt,
		"IngressSettings":            ingressSettings,
		"HttpsTriggers":              httpsTriggers,
		"MaxInstances":               maxInstances,
		"ServiceAccountEmail":        serviceAccountEmail,
		"TimeOut":                    timeOut,
		"VpcConnector":               vpcConnector,
		"VpcConnectorEgressSettings": vpcConnectorEgressSettings,
	}

	return serviceConfigData[param], nil
}
