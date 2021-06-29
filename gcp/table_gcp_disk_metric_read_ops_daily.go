package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	"google.golang.org/api/compute/v1"
)

//// TABLE DEFINITION

func tableGcpDiskMetricReadOpsDaily(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_disk_metric_read_ops_daily",
		Description: "GCP Disk metric read operation",
		List: &plugin.ListConfig{
			ParentHydrate: listComputeDisk,
			Hydrate:       listDiskMetricReadOpsDaily,
		},
		Columns: monitoringMetricColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The Compute Disk name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DimensionValue"),
			},
		}),
	}
}

//// LIST FUNCTION

func listDiskMetricReadOpsDaily(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	diskInfo := h.Item.(*compute.Disk)

	dimentionValue := "\"" + diskInfo.Name + "\""

	return listMonitorMetricStatistics(ctx, d, "DAILY", "\"compute.googleapis.com/instance/disk/read_ops_count\"", "metric.label.device_name = ", dimentionValue, diskInfo.Name)
}
