package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/compute/v1"
)

//// TABLE DEFINITION

func tableGcpComputeInstanceMetricCpuUtilization(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_instance_metric_cpu_utilization",
		Description: "GCP Compute Instance Metrics - CPU Utilization",
		List: &plugin.ListConfig{
			ParentHydrate: listComputeInstances,
			Hydrate:       listComputeInstanceMetricCpuUtilization,
			KeyColumns:    plugin.OptionalColumns([]string{"name"}),
			Tags:          map[string]string{"service": "monitoring", "action": "timeSeries.list"},
		},
		Columns: monitoringMetricColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DimensionValue"),
			},
		}),
	}
}

//// LIST FUNCTION

func listComputeInstanceMetricCpuUtilization(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instanceInfo := h.Item.(*compute.Instance)

	location := getLastPathElement(instanceInfo.Zone)
	dimensionValue := "\"" + instanceInfo.Name + "\""

	return listMonitorMetricStatistics(ctx, d, h, "FIVE_MIMUTES", "\"compute.googleapis.com/instance/cpu/utilization\"", "metric.labels.instance_name = ", dimensionValue, instanceInfo.Name, location)
}
