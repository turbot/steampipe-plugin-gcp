package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"google.golang.org/api/composer/v1"
)

//// TABLE DEFINITION

func tableGcpComposerEnvironment(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_composer_environment",
		Description: "GCP Composer Environment",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComposerEnvironment,
		},
		List: &plugin.ListConfig{
			Hydrate: listComposerEnvironments,
		},
		GetMatrixItemFunc: BuildComputeLocationList, // The package at https://pkg.go.dev/google.golang.org/api/composer/v1#ProjectsLocationsService does not provide an API to list all supported regions for the Composer service, so we utilized the `BuildComputeLocationList` function instead.
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The resource name of the environment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "uuid",
				Description: "The UUID (Universally Unique IDentifier) associated with this environment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The current state of the environment.",
				Type:        proto.ColumnType_IPADDR,
			},
			{
				Name:        "create_time",
				Description: "The time at which this environment was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "update_time",
				Description: "The time at which this environment was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "airflow_byoid_uri",
				Description: "The 'bring your own identity' variant of the URI of the Apache Airflow Web UI hosted within this environment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Config.AirflowByoidUri"),
			},
			{
				Name:        "airflow_uri",
				Description: "The URI of the Apache Airflow Web UI hosted within this environment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Config.AirflowUri"),
			},
			{
				Name:        "dag_gcs_prefix",
				Description: "The Cloud Storage prefix of the DAGs for this environment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Config.DagGcsPrefix"),
			},
			{
				Name:        "environment_size",
				Description: "The size of the Cloud Composer environment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Config.EnvironmentSize"),
			},
			{
				Name:        "gke_cluster",
				Description: "The Kubernetes Engine cluster used to run this environment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Config.GkeCluster"),
			},
			{
				Name:        "node_count",
				Description: "The number of nodes in the Kubernetes Engine cluster that will be used to run this environment.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Config.NodeCount"),
			},
			{
				Name:        "resilience_mode",
				Description: "Resilience mode of the cloud composer environment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Config.ResilienceMode"),
			},
			{
				Name:        "data_retention_config",
				Description: "The configuration setting for Airflow database data retention mechanism.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Config.DataRetentionConfig"),
			},
			{
				Name:        "database_config",
				Description: "The configuration settings for Cloud SQL instance used internally by Apache Airflow software.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Config.DatabaseConfig"),
			},
			{
				Name:        "encryption_config",
				Description: "The encryption options for the Cloud Composer environment and its dependencies. Cannot be updated.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Config.EncryptionConfig"),
			},
			{
				Name:        "maintenance_window",
				Description: "The maintenance window is the period when Cloud Composer components may undergo maintenance.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Config.MaintenanceWindow"),
			},
			{
				Name:        "master_authorized_networks_config",
				Description: "The configuration options for GKE cluster master authorized networks.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Config.MasterAuthorizedNetworksConfig"),
			},
			{
				Name:        "node_config",
				Description: "The configuration used for the Kubernetes Engine cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Config.NodeConfig"),
			},
			{
				Name:        "private_environment_config",
				Description: "The configuration used for the Private IP Cloud Composer environment.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Config.PrivateEnvironmentConfig"),
			},
			{
				Name:        "recovery_config",
				Description: "The Recovery settings configuration of an environment.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Config.RecoveryConfig"),
			},
			{
				Name:        "software_config",
				Description: "The configuration settings for software inside the environment.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Config.SoftwareConfig"),
			},
			{
				Name:        "web_server_config",
				Description: "The configuration settings for the Airflow web server App Engine instance.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Config.WebServerConfig"),
			},
			{
				Name:        "web_server_network_access_control",
				Description: "The network-level access control policy for the Airflow web server.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Config.WebServerNetworkAccessControl"),
			},
			{
				Name:        "workloads_config",
				Description: "The workloads configuration settings for the GKE cluster associated with the Cloud Composer environment.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Config.WorkloadsConfig"),
			},

			// standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(lastPathElement),
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
				Transform:   transform.FromP(composerEnvironmentTurbotData, "Akas"),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(composerEnvironmentTurbotData, "Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(composerEnvironmentTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listComposerEnvironments(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	location := d.EqualsQualString(matrixKeyLocation)

	// Create Service Connection
	service, err := ComposerService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_composer_environment.listComposerEnvironments", "service_error", err)
		return nil, err
	}

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

	parent := "projects/" + project + "/locations/" + location

	resp := service.Projects.Locations.Environments.List(parent). Context(ctx).PageSize(*pageSize)
	if err := resp.Pages(ctx, func(page *composer.ListEnvironmentsResponse) error {
		for _, item := range page.Environments {
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
		plugin.Logger(ctx).Error("gcp_composer_environment.listComposerEnvironments", "api_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getComposerEnvironment(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	location := d.EqualsQualString(matrixKeyLocation)

	// Create Service Connection
	service, err := ComposerService(ctx, d)
	if err != nil {
		return nil, err
	}

	name := d.EqualsQuals["name"].GetStringValue()

	// Empty check
	if name == "" {
		return nil, nil
	}

	// Restrict API call for other location
	if len(strings.Split(name, "/")) > 3 && strings.Split(name, "/")[3] != location {
		return nil, nil
	}

	resp, err := service.Projects.Locations.Environments.Get(name).Do()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

//// TRANSFORM FUNCTIONS

func composerEnvironmentTurbotData(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	env := d.HydrateItem.(*composer.Environment)

	param := d.Param.(string)
	splitEnv := strings.Split(env.Name, "/")

	turbotData := map[string]interface{}{
		"Project":  splitEnv[1],
		"Location": splitEnv[3],
		"SelfLink": "https://composer.googleapis.com/v1/" + env.Name,
		"Akas":     []string{"gcp://composer.googleapis.com/" + env.Name},
	}

	return turbotData[param], nil
}
