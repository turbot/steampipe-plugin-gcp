package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/alloydb/v1"
)

func tableGcpAlloyDBCluster(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_alloydb_cluster",
		Description: "GCP AlloyDB Cluster",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("display_name"),
			Hydrate:    getAlloydbCluster,
		},
		List: &plugin.ListConfig{
			Hydrate: listAlloydbClusters,
			KeyColumns: plugin.KeyColumnSlice{
				// String columns
				{Name: "location", Require: plugin.Optional, Operators: []string{"="}},
			},
		},
		GetMatrixItemFunc: BuildAlloyDBLocationList,
		Columns: []*plugin.Column{
			{
				Name:        "display_name",
				Type:        proto.ColumnType_STRING,
				Description: "User-settable and human-readable display name for the Cluster.",
				Transform:   transform.FromField("Name").Transform(getAlloyDBClusterDisplayName),
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the cluster resource.",
			},
			{
				Name:        "state",
				Type:        proto.ColumnType_STRING,
				Description: "The current serving state of the cluster.",
			},
			{
				Name:        "uid",
				Type:        proto.ColumnType_STRING,
				Description: "The system-generated UID of the resource.",
			},
			{
				Name:        "cluster_type",
				Type:        proto.ColumnType_STRING,
				Description: "The type of the cluster.",
			},
			{
				Name:        "update_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The update timestamp.",
				Transform:   transform.FromField("UpdateTime").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "create_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The create timestamp.",
				Transform:   transform.FromField("CreateTime").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "database_version",
				Type:        proto.ColumnType_STRING,
				Description: "The database engine major version.",
			},
			{
				Name:        "delete_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The delete timestamp.",
				Transform:   transform.FromField("DeleteTime").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "encryption_config",
				Type:        proto.ColumnType_JSON,
				Description: "Encryption config to encrypt the data disks.",
			},
			{
				Name:        "encryption_info",
				Type:        proto.ColumnType_JSON,
				Description: "The encryption information for the cluster.",
			},
			{
				Name:        "etag",
				Type:        proto.ColumnType_STRING,
				Description: "For Resource freshness validation.",
			},
			{
				Name:        "network",
				Type:        proto.ColumnType_STRING,
				Description: "The resource link for the VPC network.",
			},
			{
				Name:        "reconciling",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates if the cluster's current state does not match the intended state.",
			},
			{
				Name:        "satisfies_pzs",
				Type:        proto.ColumnType_BOOL,
				Description: "Reserved for future use.",
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     alloyDBClusterSelfLink,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "secondary_config",
				Type:        proto.ColumnType_JSON,
				Description: "Cross Region replication config specific to SECONDARY cluster.",
			},
			{
				Name:        "ssl_config",
				Type:        proto.ColumnType_JSON,
				Description: "SSL configuration for this AlloyDB cluster.",
			},
			{
				Name:        "annotations",
				Type:        proto.ColumnType_JSON,
				Description: "Annotations to allow client tools to store a small amount of arbitrary data.",
			},
			{
				Name:        "automated_backup_policy",
				Type:        proto.ColumnType_JSON,
				Description: "The automated backup policy for this cluster.",
			},
			{
				Name:        "backup_source",
				Type:        proto.ColumnType_JSON,
				Description: "Cluster created from backup.",
			},
			{
				Name:        "continuous_backup_config",
				Type:        proto.ColumnType_JSON,
				Description: "Continuous backup configuration for this cluster.",
			},
			{
				Name:        "continuous_backup_info",
				Type:        proto.ColumnType_JSON,
				Description: "Continuous backup properties for this cluster.",
			},
			{
				Name:        "network_config",
				Type:        proto.ColumnType_JSON,
				Description: "Network configuration details.",
			},
			{
				Name:        "primary_config",
				Type:        proto.ColumnType_JSON,
				Description: "Cross Region replication config specific to PRIMARY cluster.",
			},
			{
				Name:        "initial_user",
				Type:        proto.ColumnType_JSON,
				Description: "Initial user to setup during cluster creation.",
			},
			{
				Name:        "labels",
				Type:        proto.ColumnType_JSON,
				Description: "Labels as key-value pairs.",
			},
			{
				Name:        "migration_source",
				Type:        proto.ColumnType_JSON,
				Description: "Cluster created via DMS migration.",
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(getAlloyDBClusterDisplayName),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     gcpAlloyDBClusterTurbotData,
				Transform:   transform.FromField("Akas"),
			},

			// Standard GCP columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Hydrate:     gcpAlloyDBClusterTurbotData,
				Transform:   transform.FromField("Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Hydrate:     gcpAlloyDBClusterTurbotData,
				Transform:   transform.FromField("ProjectId"),
			},
		},
	}
}

//// LIST FUNCTION

func listAlloydbClusters(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var location string
	matrixLocation := d.EqualsQualString(matrixKeyLocation)
	// Since, when the service API is disabled, matrixLocation value will be nil
	if matrixLocation != "" {
		location = matrixLocation
	}

	// Create Service Connection
	service, err := AlloyDBService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_alloydb_cluster.listAlloydbClusters", "connection_error", err)
		return nil, err
	}

	// Max limit is set as per documentation
	pageSize := types.Int64(500)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *pageSize {
			pageSize = limit
		}
	}

	// Minimize the API call with given location
	region := d.EqualsQualString("location")
	if region != "" && region != location {
		return nil, nil
	}

	// Get project details
	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	resp := service.Projects.Locations.Clusters.List("projects/" + project + "/locations/" + location).PageSize(*pageSize)
	if err := resp.Pages(ctx, func(page *alloydb.ListClustersResponse) error {
		for _, cluster := range page.Clusters {
			d.StreamListItem(ctx, cluster)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				page.NextPageToken = ""
				return nil
			}
		}
		return nil
	}); err != nil {
		plugin.Logger(ctx).Error("gcp_alloydb_cluster.listAlloydbClusters", "api_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTION

func getAlloydbCluster(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var location string
	matrixLocation := d.EqualsQualString(matrixKeyLocation)
	// Since, when the service API is disabled, matrixLocation value will be nil
	if matrixLocation != "" {
		location = matrixLocation
	}

	// Get project details
	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	clusterName := d.EqualsQualString("display_name")

	// Empty check
	if clusterName == "" {
		return nil, nil
	}

	service, err := AlloyDBService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_alloydb_cluster.getAlloydbCluster", "connection_error", err)
		return nil, err
	}

	resp, err := service.Projects.Locations.Clusters.Get("projects/" + project + "/locations/" + location + "/clusters/" + clusterName).Do()
	if err != nil {
		plugin.Logger(ctx).Error("gcp_alloydb_instance.getAlloydbCluster", "api_error", err)
		return nil, err
	}

	return resp, nil
}

func alloyDBClusterSelfLink(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cluster := h.Item.(*alloydb.Cluster)

	selfLink := "https://alloydb.googleapis.com/v1/" + cluster.Name

	return selfLink, nil
}

func gcpAlloyDBClusterTurbotData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cluster := h.Item.(*alloydb.Cluster)

	var location string
	matrixLocation := d.EqualsQualString(matrixKeyLocation)

	// Since, when the service API is disabled, matrixLocation value will be nil
	if matrixLocation != "" {
		location = matrixLocation
	}

	turbotData := map[string]interface{}{
		"Project":  strings.Split(cluster.Name, "/")[1],
		"Location": location,
		"Akas":     []string{"gcp://alloydb.googleapis.com/" + cluster.Name},
	}

	return turbotData, nil
}

//// TRANSFORM FUNCTION

func getAlloyDBClusterDisplayName(ctx context.Context, h *transform.TransformData) (interface{}, error) {
	displayName := ""
	if h.HydrateItem != nil {
		data := h.HydrateItem.(*alloydb.Cluster)
		displayName = strings.Split(data.Name, "/")[5]
	}

	return displayName, nil
}
