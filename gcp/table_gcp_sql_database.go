package gcp

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

type sqlDatabaseInfo = struct {
	Database *sqladmin.Database
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
			Tags:       map[string]string{"service": "cloudsql", "action": "databases.get"},
		},
		List: &plugin.ListConfig{
			Hydrate:       listSQLDatabases,
			ParentHydrate: listSQLDatabaseInstances,
			Tags:          map[string]string{"service": "cloudsql", "action": "databases.list"},
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
				Transform:   transform.FromP(sqlInstanceSelfLinkToTurbotData, "Akas"),
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
				Transform:   transform.FromP(sqlInstanceSelfLinkToTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listSQLDatabases(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listSQLDatabases")

	// Get the details of Cloud SQL instance
	instance := h.Item.(*sqladmin.DatabaseInstance)

	// ERROR: rpc error: code = Unknown desc = googleapi: Error 400: Invalid request: Invalid request since instance is not running., invalid
	// Return nil, if the instance not in running state
	if instance.State != "RUNNABLE" {
		return nil, nil
	}

	// Create service connection
	service, err := CloudSQLAdminService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	resp, err := service.Databases.List(project, instance.Name).Do()
	// apply rate limiting
	d.WaitForListRateLimit(ctx)
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
	plugin.Logger(ctx).Trace("getSQLDatabase")

	// Create service connection
	service, err := CloudSQLAdminService(ctx, d)
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
	instanceName := d.EqualsQuals["instance_name"].GetStringValue()

	// Return nil, if no input provided
	if name == "" || instanceName == "" {
		return nil, nil
	}

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

func sqlInstanceSelfLinkToTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(sqlDatabaseInfo)
	param := d.Param.(string)

	project := strings.Split(data.Database.SelfLink, "/")[6]

	turbotData := map[string]interface{}{
		"Project": project,
		"Akas":    []string{"gcp://cloudsql.googleapis.com/projects/" + project + "/instances/" + data.Database.Instance + "/databases/" + data.Database.Name},
	}

	return turbotData[param], nil
}
