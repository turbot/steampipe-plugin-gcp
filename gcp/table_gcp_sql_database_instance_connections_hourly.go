package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

//// TABLE DEFINITION

func tableGcpSQLDatabaseInstanceConnectionsMetricHourly(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_sql_database_instance_connections_hourly",
		Description: "GCP SQL Database Instance Hourly connections",
		List: &plugin.ListConfig{
			ParentHydrate: listSQLDatabaseInstances,
			Hydrate:       listSQLDatabaseInstanceMetricConnectionsHourly,
		},
		Columns: monitoringMetricColumns([]*plugin.Column{
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

func listSQLDatabaseInstanceMetricConnectionsHourly(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instanceInfo := h.Item.(*sqladmin.DatabaseInstance)

	// Get project details
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}
	project := projectData.Project
	dimensionValue := "\"" + project + ":" + instanceInfo.Name + "\""

	return listMonitorMetricStatistics(ctx, d, "HOURLY", "\"cloudsql.googleapis.com/database/network/connections\"", "resource.label.database_id = ", dimensionValue, instanceInfo.Name)
}
