package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

//// TABLE DEFINITION

func tableGcpSQLDatabaseInstanceMetricConnections(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_sql_database_instance_metric_connections",
		Description: "GCP SQL Database Instance Metrics Connections",
		List: &plugin.ListConfig{
			ParentHydrate: listSQLDatabaseInstances,
			Hydrate:       listSQLDatabaseInstanceMetricConnections,
		},
		Columns: monitoringMetricColumns([]*plugin.Column{
			{
				Name:        "instance_id",
				Description: "The ID of the instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DimensionValue"),
			},
		}),
	}
}

//// LIST FUNCTION

func listSQLDatabaseInstanceMetricConnections(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instanceInfo := h.Item.(*sqladmin.DatabaseInstance)

	// Get project details
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}
	project := projectData.Project
	location := instanceInfo.Region
	dimensionValue := "\"" + project + ":" + instanceInfo.Name + "\""

	return listMonitorMetricStatistics(ctx, d, "FIVE_MINUTES", "\"cloudsql.googleapis.com/database/network/connections\"", "resource.label.database_id = ", dimensionValue, instanceInfo.Name, location)
}
