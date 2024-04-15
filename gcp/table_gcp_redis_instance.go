package gcp

import (
	"context"
	"strings"

	"cloud.google.com/go/redis/apiv1/redispb"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/iterator"
)

//// TABLE DEFINITION

func tableGcpRedisInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_redis_instance",
		Description: "GCP Redis Instance",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "location"}),
			Hydrate:    getGcpRedisInstance,
			Tags:       map[string]string{"service": "redis", "action": "instances.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listGcpRedisInstances,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "location", Require: plugin.Optional},
			},
			Tags: map[string]string{"service": "redis", "action": "instances.list"},
		},
		GetMatrixItemFunc: BuildRedisLocationList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "Unique name of the resource in this scope including project and location.",
				Transform:   transform.FromField("Name").Transform(lastPathElement),
			},
			{
				Name:        "display_name",
				Description: "An arbitrary and optional user-provided name for the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "labels",
				Description: "Resource labels to represent user provided metadata.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "location_id",
				Description: "The zone where the instance will be provisioned. If not provided, the service will choose a zone from the specified region for the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "alternative_location_id",
				Description: "If specified, at least one node will be provisioned in this zone in addition to the zone specified in location_id. Only applicable to standard tier. If provided, it must be a different zone from the one provided in [location_id]. Additional nodes beyond the first 2 will be placed in zones selected by the service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "redis_version",
				Description: "The version of Redis software.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "reserved_ip_range",
				Description: "The reserved IP range for the instnce.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "secondary_ip_range",
				Description: "Additional IP range for node placement.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "host",
				Description: "Hostname or IP address of the exposed Redis endpoint used by clients to connect to the service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "port",
				Description: "The port number of the exposed Redis endpoint.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "current_location_id",
				Description: "The current zone where the Redis primary node is located. In basic tier, this will always be the same as [location_id]. In standard tier, this can be the zone of any node in the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "The time the instance was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromP(gcpRedisInstanceCreateTime, "CreateTime"),
			},
			{
				Name:        "state",
				Description: "The current state of this instance.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "status_message",
				Description: "Additional information about the current status of this instance, if available.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "redis_configs",
				Description: "Redis configuration parameters.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tier",
				Description: "The service tier of the instance.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "memory_size_gb",
				Description: "Redis memory size in GiB.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "authorized_network",
				Description: "The full name of the Google Compute Engine to which the instance is connected. If left unspecified, the `default` network will be used.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "persistence_iam_identity",
				Description: "Cloud IAM identity used by import / export operations to transfer data to/from Cloud Storage. Format is `serviceAccount:<service_account_email>`.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "connect_mode",
				Description: "The network connect mode of the Redis instance.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "auth_enabled",
				Description: "Indicates whether OSS Redis AUTH is enabled for the instance.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "server_ca_certs",
				Description: "List of server CA certificates for the instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "transit_encryption_mode",
				Description: "The TLS mode of the Redis instance.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "maintenance_policy",
				Description: "The maintenance policy for the instance. If not provided, maintenance events can be performed at any time.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "maintenance_schedule",
				Description: "Date and time of upcoming maintenance events which have been scheduled.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "replica_count",
				Description: "The number of replica nodes.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "nodes",
				Description: "Info per node.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "read_endpoint",
				Description: "Hostname or IP address of the exposed readonly Redis endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "read_endpoint_port",
				Description: "The port number of the exposed readonly redis endpoint.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "read_replicas_mode",
				Description: "Read replicas mode for the instance.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "customer_managed_key",
				Description: "The KMS key reference that the customer provides when trying to create the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "persistence_config",
				Description: "Persistence configuration parameters.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "suspension_reasons",
				Description: "Reasons that causes instance in `SUSPENDED` state.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "maintenance_version",
				Description: "The self service update maintenance version.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "available_maintenance_versions",
				Description: "The available maintenance versions that an instance could update to.",
				Type:        proto.ColumnType_JSON,
			},

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
				Transform:   transform.FromP(gcpRedisInstanceTurbotData, "Akas"),
			},

			// Standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpRedisInstanceTurbotData, "location"),
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

func listGcpRedisInstances(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Create Service Connection
	service, err := RedisService(ctx, d)
	if err != nil {
		logger.Error("gcp_redis_instance.listGcpRedisInstances", "connection_error", err)
		return nil, err
	}

	location := d.EqualsQualString("location")
	matrixLocation := d.EqualsQualString(matrixKeyRedisLocation)
	// Since, when the service API is disabled, matrixLocation value will be nil
	if location != "" && location != matrixLocation {
		return nil, nil
	}

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		logger.Error("gcp_redis_instance.listGcpRedisInstances", "cache_error", err)
		return nil, err
	}
	project := projectId.(string)

	parent := "projects/" + project + "/locations/" + matrixLocation
	req := &redispb.ListInstancesRequest{
		Parent: parent,
	}

	it := service.ListInstances(ctx, req)
	for {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		resp, err := it.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			logger.Error("gcp_redis_instance.listGcpRedisInstances", "api_error", err)
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

func getGcpRedisInstance(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	instanceName := d.EqualsQualString("name")

	// Create Service Connection
	service, err := RedisService(ctx, d)
	if err != nil {
		logger.Error("gcp_redis_instance.getGcpRedisInstance", "connection_error", err)
		return nil, err
	}

	location := d.EqualsQualString("location")
	matrixLocation := d.EqualsQualString(matrixKeyRedisLocation)
	// Since, when the service API is disabled, matrixLocation value will be nil
	if location != "" && location != matrixLocation {
		return nil, nil
	}

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		logger.Error("gcp_redis_instance.getGcpRedisInstance", "cache_error", err)
		return nil, err
	}
	project := projectId.(string)

	name := "projects/" + project + "/locations/" + matrixLocation + "/instances/" + instanceName

	req := &redispb.GetInstanceRequest{
		Name: name,
	}

	op, err := service.GetInstance(ctx, req)
	if err != nil {
		logger.Error("gcp_redis_instance.getGcpRedisInstance", "api_error", err)
		return nil, err
	}

	return op, nil
}

/// TRANSFORM FUNCTIONS

func gcpRedisInstanceTurbotData(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	instance := d.HydrateItem.(*redispb.Instance)
	param := d.Param.(string)
	akas := []string{"gcp://redis.googleapis.com/" + instance.Name}
	locationId := strings.Split(instance.LocationId, "-")
	location := strings.Join(locationId[:len(locationId)-1], "-")
	data := make(map[string]interface{}, 0)
	data["akas"] = akas
	data["location"] = location

	return data[param], nil
}

func gcpRedisInstanceCreateTime(_ context.Context, d *transform.TransformData) (interface{}, error) {
	instanceCreateTime := d.HydrateItem.(*redispb.Instance).CreateTime
	if instanceCreateTime == nil {
		return nil, nil
	}
	return instanceCreateTime.AsTime(), nil
}
