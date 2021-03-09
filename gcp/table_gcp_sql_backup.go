package gcp

import (
	"context"
	"strconv"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	sql "google.golang.org/api/sql/v1beta4"
)

//// TABLE DEFINITION

func tableGcpSQLBackup(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_sql_backup",
		Description: "GCP SQL Backup",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id", "instance_name"}),
			Hydrate:    getSQLBackup,
		},
		List: &plugin.ListConfig{
			Hydrate:       listSQLBackups,
			ParentHydrate: listSQLDatabaseInstances,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "An unique identifier for the backup run.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "instance_name",
				Description: "The name of the Cloud SQL instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Instance"),
			},
			{
				Name:        "kind",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A user-specified, human-readable description of the backup run.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "Specifies the status of the backup run.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "end_time",
				Description: "Specifies the time when the backup operation completed.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "enqueued_time",
				Description: "Specifies the time when the run was enqueued.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "start_time",
				Description: "Specifies the time when the backup operation actually started.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "window_start_time",
				Description: "Specifies the start time of the backup window during which this the backup was attempted.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "self_link",
				Description: "The server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "Specifies the type of the backup run. Value can be either 'AUTOMATED' or 'ON_DEMAND'.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "disk_encryption_configuration",
				Description: "Specifies the encryption configuration for the disk.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "disk_encryption_status",
				Description: "Specifies the encryption status of the disk.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "error",
				Description: "Information about why the backup operation failed. This is only present if the run has the FAILED status.",
				Type:        proto.ColumnType_JSON,
			},

			// standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Id"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSQLBackupAka,
				Transform:   transform.FromValue(),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
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

func listSQLBackups(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listSQLBackups")

	// Get the details of Cloud SQL instance
	instance := h.Item.(*sql.DatabaseInstance)

	// Create service connection
	service, err := CloudSQLAdminService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}
	project := projectData.Project

	resp, err := service.BackupRuns.List(project, instance.Name).Do()
	if err != nil {
		return nil, err
	}
	for _, backup := range resp.Items {
		d.StreamLeafListItem(ctx, backup)
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSQLBackup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSQLBackup")

	// Create service connection
	service, err := CloudSQLAdminService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}
	project := projectData.Project

	id := d.KeyColumnQuals["id"].GetInt64Value()
	instanceName := d.KeyColumnQuals["instance_name"].GetStringValue()

	resp, err := service.BackupRuns.Get(project, instanceName, id).Do()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func getSQLBackupAka(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	backup := h.Item.(*sql.BackupRun)

	// Get project details
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}
	project := projectData.Project

	akas := []string{"gcp://cloudsql.googleapis.com/projects/" + project + "/instances/" + backup.Instance + "/backupRuns/" + strconv.Itoa(int(backup.Id))}

	return akas, nil
}
