package gcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/metastore/v1"
)

func tableGcpDataprocMetastoreService(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_dataproc_metastore_service",
		Description: "GCP Dataproc Metastore Service",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getDataprocMetastoreService,
		},
		List: &plugin.ListConfig{
			Hydrate: listDataprocMetastoreServices,
			KeyColumns: plugin.KeyColumnSlice{
				// String columns
				{Name: "state", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "database_type", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "uid", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "release_channel", Require: plugin.Optional, Operators: []string{"="}},
			},
		},
		GetMatrixItemFunc: BuildDataprocMetastoreLocationList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The relative resource name of the metastore service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "artifact_gcs_uri",
				Description: "A Cloud Storage URI (starting with gs://) that specifies where artifacts related to the metastore service are stored.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "The time when the metastore service was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "database_type",
				Description: "The database type that the Metastore service stores its data.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "endpoint_uri",
				Description: "The URI of the endpoint used to access the metastore service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tier",
				Description: "The tier of the service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "uid",
				Description: "The globally unique resource identifier of the metastore service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "update_time",
				Description: "The time when the metastore service was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "network",
				Description: "The relative resource name of the VPC network on which the instance can be accessed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "port",
				Description: "The TCP port at which the metastore service is reached. Default: 9083.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "release_channel",
				Description: "The release channel of the service. If unspecified, defaults to STABLE.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The current state of the metastore service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state_message",
				Description: "Additional information about the current state of the metastore service, if available.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     dataprocMetastoreServiceSelfLink,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "telemetry_config",
				Description: "The configuration specifying telemetry settings for the Dataproc Metastore service.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "network_config",
				Description: "The configuration specifying the network settings for the Dataproc Metastore service.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "scaling_config",
				Description: "Scaling configuration of the metastore service.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "hive_metastore_config",
				Description: "Configuration information specific to running Hive metastore software as the metastore service.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "maintenance_window",
				Description: "The one-hour maintenance window of the metastore service.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("MaintenanceWindow"),
			},
			{
				Name:        "encryption_config",
				Description: "Information used to configure the Dataproc Metastore service to encrypt customer data at rest.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("EncryptionConfig"),
			},
			{
				Name:        "metadata_integration",
				Description: "The setting that defines how metastore metadata should be integrated with external services and systems.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("MetadataIntegration"),
			},
			{
				Name:        "metadata_management_activity",
				Description: "The metadata management activities of the metastore service.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("MetadataManagementActivity"),
			},

			// Steampipe standard columns
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
				Hydrate:     gcpMetastoreServiceTurbotData,
				Transform:   transform.FromField("Akas"),
			},

			// Standard GCP columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Hydrate:     gcpMetastoreServiceTurbotData,
				Transform:   transform.FromField("Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Hydrate:     gcpMetastoreServiceTurbotData,
				Transform:   transform.FromField("Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listDataprocMetastoreServices(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var location string
	matrixLocation := d.EqualsQualString(matrixKeyLocation)
	// Since, when the service API is disabled, matrixLocation value will be nil
	if matrixLocation != "" {
		location = matrixLocation
	}

	// Create Service Connection
	service, err := DataprocMetastoreService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_dataproc_metastore_service.listDataprocMetastoreServices", "connection_error", err)
		return nil, err
	}

	var filters []string

	if d.EqualsQualString("state") != "" {
		filters = append(filters, fmt.Sprintf("state=\"%s\"", d.EqualsQualString("state")))
	}
	if d.EqualsQualString("database_type") != "" {
		filters = append(filters, fmt.Sprintf("databaseType=\"%s\"", d.EqualsQualString("database_type")))
	}
	if d.EqualsQualString("uid") != "" {
		filters = append(filters, fmt.Sprintf("uid=\"%s\"", d.EqualsQualString("uid")))
	}
	if d.EqualsQualString("release_channel") != "" {
		filters = append(filters, fmt.Sprintf("releaseChannel=\"%s\"", d.EqualsQualString("release_channel")))
	}

	filterString := ""
	if len(filters) > 0 {
		filterString = strings.Join(filters, "AND")
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

	// projects/{project_number}/locations/{location_id}
	parent := "projects/" + project + "/locations/" + location

	resp := service.Projects.Locations.Services.List(parent).PageSize(*pageSize).Filter(filterString)
	if err := resp.Pages(ctx, func(page *metastore.ListServicesResponse) error {
		for _, service := range page.Services {
			d.StreamListItem(ctx, service)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				page.NextPageToken = ""
				return nil
			}
		}
		return nil
	}); err != nil {
		plugin.Logger(ctx).Error("gcp_dataproc_metastore_service.listDataprocMetastoreServices", "api_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getDataprocMetastoreService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	matrixLocation := d.EqualsQualString(matrixKeyLocation)
	name := d.EqualsQuals["name"].GetStringValue()

	if len(name) < 1 {
		return nil, nil
	}

	// Restrict the APi call for other locations
	if len(strings.Split(name, "/")) > 3 && strings.Split(name, "/")[3] != matrixLocation {
		return nil, nil
	}
	// Create Service Connection
	service, err := DataprocMetastoreService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_dataproc_metastore_service.getDataprocMetastoreService", "connection_error", err)
		return nil, err
	}

	resp, err := service.Projects.Locations.Services.Get(name).Do()
	if err != nil {
		plugin.Logger(ctx).Error("gcp_dataproc_metastore_service.getDataprocMetastoreService", "api_error", err)
		return nil, err
	}

	return resp, nil
}

func gcpMetastoreServiceTurbotData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	service := h.Item.(*metastore.Service)

	splitName := strings.Split(service.Name, "/")

	turbotData := map[string]interface{}{
		"Project":  splitName[1],
		"Location": splitName[3],
		"Akas":     []string{"gcp://metastore.googleapis.com/" + service.Name},
	}

	return turbotData, nil
}

func dataprocMetastoreServiceSelfLink(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	service := h.Item.(*metastore.Service)

	selfLink := "https://metastore.googleapis.com/v1/" + service.Name

	return selfLink, nil
}
