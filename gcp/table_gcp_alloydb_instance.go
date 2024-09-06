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

func tableGcpAlloyDBInstance(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_alloydb_instance",
		Description: "GCP AlloyDB Instance",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"cluster_name", "instance_display_name"}),
			Hydrate:    getAlloydbInstance,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAlloydbClusters,
			Hydrate:       listAlloydbInstances,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "location", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "cluster_name", Require: plugin.Optional, Operators: []string{"="}},
			},
		},
		GetMatrixItemFunc: BuildAlloyDBLocationList,
		Columns: []*plugin.Column{
			// Changed the column name to instance_display_time because:
			// This table is associated with gcp_alloydb_cluster.
			// There is already a column named display_name in the gcp_alloydb_cluster table.
			// Using the same column name would cause ambiguity when querying this table in conjunction with its parent table using parent qualifiers.
			{
				Name:        "instance_display_name",
				Type:        proto.ColumnType_STRING,
				Description: "User-settable display name for the instance.",
				Transform:   transform.FromField("Name").Transform(getAlloyDBInstanceDisplayName),
			},
			{
				Name:        "cluster_name",
				Type:        proto.ColumnType_STRING,
				Description: "User-settable display name for the cluster.",
				Transform:   transform.FromField("Name").Transform(getAlloyDBClusterDisplayNameTransform),
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The resource name of the instance.",
			},
			{
				Name:        "state",
				Type:        proto.ColumnType_STRING,
				Description: "Current serving state of the instance.",
			},
			{
				Name:        "uid",
				Type:        proto.ColumnType_STRING,
				Description: "System-generated UID of the resource.",
			},
			{
				Name:        "create_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Creation time of the instance.",
				Transform:   transform.FromField("CreateTime").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "availability_type",
				Type:        proto.ColumnType_STRING,
				Description: "Availability type of the instance.",
			},
			{
				Name:        "delete_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Deletion time of the instance.",
				Transform:   transform.FromField("DeleteTime").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "etag",
				Type:        proto.ColumnType_STRING,
				Description: "Resource freshness validation.",
			},
			{
				Name:        "gce_zone",
				Type:        proto.ColumnType_STRING,
				Description: "Compute Engine zone where the instance is located.",
			},
			{
				Name:        "instance_type",
				Type:        proto.ColumnType_STRING,
				Description: "Type of the instance.",
			},
			{
				Name:        "ip_address",
				Type:        proto.ColumnType_IPADDR,
				Description: "IP address assigned to the instance.",
			},
			{
				Name:        "reconciling",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates if the instance is reconciling.",
			},
			{
				Name:        "satisfies_pzs",
				Type:        proto.ColumnType_BOOL,
				Description: "Reserved for future use.",
			},
			{
				Name:        "update_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Last update time of the instance.",
				Transform:   transform.FromField("UpdateTime").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     alloyDBInstanceSelfLink,
				Transform:   transform.FromValue(),
			},
			// JSON Fields
			{
				Name:        "annotations",
				Type:        proto.ColumnType_JSON,
				Description: "Annotations to allow client tools to store a small amount of arbitrary data.",
			},
			{
				Name:        "client_connection_config",
				Type:        proto.ColumnType_JSON,
				Description: "Client connection specific configurations.",
			},
			{
				Name:        "database_flags",
				Type:        proto.ColumnType_JSON,
				Description: "Database flags set at the instance level.",
			},
			{
				Name:        "labels",
				Type:        proto.ColumnType_JSON,
				Description: "Labels as key-value pairs.",
			},
			{
				Name:        "machine_config",
				Type:        proto.ColumnType_JSON,
				Description: "Configurations for the machines that host the underlying database engine.",
			},
			{
				Name:        "nodes",
				Type:        proto.ColumnType_JSON,
				Description: "List of available read-only VMs in this instance.",
			},
			{
				Name:        "query_insights_config",
				Type:        proto.ColumnType_JSON,
				Description: "Configuration for query insights.",
			},
			{
				Name:        "read_pool_config",
				Type:        proto.ColumnType_JSON,
				Description: "Read pool instance configuration.",
			},
			{
				Name:        "writable_node",
				Type:        proto.ColumnType_JSON,
				Description: "This is set for the read-write VM of the PRIMARY instance only.",
			},
			{
				Name:        "connection_info",
				Type:        proto.ColumnType_JSON,
				Description: "ConnectionInfo singleton resource.",
				Hydrate:     getAlloydbInstanceConnectionInfo,
				Transform:   transform.FromValue(),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(getAlloyDBInstanceDisplayName),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     gcpAlloyDBInstanceTurbotData,
				Transform:   transform.FromField("Akas"),
			},

			// Standard GCP columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Hydrate:     gcpAlloyDBInstanceTurbotData,
				Transform:   transform.FromField("Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Hydrate:     gcpAlloyDBInstanceTurbotData,
				Transform:   transform.FromField("Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listAlloydbInstances(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cluster := h.Item.(*alloydb.Cluster)
	clusterName := strings.Split(cluster.Name, "/")[5]

	var location string
	matrixLocation := d.EqualsQualString(matrixKeyLocation)
	// Since, when the service API is disabled, matrixLocation value will be nil
	if matrixLocation != "" {
		location = matrixLocation
	}

	// Create Service Connection
	service, err := AlloyDBService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_alloydb_instance.listAlloydbInstances", "connection_error", err)
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

	parentClusterName := d.EqualsQualString("cluster_name")
	if parentClusterName != "" && parentClusterName != clusterName {
		return nil, nil
	}

	// Get project details
	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	resp := service.Projects.Locations.Clusters.Instances.List("projects/" + project + "/locations/" + location + "/clusters/" + clusterName).PageSize(*pageSize)
	if err := resp.Pages(ctx, func(page *alloydb.ListInstancesResponse) error {
		for _, instance := range page.Instances {
			d.StreamListItem(ctx, instance)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				page.NextPageToken = ""
				return nil
			}
		}
		return nil
	}); err != nil {
		plugin.Logger(ctx).Error("gcp_alloydb_instance.listAlloydbInstances", "api_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTION

func getAlloydbInstance(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

	instanceName := d.EqualsQualString("instance_display_name")
	clusterName := d.EqualsQualString("cluster_name")

	// Empty check
	if instanceName == "" || clusterName == "" {
		return nil, nil
	}

	service, err := AlloyDBService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_alloydb_instance.getAlloydbInstance", "connection_error", err)
		return nil, err
	}

	resp, err := service.Projects.Locations.Clusters.Instances.Get("projects/" + project + "/locations/" + location + "/clusters/" + clusterName + "/instances/" + instanceName).Do()
	if err != nil {
		plugin.Logger(ctx).Error("gcp_alloydb_instance.getAlloydbInstance", "api_error", err)
		return nil, err
	}

	return resp, nil
}

func getAlloydbInstanceConnectionInfo(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	alloyDBInstance := h.Item.(*alloydb.Instance)

	service, err := AlloyDBService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_alloydb_instance.getAlloydbInstanceConnectionInfo", "connection_error", err)
		return nil, err
	}

	resp, err := service.Projects.Locations.Clusters.Instances.GetConnectionInfo(alloyDBInstance.Name).Do()
	if err != nil {
		plugin.Logger(ctx).Error("gcp_alloydb_instance.getAlloydbInstanceConnectionInfo", "api_error", err)
		return nil, err
	}

	return resp, nil
}

func alloyDBInstanceSelfLink(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instance := h.Item.(*alloydb.Instance)

	selfLink := "https://alloydb.googleapis.com/v1/" + instance.Name

	return selfLink, nil
}

func gcpAlloyDBInstanceTurbotData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instance := h.Item.(*alloydb.Instance)

	var location string
	matrixLocation := d.EqualsQualString(matrixKeyLocation)

	// Since, when the service API is disabled, matrixLocation value will be nil
	if matrixLocation != "" {
		location = matrixLocation
	}

	turbotData := map[string]interface{}{
		"Project":  strings.Split(instance.Name, "/")[1],
		"Location": location,
		"Akas":     []string{"gcp://alloydb.googleapis.com/" + instance.Name},
	}

	return turbotData, nil
}

//// TRANSFORM FUNCTION

func getAlloyDBInstanceDisplayName(ctx context.Context, h *transform.TransformData) (interface{}, error) {
	displayName := ""
	if h.HydrateItem != nil {
		data := h.HydrateItem.(*alloydb.Instance)
		displayName = strings.Split(data.Name, "/")[7]
	}

	return displayName, nil
}

func getAlloyDBClusterDisplayNameTransform(ctx context.Context, h *transform.TransformData) (interface{}, error) {
	displayName := ""
	if h.HydrateItem != nil {
		data := h.HydrateItem.(*alloydb.Instance)
		displayName = strings.Split(data.Name, "/")[5]
	}

	return displayName, nil
}
