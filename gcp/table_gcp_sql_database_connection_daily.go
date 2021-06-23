package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

//// TABLE DEFINITION

func tableGcpSQLDatabaseConnectionsMetricDaily(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "table_gcp_sql_database_connection_daily",
		Description: "GCP SQL Database Daily connections",
		List: &plugin.ListConfig{
			ParentHydrate: listSQLDatabaseInstances,
			Hydrate:       listSQLDatabaseMetricConnectionsDaily,
		},
		Columns: monitoringMetricColumn([]*plugin.Column{
			{
				Name:        "instance_id",
				Description: "The SQL Instance name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DimensionValue"),
			},
		}),
	}
}

//// LIST FUNCTION

func listSQLDatabaseMetricConnectionsDaily(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	databaseinfo := h.Item.(*sqladmin.DatabaseInstance)

	// Get project details
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}
	project := projectData.Project
	dimentionValue := "\"" + project + ":" + databaseinfo.Name + "\""

	return listMonitorMetricStatistics(ctx, d, "DAILY", "\"cloudsql.googleapis.com/database/network/connections\"", "resource.label.database_id = ", dimentionValue, databaseinfo.Name)
}
