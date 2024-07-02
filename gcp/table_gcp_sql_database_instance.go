package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"google.golang.org/api/googleapi"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

//// TABLE DEFINITION

func tableGcpSQLDatabaseInstance(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_sql_database_instance",
		Description: "GCP SQL Database Instance",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getSQLDatabaseInstance,
			Tags:       map[string]string{"service": "cloudsql", "action": "instances.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listSQLDatabaseInstances,
			KeyColumns: plugin.KeyColumnSlice{
				// String columns
				{Name: "instance_type", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "state", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "database_version", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "backend_type", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "gce_zone", Require: plugin.Optional, Operators: []string{"<>", "="}},
			},
			Tags: map[string]string{"service": "cloudsql", "action": "instances.list"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getSQLDatabaseInstanceUsers,
				Tags: map[string]string{"service": "cloudsql", "action": "users.list"},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "A friendly name that identifies the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "Specifies the current serving state of the Cloud SQL instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_type",
				Description: "Specifies the type of the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "database_version",
				Description: "Specifies the type and version of the database engine.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "machine_type",
				Description: "Specifies the tier or machine type for this instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Settings.Tier"),
			},
			{
				Name:        "data_disk_type",
				Description: "Specifies the type of the data disk used for this instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Settings.DataDiskType"),
			},
			{
				Name:        "data_disk_size_gb",
				Description: "Specifies the size of the data disk, in GB. Minimum size is 10GB. Not used for First Generation instances.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Settings.DataDiskSizeGb"),
			},
			{
				Name:        "storage_auto_resize",
				Description: "Specifies whether the configuration for automatic increment of the the storage size is enabled, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Settings.StorageAutoResize"),
				Default:     true,
			},
			{
				Name:        "enable_point_in_time_recovery",
				Description: "Allows user to recover data from a specific point in time, down to a fraction of a second.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Settings.BackupConfiguration.PointInTimeRecoveryEnabled"),
				Default:     true,
			},
			{
				Name:        "binary_log_enabled",
				Description: "Indicates whether binary log is enabled, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Settings.BackupConfiguration.BinaryLogEnabled"),
			},
			{
				Name:        "kind",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "activation_policy",
				Description: "Describes the activation policy specifies when the instance is activated.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Settings.ActivationPolicy"),
			},
			{
				Name:        "availability_type",
				Description: "Specifies the availability type of the instance. This field is used only for PostgreSQL and MySQL instances.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Settings.AvailabilityType"),
			},
			{
				Name:        "backend_type",
				Description: "Specifies the backend type. Possible values are: FIRST_GEN, SECOND_GEN, EXTERNAL, and SQL_BACKEND_TYPE_UNSPECIFIED.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "backup_enabled",
				Description: "Indicates whether backup configuration is enabled, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Settings.BackupConfiguration.Enabled"),
			},
			{
				Name:        "backup_location",
				Description: "Specifies the backup location.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Settings.BackupConfiguration.Location"),
			},
			{
				Name:        "backup_replication_log_archiving_enabled",
				Description: "Indicates whether backup replication log archiving is enabled, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Settings.BackupConfiguration.ReplicationLogArchivingEnabled"),
			},
			{
				Name:        "backup_start_time",
				Description: "Specifies the start time for the daily backup configuration.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Settings.BackupConfiguration.StartTime"),
			},
			{
				Name:        "connection_name",
				Description: "Specifies the connection name of the Cloud SQL instance used in connection strings.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "crash_safe_replication_enabled",
				Description: "Specifies whether the database flags for crash-safe replication are enabled, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Settings.CrashSafeReplicationEnabled"),
			},
			{
				Name:        "current_disk_size",
				Description: "Specifies the current disk usage of the instance in bytes.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "database_replication_enabled",
				Description: "Specifies whether the replication of database is enabled, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Settings.DatabaseReplicationEnabled"),
			},
			{
				Name:        "gce_zone",
				Description: "Specifies the Compute Engine zone that the instance is currently serving from.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ipv6_address",
				Description: "Specifies the IPv6 address assigned to the instance. This property is applicable only to First Generation instances.",
				Type:        proto.ColumnType_IPADDR,
			},
			{
				Name:        "kms_key_name",
				Description: "Specifies the resource name of KMS key used for disk encryption.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DiskEncryptionConfiguration.KmsKeyName"),
			},
			{
				Name:        "kms_key_version_name",
				Description: "Specifies the KMS key version used to encrypt the Cloud SQL instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DiskEncryptionStatus.KmsKeyVersionName"),
			},
			{
				Name:        "master_instance_name",
				Description: "Specifies the name of the instance which will act as master in the replication setup.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "max_disk_size",
				Description: "Specifies the maximum disk size of the instance in bytes.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "pricing_plan",
				Description: "Specifies the pricing plan for this instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Settings.PricingPlan"),
			},
			{
				Name:        "replication_type",
				Description: "Specifies the type of replication this instance uses.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Settings.ReplicationType"),
			},
			{
				Name:        "self_link",
				Description: "The server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_account_email_address",
				Description: "The service account email address assigned to the instance. This property is applicable only to Second Generation instances.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "settings_version",
				Description: "Specifies the version of instance settings.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Settings.SettingsVersion"),
			},
			{
				Name:        "storage_auto_resize_limit",
				Description: "Specifies the maximum size to which storage capacity can be automatically increased.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Settings.StorageAutoResizeLimit"),
			},
			{
				Name:        "failover_replica_available",
				Description: "The availability status of the failover replica. A false status indicates that the failover replica is out of sync.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("FailoverReplica.Available"),
			},
			{
				Name:        "failover_replica_name",
				Description: "The name of the failover replica. If specified at instance creation, a failover replica is created for the instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FailoverReplica.Name"),
			},
			{
				Name:        "can_defer_maintenance",
				Description: "Indicates whether the scheduled maintenance can be deferred, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ScheduledMaintenance.CanDefer"),
			},
			{
				Name:        "can_reschedule_maintenance",
				Description: "Indicates whether the scheduled maintenance can be rescheduled, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ScheduledMaintenance.CanReschedule"),
			},
			{
				Name:        "maintenance_start_time",
				Description: "The start time of any upcoming scheduled maintenance for this instance.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ScheduledMaintenance.StartTime"),
			},
			{
				Name:        "authorized_gae_applications",
				Description: "A list of App Engine app IDs, that can access this instance.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Settings.AuthorizedGaeApplications"),
			},
			{
				Name:        "database_flags",
				Description: "A list of database flags passed to the instance at startup.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Settings.DatabaseFlags"),
			},
			{
				Name:        "instance_users",
				Description: "A list of users in the specified Cloud SQL instance.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSQLDatabaseInstanceUsers,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "ip_addresses",
				Description: "A list of assigned IP addresses for the instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "ip_configuration",
				Description: "Describes the settings for IP management. It allows to enable or disable the instance IP and manage which external networks can connect to the instance.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Settings.IpConfiguration"),
			},
			{
				Name:        "location_preference",
				Description: "Describes the location preference settings. This allows the instance to be located as near as possible to either an App Engine app or Compute Engine zone for better performance.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Settings.LocationPreference"),
			},
			{
				Name:        "labels",
				Description: "A label is a key-value pair that helps you organize your Google Cloud instances. You can attach a label to each resource, then filter the resources based on their labels.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Settings.UserLabels"),
			},
			{
				Name:        "maintenance_window",
				Description: "Describes the maintenance window for this instance.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Settings.MaintenanceWindow"),
			},
			{
				Name:        "on_premises_configuration",
				Description: "Describes the configurations specific to on-premises instances.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "replica_names",
				Description: "A list of replicas of the instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "replication_configuration",
				Description: "Describes the configurations specific to failover replicas and read replicas.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "ssl_configuration",
				Description: "Describes the SSL configuration of the instance.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ServerCaCert"),
			},
			{
				Name:        "suspension_reason",
				Description: "A list of reasons for the suspension, if the instance state is SUSPENDED.",
				Type:        proto.ColumnType_JSON,
			},

			// standard steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Settings.UserLabels"),
			},
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
				Transform:   transform.From(sqlDatabaseInstanceAka),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Region"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
			},
		},
	}
}

//// LIST FUNCTION

func listSQLDatabaseInstances(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listSQLDatabaseInstances")

	// Create service connection
	service, err := CloudSQLAdminService(ctx, d)
	if err != nil {
		return nil, err
	}

	filterQuals := []filterQualMap{
		{"instance_type", "instanceType", "string"},
		{"state", "state", "string"},
		{"database_version", "databaseVersion", "string"},
		{"backend_type", "backendType", "string"},
		{"gce_zone", "gceZone", "string"},
	}

	filters := buildQueryFilterFromQuals(filterQuals, d.Quals)
	filterString := ""
	if len(filters) > 0 {
		filterString = strings.Join(filters, " ")
	}

	// Max limit isn't mentioned in the documentation
	// Default limit is set as 1000
	pageSize := types.Int64(1000)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *pageSize {
			pageSize = limit
		}
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	resp := service.Instances.List(project).Filter(filterString).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *sqladmin.InstancesListResponse) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, instance := range page.Items {
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
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSQLDatabaseInstance(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSQLDatabaseInstance")

	// Create service connection
	service, err := CloudSQLAdminService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	name := d.EqualsQuals["name"].GetStringValue()

	resp, err := service.Instances.Get(project, name).Do()
	if err != nil {
		return nil, err
	}

	// If the name field kept as empty, API does not returns any error
	if len(resp.Name) < 1 {
		return nil, nil
	}

	return resp, nil
}

func getSQLDatabaseInstanceUsers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create service connection
	service, err := CloudSQLAdminService(ctx, d)
	if err != nil {
		return nil, err
	}

	instance := h.Item.(*sqladmin.DatabaseInstance)
	project := instance.Project

	req, err := service.Users.List(project, instance.Name).Do()
	if err != nil {
		if err.(*googleapi.Error).Message == "Invalid request: Invalid request since instance is not running." {
			return nil, nil
		}
		return nil, err
	}

	return req.Items, nil
}

//// TRANSFORM FUNCTIONS

func sqlDatabaseInstanceAka(_ context.Context, d *transform.TransformData) (interface{}, error) {
	instance := d.HydrateItem.(*sqladmin.DatabaseInstance)

	akas := []string{"gcp://cloudsql.googleapis.com/projects/" + instance.Project + "/regions/" + instance.Region + "/instances/" + instance.Name}

	return akas, nil
}
