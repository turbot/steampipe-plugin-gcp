package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	"google.golang.org/api/compute/v1"
)

//// TABLE DEFINITION

func tableComputeInstanceCpuUtilizationMetricDaily(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_instance_metric_cpu_utilization_daily",
		Description: "GCP Compute Instance Metrics - CPU Utilization (Daily)",
		List: &plugin.ListConfig{
			ParentHydrate: listComputeInstances,
			Hydrate:       listComputeInstanceMetricCpuUtilizationDaily,
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

func listComputeInstanceMetricCpuUtilizationDaily(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instanceInfo := h.Item.(*compute.Instance)

	location := getLastPathElement(instanceInfo.Zone)
	dimensionValue := "\"" + instanceInfo.Name + "\""

	return listMonitorMetricStatistics(ctx, d, "DAILY", "\"compute.googleapis.com/instance/cpu/utilization\"", "metric.labels.instance_name = ", dimensionValue, instanceInfo.Name, location)
}
