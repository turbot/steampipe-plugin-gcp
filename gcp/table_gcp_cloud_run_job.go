package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"google.golang.org/api/run/v2"
)

//// TABLE DEFINITION

func tableGcpCloudRunJob(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_cloud_run_job",
		Description: "GCP Cloud Run Job",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "location"}),
			Hydrate:    getCloudRunJob,
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudRunJobs,
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "location",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildComputeLocationList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The fully qualified name of this Job.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(lastPathElement),
			},
			{
				Name:        "client",
				Description: "Arbitrary identifier for the API client.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "client_version",
				Description: "Arbitrary version identifier for the API client.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "The creation time.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromGo().NullIfZero(),
			},
			{
				Name:        "delete_time",
				Description: "The deletion time.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromGo().NullIfZero(),
			},
			{
				Name:        "creator",
				Description: "Email address of the authenticated creator.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "etag",
				Description: "A system-generated fingerprint for this version of the resource.",
				Type:        proto.ColumnType_STRING,
			},

			{
				Name:        "execution_count",
				Description: "Number of executions created for this job.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "expire_time",
				Description: "For a deleted resource, the time after which it will be permamently deleted.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromGo().NullIfZero(),
			},
			{
				Name:        "generation",
				Description: "A number that monotonically increases every time the user modifies the desired state.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "self_link",
				Description: "The server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     cloudRunJobSelfLink,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "last_modifier",
				Description: "Email address of the last authenticated modifier.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "launch_stage",
				Description: "The launch stage as defined by Google Cloud Platform Launch Stages (https://cloud.google.com/terms/launch-stages). Cloud Run supports `ALPHA`, `BETA`, and `GA`. If no value is specified, GA is assumed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "observed_generation",
				Description: "The generation of this Job.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "reconciling",
				Description: "Returns true if the Job is currently being acted upon by the system to bring it into the desired state.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "run_execution_token",
				Description: "A unique string used as a suffix for creating a new execution. The Job will become ready when the execution is successfully completed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "satisfies_pzs",
				Description: "Reserved for future use.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "start_execution_token",
				Description: "A unique string used as a suffix creating a new execution. The Job will become ready when the execution is successfully started.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "uid",
				Description: "Server assigned unique identifier for the trigger.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "update_time",
				Description: "The last-modified time.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromGo().NullIfZero(),
			},

			// JSON fields
			{
				Name:        "annotations",
				Description: "Unstructured key value map that may be set by external tools to store and arbitrary metadata.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "binary_authorization",
				Description: "Settings for the Binary Authorization feature.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "conditions",
				Description: "The Conditions of all other associated sub-resources.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "labels",
				Description: "Unstructured key value map that can be used to organize and categorize objects.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "latest_created_execution",
				Description: "The last created execution.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "template",
				Description: "The template used to create executions for this Job.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "terminal_condition",
				Description: "The Condition of this Job, containing its readiness status, and detailed error information in case it did not reach a serving state.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "iam_policy",
				Description: "An Identity and Access Management (IAM) policy, which specifies access controls for Google Cloud resources.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudRunJobIamPolicy,
				Transform:   transform.FromValue(),
			},

			// Standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(cloudRunJobData, "Title"),
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
				Transform:   transform.FromP(cloudRunJobData, "Akas"),
			},

			// Standard GCP columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(cloudRunJobData, "Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(cloudRunJobData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listCloudRunJobs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString("region")

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

	// Create Service Connection
	service, err := CloudRunService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_cloud_run_job.listCloudRunJobs", "service_error", err)
		return nil, err
	}

	// Max limit is set as per documentation
	pageSize := types.Int64(500)
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

	input := "projects/" + project + "/locations/" + location

	resp := service.Projects.Locations.Jobs.List(input).PageSize(*pageSize)
	if err := resp.Pages(ctx, func(page *run.GoogleCloudRunV2ListJobsResponse) error {
		for _, item := range page.Jobs {
			d.StreamListItem(ctx, item)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				page.NextPageToken = ""
				return nil
			}
		}
		return nil
	}); err != nil {
		plugin.Logger(ctx).Error("gcp_cloud_run_job.listCloudRunJobs", "api_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCloudRunJob(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := CloudRunService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_cloud_run_job.getCloudRunJob", "service_error", err)
		return nil, err
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	jobName := d.EqualsQuals["name"].GetStringValue()
	location := d.EqualsQuals["location"].GetStringValue()

	// Empty Check
	if jobName == "" || location == "" {
		return nil, nil
	}

	input := "projects/" + project + "/locations/" + location + "/jobs/" + jobName

	resp, err := service.Projects.Locations.Jobs.Get(input).Do()
	if err != nil {
		plugin.Logger(ctx).Error("gcp_cloud_run_job.getCloudRunJob", "api_error", err)
		return nil, err
	}
	return resp, err
}

func getCloudRunJobIamPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	data := h.Item.(*run.GoogleCloudRunV2Job)
	jobName := strings.Split(data.Name, "/")[5]
	location := strings.Split(data.Name, "/")[3]

	// Create Service Connection
	service, err := CloudRunService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_cloud_run_job.getCloudRunJobIamPolicy", "service_error", err)
		return nil, err
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	input := "projects/" + project + "/locations/" + location + "/jobs/" + jobName

	resp, err := service.Projects.Locations.Jobs.GetIamPolicy(input).Do()
	if err != nil {
		plugin.Logger(ctx).Error("gcp_cloud_run_job.getCloudRunJobIamPolicy", "api_error", err)
		return nil, err
	}

	return resp, err
}

//// TRANSFORM FUNCTIONS

func cloudRunJobSelfLink(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(*run.GoogleCloudRunV2Job)

	var location string
	matrixLocation := d.EqualsQualString(matrixKeyLocation)
	// Since, when the service API is disabled, matrixLocation value will be nil
	if matrixLocation != "" {
		location = matrixLocation
	}

	projectID := strings.Split(data.Name, "/")[1]
	name := strings.Split(data.Name, "/")[5]

	selfLink := "https://run.googleapis.com/v2/projects/" + projectID + "/regions/" + location + "/repositories/" + name

	return selfLink, nil
}

//// TRANSFORM FUNCTIONS

func cloudRunJobData(ctx context.Context, h *transform.TransformData) (interface{}, error) {
	data := h.HydrateItem.(*run.GoogleCloudRunV2Job)
	param := h.Param.(string)

	projectID := strings.Split(data.Name, "/")[1]
	name := strings.Split(data.Name, "/")[5]
	location := strings.Split(data.Name, "/")[3]

	turbotData := map[string]interface{}{
		"Project":  projectID,
		"Title":    name,
		"Location": location,
		"Akas":     []string{"gcp://run.googleapis.com/projects/" + projectID + "/repositories/" + name},
	}

	return turbotData[param], nil
}
