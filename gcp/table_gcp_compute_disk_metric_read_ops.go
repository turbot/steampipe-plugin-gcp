package gcp

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	"google.golang.org/api/compute/v1"
)

//// TABLE DEFINITION

func tableGcpComputeDiskMetricReadOps(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_disk_metric_read_ops",
		Description: "GCP Compute Disk Metrics - Read Ops",
		List: &plugin.ListConfig{
			ParentHydrate:     listComputeDisk,
			Hydrate:           listComputeDiskMetricReadOps,
			ShouldIgnoreError: isNotFoundError([]string{"403"}),
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

func listComputeDiskMetricReadOps(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	diskInfo := h.Item.(*compute.Disk)

	// Get location
	zoneName := getLastPathElement(types.SafeString(diskInfo.Zone))
	regionName := getLastPathElement(types.SafeString(diskInfo.Region))
	location := zoneName
	if zoneName == "" {
		location = regionName
	}
	dimensionValue := "\"" + diskInfo.Name + "\""

	return listMonitorMetricStatistics(ctx, d, "5_MIN", "\"compute.googleapis.com/instance/disk/read_ops_count\"", "metric.label.device_name = ", dimensionValue, diskInfo.Name, location)
}
