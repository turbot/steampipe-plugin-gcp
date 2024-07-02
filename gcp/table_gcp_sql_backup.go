package gcp

import (
	"context"
	"strconv"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

//// TABLE DEFINITION

func tableGcpSQLBackup(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_sql_backup",
		Description: "GCP SQL Backup",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id", "instance_name"}),
			Hydrate:    getSQLBackup,
			Tags:       map[string]string{"service": "cloudsql", "action": "backupRuns.get"},
		},
		List: &plugin.ListConfig{
			Hydrate:       listSQLBackups,
			ParentHydrate: listSQLDatabaseInstances,
			Tags:          map[string]string{"service": "cloudsql", "action": "backupRuns.list"},
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
				Transform:   transform.FromField("EndTime").NullIfZero(),
			},
			{
				Name:        "enqueued_time",
				Description: "Specifies the time when the run was enqueued.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("EnqueuedTime").NullIfZero(),
			},
			{
				Name:        "start_time",
				Description: "Specifies the time when the backup operation actually started.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("StartTime").NullIfZero(),
			},
			{
				Name:        "window_start_time",
				Description: "Specifies the start time of the backup window during which this the backup was attempted.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("WindowStartTime").NullIfZero(),
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
	instance := h.Item.(*sqladmin.DatabaseInstance)

	// Create service connection
	service, err := CloudSQLAdminService(ctx, d)
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

	resp := service.BackupRuns.List(project, instance.Name).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *sqladmin.BackupRunsListResponse) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, backup := range page.Items {
			d.StreamListItem(ctx, backup)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				page.NextPageToken = ""
				return nil
			}
		}
		return nil
	}); err != nil {
		return nil, err
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

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	id := d.EqualsQuals["id"].GetInt64Value()
	instanceName := d.EqualsQuals["instance_name"].GetStringValue()

	resp, err := service.BackupRuns.Get(project, instanceName, id).Do()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func getSQLBackupAka(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	backup := h.Item.(*sqladmin.BackupRun)

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	akas := []string{"gcp://cloudsql.googleapis.com/projects/" + project + "/instances/" + backup.Instance + "/backupRuns/" + strconv.Itoa(int(backup.Id))}

	return akas, nil
}
