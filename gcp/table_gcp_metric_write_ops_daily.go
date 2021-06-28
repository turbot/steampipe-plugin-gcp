package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	compute "google.golang.org/api/compute/v1"
)

//// TABLE DEFINITION

func tableGcpDiskMetricWriteOpsDaily(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_disk_metric_write_ops_hourly",
		Description: "GCP SQL Database Daily connections",
		List: &plugin.ListConfig{
			ParentHydrate: listComputeInstances,
			Hydrate:       listDiskMetricWriteOpsDaily,
		},
		Columns: monitoringMetricColumns([]*plugin.Column{
			{
				Name:        "instance_id",
				Description: "",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DimensionValue"),
			},
		}),
	}
}

//// LIST FUNCTION

func listDiskMetricWriteOpsDaily(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instanceInfo := h.Item.(*compute.Instance)

	dimentionValue := "\"" + instanceInfo.Name + "\""

	return listMonitorMetricStatistics(ctx, d, "DAILY", "\"compute.googleapis.com/instance/disk/write_ops_count\"", "resource.label.instance_name = ", dimentionValue, "")
}
