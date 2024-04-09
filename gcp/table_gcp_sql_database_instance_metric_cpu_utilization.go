package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

//// TABLE DEFINITION

func tableGcpSQLDatabaseInstanceMetricCpuUtilization(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_sql_database_instance_metric_cpu_utilization",
		Description: "GCP SQL Database Instance Metrics - CPU Utilization",
		List: &plugin.ListConfig{
			ParentHydrate: listSQLDatabaseInstances,
			Hydrate:       listSQLDatabaseInstanceMetricCpuUtilization,
			Tags:          map[string]string{"service": "monitoring", "action": "timeSeries.list"},
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

func listSQLDatabaseInstanceMetricCpuUtilization(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instanceInfo := h.Item.(*sqladmin.DatabaseInstance)

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)
	location := instanceInfo.Region
	dimensionValue := "\"" + project + ":" + instanceInfo.Name + "\""

	return listMonitorMetricStatistics(ctx, d, h, "FIVE_MINUTE", "\"cloudsql.googleapis.com/database/cpu/utilization\"", "resource.label.database_id = ", dimensionValue, instanceInfo.Name, location)
}
