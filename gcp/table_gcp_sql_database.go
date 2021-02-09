package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	sql "google.golang.org/api/sql/v1beta4"
)

type sqlDatabaseInfo = struct {
	Database *sql.Database
	Region   string
}

//// TABLE DEFINITION

func tableGcpSQLDatabase(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_sql_database",
		Description: "GCP SQL Database",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "instance_name"}),
			Hydrate:    getSQLDatabase,
		},
		List: &plugin.ListConfig{
			Hydrate:       listSQLDatabases,
			ParentHydrate: listSQLDatabaseInstances,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "A friendly name that identifies the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Database.Name"),
			},
			{
				Name:        "instance_name",
				Description: "The name of the Cloud SQL instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Database.Instance"),
			},
			{
				Name:        "kind",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Database.Kind"),
			},
			{
				Name:        "charset",
				Description: "Specifies the MySQL charset value.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Database.Charset"),
			},
			{
				Name:        "collation",
				Description: "Specifies the MySQL collation value.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Database.Collation"),
			},
			{
				Name:        "self_link",
				Description: "The server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Database.SelfLink"),
			},
			{
				Name:        "sql_server_database_compatibility_level",
				Description: "The version of SQL Server with which the database is to be made compatible.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Database.SqlserverDatabaseDetails.CompatibilityLevel"),
			},
			{
				Name:        "sql_server_database_recovery_model",
				Description: "The recovery model of a SQL Server database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Database.SqlserverDatabaseDetails.RecoveryModel"),
			},

			// standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Database.Name"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(sqlDatabaseAka),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Region"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromConstant(activeProject()),
			},
		},
	}
}

//// LIST FUNCTION

func listSQLDatabases(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listSQLDatabases")

	// Get the details of Cloud SQL instance
	instance := h.Item.(*sql.DatabaseInstance)

	service, err := sql.NewService(ctx)
	if err != nil {
		return nil, err
	}

	project := activeProject()
	resp, err := service.Databases.List(project, instance.Name).Do()
	if err != nil {
		return nil, err
	}
	for _, database := range resp.Items {
		d.StreamLeafListItem(ctx, sqlDatabaseInfo{database, instance.Region})
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSQLDatabase(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	service, err := sql.NewService(ctx)
	if err != nil {
		return nil, err
	}

	name := d.KeyColumnQuals["name"].GetStringValue()
	instanceName := d.KeyColumnQuals["instance_name"].GetStringValue()
	project := activeProject()

	// Get the region where the specified instance is created
	instanceData, err := service.Instances.Get(project, instanceName).Do()
	if err != nil {
		return nil, err
	}
	region := instanceData.Region

	resp, err := service.Databases.Get(project, instanceName, name).Do()
	if err != nil {
		return nil, err
	}

	return sqlDatabaseInfo{resp, region}, nil
}

//// TRANSFORM FUNCTIONS

func sqlDatabaseAka(_ context.Context, d *transform.TransformData) (interface{}, error) {
	database := d.HydrateItem.(sqlDatabaseInfo)

	akas := []string{"gcp://cloudsql.googleapis.com/projects/" + activeProject() + "/instances/" + database.Database.Instance + "/databases/" + database.Database.Name}

	return akas, nil
}
