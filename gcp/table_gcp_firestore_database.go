package gcp

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"

	"google.golang.org/api/firestore/v1"
)

//// TABLE DEFINITION

func tableGcpFirestoreDatabase(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_firestore_database",
		Description: "GCP Firestore Database",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getFirestoreDatabase,
		},
		List: &plugin.ListConfig{
			Hydrate: listFirestoreDatabases,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "show_deleted", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "show_deleted",
				Type:        proto.ColumnType_BOOL,
				Description: "Set to true to include deleted databases in the query results. By default, deleted databases are excluded.",
				Transform:   transform.FromQual("show_deleted"),
				Default:     false,
			},
			{
				Name:        "app_engine_integration_mode",
				Description: "The App Engine integration mode to use for this database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "concurrency_mode",
				Description: "The concurrency control mode to use for this database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "The time at which the database was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "delete_protection_state",
				Description: "The delete protection state of the database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "delete_time",
				Description: "The time at which the database was deleted.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("DeleteTime").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "earliest_version_time",
				Description: "The earliest timestamp at which older versions of the data can be read.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "etag",
				Description: "This checksum is computed by the server based on the value of other fields.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "key_prefix",
				Description: "The key prefix for this database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "The name of the database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "point_in_time_recovery_enablement",
				Description: "Whether to enable the PITR feature on this database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "previous_id",
				Description: "The database resource's prior database ID. This field is only populated for deleted databases.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(firestoreDatabaseSelfLink),
			},
			{
				Name:        "type",
				Description: "The type of the database. Can be FIRESTORE_NATIVE or DATASTORE_MODE.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "uid",
				Description: "The system-generated UUID4 for this Database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "update_time",
				Description: "The time at which the database was most recently updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "version_retention_period",
				Description: "The period during which past versions of data are retained in the database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cmek_config",
				Description: "The CMEK (Customer Managed Encryption Key) configuration for the database.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "source_info",
				Description: "Information about the provenance of this database.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(firestoreDatabaseTurbotData, "Title"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(firestoreDatabaseTurbotData, "Akas"),
			},

			// Standard GCP columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LocationId"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(firestoreDatabaseTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listFirestoreDatabases(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := FirestoreDatabaseService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_firestore_database.listFirestoreDatabases", "service_err", err)
		return nil, err
	}

	// Get project details
	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	// Get showDeleted qual
	showDeleted := false
	if v, ok := d.EqualsQuals["show_deleted"]; ok {
		showDeleted = v.GetBoolValue()
	}

	// List databases
	parent := "projects/" + project
	resp, err := service.Projects.Databases.List(parent).ShowDeleted(showDeleted).Do()
	if err != nil {
		plugin.Logger(ctx).Error("gcp_firestore_database.listFirestoreDatabases", "api", err)
		return nil, err
	}

	for _, database := range resp.Databases {
		d.StreamListItem(ctx, database)
		// Check if context has been cancelled or if the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			break
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getFirestoreDatabase(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := FirestoreDatabaseService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_firestore_database.getFirestoreDatabase", "service_err", err)
		return nil, err
	}

	// Get project details
	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	name := d.EqualsQuals["name"].GetStringValue()
	if name == "" {
		return nil, nil
	}

	// If the name doesn't start with "projects/", add the prefix
	if !strings.HasPrefix(name, "projects/") {
		name = "projects/" + project + "/databases/" + name
	}

	// Get database info
	database, err := service.Projects.Databases.Get(name).Do()
	if err != nil {
		plugin.Logger(ctx).Error("gcp_firestore_database.getFirestoreDatabase", "api", err)
		return nil, err
	}

	return database, nil
}

//// TRANSFORM FUNCTIONS

func firestoreDatabaseTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*firestore.GoogleFirestoreAdminV1Database)
	param := d.Param.(string)

	split := strings.Split(data.Name, "/")
	project := split[1]
	title := split[3]

	turbotData := map[string]interface{}{
		"Akas":    []string{"gcp://firestore.googleapis.com/" + data.Name},
		"Project": project,
		"Title":   title,
	}

	return turbotData[param], nil
}

func firestoreDatabaseSelfLink(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*firestore.GoogleFirestoreAdminV1Database)

	selfLink := "https://firestore.googleapis.com/v1/" + data.Name
	return selfLink, nil
}
