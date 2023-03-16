package gcp

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/compute/v1"
)

//// TABLE DEFINITION

func tableGcpComputeDiskMetricWriteOpsDaily(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_disk_metric_write_ops_daily",
		Description: "GCP Compute Disk Metrics - Write Ops (Daily)",
		List: &plugin.ListConfig{
			ParentHydrate: listComputeDisk,
			Hydrate:       listComputeDiskMetricWriteOpsDaily,
			KeyColumns:    plugin.OptionalColumns([]string{"name"}),
		},
		Columns: monitoringMetricColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the disk.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DimensionValue"),
			},
		}),
	}
}

//// LIST FUNCTION

func listComputeDiskMetricWriteOpsDaily(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	diskInfo := h.Item.(*compute.Disk)

	// Get location
	zoneName := getLastPathElement(types.SafeString(diskInfo.Zone))
	regionName := getLastPathElement(types.SafeString(diskInfo.Region))
	location := zoneName
	if zoneName == "" {
		location = regionName
	}
	dimensionValue := "\"" + diskInfo.Name + "\""

	return listMonitorMetricStatistics(ctx, d, h, "DAILY", "\"compute.googleapis.com/instance/disk/write_ops_count\"", "metric.label.device_name = ", dimensionValue, diskInfo.Name, location)
}
