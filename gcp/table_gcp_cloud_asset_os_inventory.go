package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"

	"google.golang.org/api/cloudasset/v1"
)

//// TABLE DEFINITION

func tableGcpCloudAssetOsInventory(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_cloud_asset_os_inventory",
		Description: "GCP Cloud Asset OS Inventory",
		List: &plugin.ListConfig{
			Hydrate: listCloudAssetOsInventories,
			Tags:    map[string]string{"service": "cloudasset", "action": "assets.listOsInventory"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The full name of the asset.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "asset_type",
				Description: "The type of the asset.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "update_time",
				Description: "The last update timestamp of an asset.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "hostname",
				Description: "The VM hostname.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("OsInventory.OsInfo.Hostname"),
			},
			{
				Name:        "os_long_name",
				Description: "The operating system long name, e.g., Debian GNU/Linux 9 (stretch).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("OsInventory.OsInfo.LongName"),
			},
			{
				Name:        "os_short_name",
				Description: "The operating system short name, e.g., windows or debian.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("OsInventory.OsInfo.ShortName"),
			},
			{
				Name:        "os_version",
				Description: "The version of the operating system.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("OsInventory.OsInfo.Version"),
			},
			{
				Name:        "architecture",
				Description: "The system architecture of the operating system.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("OsInventory.OsInfo.Architecture"),
			},
			{
				Name:        "kernel_version",
				Description: "The kernel version of the operating system.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("OsInventory.OsInfo.KernelVersion"),
			},
			{
				Name:        "kernel_release",
				Description: "The kernel release of the operating system.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("OsInventory.OsInfo.KernelRelease"),
			},
			{
				Name:        "osconfig_agent_version",
				Description: "The current version of the OS Config agent running on the VM.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("OsInventory.OsInfo.OsconfigAgentVersion"),
			},
			{
				Name:        "inventory_update_time",
				Description: "Timestamp of the last reported inventory for the VM.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("OsInventory.UpdateTime"),
			},
			{
				Name:        "items",
				Description: "Inventory items related to the VM keyed by an opaque unique identifier for each inventory item, such as installed packages and available patches.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("OsInventory.Items"),
			},
			{
				Name:        "ancestors",
				Description: "The ancestry path of an asset in Google Cloud resource hierarchy.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},

			// Standard GCP columns
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Hydrate:     getProject,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listCloudAssetOsInventories(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return listCloudAssetsByContentType(ctx, d, h, "OS_INVENTORY", "gcp_cloud_asset_os_inventory.listCloudAssetOsInventories", func(ctx context.Context, d *plugin.QueryData, asset *cloudasset.Asset) {
		d.StreamListItem(ctx, asset)
	})
}
