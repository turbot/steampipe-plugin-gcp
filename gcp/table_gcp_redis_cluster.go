package gcp

import (
	"context"
	"strings"

	"cloud.google.com/go/redis/cluster/apiv1/clusterpb"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/iterator"
)

//// TABLE DEFINITION

func tableGcpRedisCluster(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_redis_cluster",
		Description: "GCP Redis Cluster",
		Get: &plugin.GetConfig{
			Hydrate:    getGcpRedisCluster,
			KeyColumns: plugin.AllColumns([]string{"name", "location"}),
			Tags:       map[string]string{"service": "rediscluster", "action": "ListClusters"},
		},
		List: &plugin.ListConfig{
			Hydrate: listGcpRedisClusters,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "location", Require: plugin.Optional},
			},
			Tags: map[string]string{"service": "rediscluster", "action": "GetCluster"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "Unique name of the resource in this scope including project and location.",
				Transform:   transform.FromField("Name").Transform(lastPathElement),
			},
			{
				Name:        "create_time",
				Description: "The timestamp associated with the cluster creation request.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromP(gcpRedisClusterCreateTime, "CreateTime"),
			},
			{
				Name:        "state",
				Description: "The current state of this cluster.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "uid",
				Description: "System assigned, unique identifier for the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "authorization_mode",
				Description: "The authorization mode of the Redis cluster.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "transit_encryption_mode",
				Description: "The in-transit encryption for the Redis cluster.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "psc_configs",
				Description: "Each PscConfig configures the consumer network where IPs will be designated to the cluster for client access through Private Service Connect Automation.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "discovery_endpoints",
				Description: "Endpoints created on each given network, for Redis clients to connect to the cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "psc_connections",
				Description: "PSC connections for discovery of the cluster topology and accessing the cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "state_info",
				Description: "Additional information about the current state of the cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "node_type",
				Description: "The type of Redis nodes in the cluster that determines the underlying machine-type.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "persistence_config",
				Description: "Persistence config (RDB, AOF) for the cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "redis_configs",
				Description: "Key/Value pairs of customer overrides for mutable Redis Configs.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "zone_distribution_config",
				Description: "This config will be used to determine how the customer wants us to distribute cluster resources within the region.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "cross_cluster_replication_config",
				Description: "Cross cluster replication config.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "replica_count",
				Description: "The number of replica nodes per shard.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "size_gb",
				Description: "Redis memory size in GB for the entire cluster rounded up to the next integer.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "shard_count",
				Description: "Number of shards for the Redis cluster.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "precise_size_gb",
				Description: "Precise value of redis memory size in GB for the entire cluster.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "deletion_protection_enabled",
				Description: "The delete operation will fail when the value is set to true.",
				Type:        proto.ColumnType_JSON,
			},
			// FIXME: this is missing from the Go SDK
			// https://github.com/googleapis/google-cloud-go/issues/11061
			// {
			// 	Name:        "maintenance_policy",
			// 	Description: "The maintenance policy for the cluster. If not provided, maintenance events can be performed at any time.",
			// 	Type:        proto.ColumnType_JSON,
			// },
			// {
			// 	Name:        "maintenance_schedule",
			// 	Description: "Date and time of upcoming maintenance events which have been scheduled.",
			// 	Type:        proto.ColumnType_JSON,
			// },

			// Standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(gcpRedisClusterTurbotData, "Akas"),
			},

			// Standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpRedisClusterTurbotData, "location"),
			},
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

//// LIST FUNCTIONS

func listGcpRedisClusters(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Create Service Connection
	service, err := RedisClusterService(ctx, d)
	if err != nil {
		logger.Error("gcp_redis_cluster.listGcpRedisClusters", "connection_error", err)
		return nil, err
	}

	location := d.EqualsQualString("location")
	if location == "" {
		// Wildcard to query all locations at once
		// https://cloud.google.com/memorystore/docs/cluster/reference/rest/v1/projects.locations.clusters/list
		location = "-"
	}

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		logger.Error("gcp_redis_cluster.listGcpRedisClusters", "cache_error", err)
		return nil, err
	}
	project := projectId.(string)

	parent := "projects/" + project + "/locations/" + location
	req := &clusterpb.ListClustersRequest{
		Parent: parent,
	}

	it := service.ListClusters(ctx, req)
	for {
		resp, err := it.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			logger.Error("gcp_redis_cluster.listGcpRedisClusters", "api_error", err)
			return nil, err
		}

		d.StreamListItem(ctx, resp)

		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getGcpRedisCluster(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	clusterName := d.EqualsQualString("name")

	// Create Service Connection
	service, err := RedisClusterService(ctx, d)
	if err != nil {
		logger.Error("gcp_redis_cluster.getGcpRedisCluster", "connection_error", err)
		return nil, err
	}

	location := d.EqualsQualString("location")

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		logger.Error("gcp_redis_cluster.getGcpRedisCluster", "cache_error", err)
		return nil, err
	}
	project := projectId.(string)

	name := "projects/" + project + "/locations/" + location + "/clusters/" + clusterName

	req := &clusterpb.GetClusterRequest{
		Name: name,
	}

	op, err := service.GetCluster(ctx, req)
	if err != nil {
		logger.Error("gcp_redis_cluster.getGcpRedisCluster", "api_error", err)
		return nil, err
	}

	return op, nil
}

/// TRANSFORM FUNCTIONS

func gcpRedisClusterTurbotData(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	cluster := d.HydrateItem.(*clusterpb.Cluster)
	param := d.Param.(string)
	akas := []string{"gcp://rediscluster.googleapis.com/" + cluster.Name}
	location := strings.Split(cluster.Name, "/")[3]
	data := make(map[string]interface{}, 0)
	data["akas"] = akas
	data["location"] = location

	return data[param], nil
}

func gcpRedisClusterCreateTime(_ context.Context, d *transform.TransformData) (interface{}, error) {
	clusterCreateTime := d.HydrateItem.(*clusterpb.Cluster).CreateTime
	if clusterCreateTime == nil {
		return nil, nil
	}
	return clusterCreateTime.AsTime(), nil
}
