package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	compute "google.golang.org/api/compute/v1"
)

//// TABLE DEFINITION

func tableGcpDiskMetricReadOps(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_disk_metric_read_ops",
		Description: "GCP Disk metric read operations",
		List: &plugin.ListConfig{
			ParentHydrate: listComputeInstances,
			Hydrate:       listDiskMetricReadOps,
		},
		Columns: monitoringMetricColumns([]*plugin.Column{
			{
				Name:        "instance_id",
				Description: "The VM Instance name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DimensionValue"),
			},
		}),
	}
}

//// LIST FUNCTION

func listDiskMetricReadOps(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instanceInfo := h.Item.(*compute.Instance)

	dimentionValue := "\"" + instanceInfo.Name + "\""

	return listMonitorMetricStatistics(ctx, d, "5_MIN", "\"compute.googleapis.com/instance/disk/read_ops_count\"", "metric.label.instance_name = ", dimentionValue, instanceInfo.Name)
}
