package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	"google.golang.org/api/compute/v1"
)

//// TABLE DEFINITION

func tableComputeInstanceCpuUtilizationMetricHourly(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_instance_metric_cpu_utilization_hourly",
		Description: "GCP Compute Instance Metrics - CPU Utilization (Hourly)",
		List: &plugin.ListConfig{
			ParentHydrate: listComputeInstances,
			Hydrate:       listComputeInstanceMetricCpuUtilizationHourly,
		},
		Columns: monitoringMetricColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The Compute Instance name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DimensionValue"),
			},
		}),
	}
}

//// LIST FUNCTION

func listComputeInstanceMetricCpuUtilizationHourly(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instanceInfo := h.Item.(*compute.Instance)

	dimensionValue := "\"" + instanceInfo.Name + "\""

	return listMonitorMetricStatistics(ctx, d, "HOURLY", "\"compute.googleapis.com/instance/cpu/utilization\"", "metric.labels.instance_name = ", dimensionValue, instanceInfo.Name)
}
