package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// List of zones where TPUs are available
// Source: https://cloud.google.com/tpu/docs/regions-zones
var tpuSupportedZones = []string{
	"us-central1-a",
	"us-central1-b",
	"us-central1-c",
	"europe-west4-a",
	"asia-east1-c",
}

// BuildTpuZoneList :: return a list of matrix items, one per supported TPU zone
func BuildTpuZoneList(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {
	matrix := make([]map[string]interface{}, len(tpuSupportedZones))
	for i, zone := range tpuSupportedZones {
		matrix[i] = map[string]interface{}{
			"zone": zone,
		}
	}
	return matrix
}
