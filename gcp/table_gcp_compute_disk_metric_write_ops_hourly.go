package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	compute "google.golang.org/api/compute/v1"
)

//// TABLE DEFINITION

func tableGcpComputeDiskMetricWriteOpsHourly(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_disk_metric_write_ops_hourly",
		Description: "GCP Disk metric write operations",
		List: &plugin.ListConfig{
			ParentHydrate: listComputeDisk,
			Hydrate:       listDiskMetricWriteOpsHourly,
		},
		Columns: monitoringMetricColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the Disk.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DimensionValue"),
			},
		}),
	}
}

//// LIST FUNCTION

func listDiskMetricWriteOpsHourly(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	diskInfo := h.Item.(*compute.Disk)

	dimentionValue := "\"" + diskInfo.Name + "\""

	return listMonitorMetricStatistics(ctx, d, "HOURLY", "\"compute.googleapis.com/instance/disk/write_ops_count\"", "metric.label.device_name = ", dimentionValue, diskInfo.Name)
}