/*
Package gcp implements a steampipe plugin for gcp.

This plugin provides data that Steampipe uses to present foreign
tables that represent GCP resources.
*/
package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/rate_limiter"
)

const pluginName = "steampipe-plugin-gcp"

// Plugin creates this (gcp) plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             pluginName,
		DefaultTransform: transform.FromCamel(),
		DefaultGetConfig: &plugin.GetConfig{
			ShouldIgnoreError: isIgnorableError([]string{"404", "400"}),
		},
		// Default ignore config for the plugin
		DefaultIgnoreConfig: &plugin.IgnoreConfig{
			ShouldIgnoreErrorFunc: shouldIgnoreErrorPluginDefault(),
		},
		RateLimiters: []*rate_limiter.Definition{
			// API Requests per 100 seconds: 5,000
			// https://cloud.google.com/memorystore/docs/redis/quotas#per-second_api_requests_quota
			{
				Name:       "gcp_redis_list_instances",
				FillRate:   50,
				BucketSize: 5000,
				Scope:      []string{"connection", "service", "action"},
				Where:      "service = 'redis' and action = 'ListInstances'",
			},
			{
				Name:       "gcp_redis_get_instance",
				FillRate:   50,
				BucketSize: 5000,
				Scope:      []string{"connection", "service", "action"},
				Where:      "service = 'redis' and action = 'GetInstance'",
			},
			// Redis Cluster requests per project per minute: 60
			// https://cloud.google.com/memorystore/docs/cluster/quotas#per-minute_api_requests_quota
			{
				Name:       "gcp_rediscluster_list_clusters",
				FillRate:   1,
				BucketSize: 60,
				Scope:      []string{"connection", "service", "action"},
				Where:      "service = 'rediscluster' and action = 'ListClusters'",
			},
			{
				Name:       "gcp_rediscluster_get_cluster",
				FillRate:   1,
				BucketSize: 60,
				Scope:      []string{"connection", "service", "action"},
				Where:      "service = 'rediscluster' and action = 'GetCluster'",
			},

			// Cloud Resource Manager & Service Usage API rate quota: 1,200 requests/minute per user
			// Doc: https://cloud.google.com/resource-manager/quotas (see API rate quotas) and https://cloud.google.com/service-usage/quotas
			// Tables: gcp_project, gcp_organization, gcp_organization_project, gcp_project_organization_policy, gcp_project_service, gcp_iam_policy
			{
				Name:       "gcp_resourcemanager",
				FillRate:   20,
				BucketSize: 1200,
				Scope:      []string{"connection", "service", "action"},
				Where:      "service in ('resourcemanager', 'serviceusage') and action in ('organizations.get', 'projects.list', 'projects.getIamPolicy', 'services.list', 'services.get')",
			},

			// Cloud Resource Manager API rate quota: 600 read requests per minute per project (10 per second)
			// FIXME: Limits are per API consumer project so we need to find a way to take quota_project into account instead of connection
			// Doc: https://cloud.google.com/resource-manager/docs/limits
			// APIs: projects.getAccessApprovalSettings, projects.getAncestry
			// Tables: gcp_project, gcp_organization_project
			// FillRate and BucketSize set to 10 to match the per-second rate and allow for short bursts within the minute quota.
			{
				Name:       "gcp_cloudresourcemanager_projects_get_access_approval_settings",
				FillRate:   10,
				BucketSize: 60,
				Scope:      []string{"connection", "service", "action"},
				Where:      "service = 'resourcemanager' and action = 'projects.getAccessApprovalSettings'",
			},
			{
				Name:       "gcp_cloudresourcemanager_projects_get_ancestry",
				FillRate:   10,
				BucketSize: 60,
				Scope:      []string{"connection", "service", "action"},
				Where:      "service = 'resourcemanager' and action = 'projects.getAncestry'",
			},

			// Compute Engine API rate quotas are enforced per minute and vary by project and method group.
			// This limiter targets a conservative 10 rps with a burst of 20 for common read/list methods.
			// Doc: https://cloud.google.com/compute/api-quota and metrics: https://cloud.google.com/compute/docs/api/compute-api-quota-metrics
			// Tables: gcp_compute_instance, gcp_compute_disk, gcp_compute_image, gcp_compute_snapshot, gcp_compute_address, gcp_compute_global_address, gcp_compute_network, gcp_compute_subnetwork, gcp_compute_firewall, gcp_compute_route, gcp_compute_ssl_policy, gcp_compute_url_map, gcp_compute_forwarding_rule, gcp_compute_global_forwarding_rule
			{
				Name:       "gcp_compute",
				FillRate:   10,
				BucketSize: 20,
				Scope:      []string{"connection", "service", "action"},
				Where:      "service = 'compute' and action in ('instances.get', 'instances.list', 'instances.getIamPolicy', 'disks.get', 'disks.list', 'disks.getIamPolicy', 'images.list', 'images.get', 'images.getIamPolicy', 'snapshots.list', 'snapshots.get', 'addresses.list', 'addresses.get', 'globalAddresses.list', 'globalAddresses.get', 'networks.list', 'networks.get', 'subnetworks.list', 'subnetworks.get', 'subnetworks.getIamPolicy', 'firewalls.list', 'firewalls.get', 'sslPolicies.list', 'sslPolicies.get', 'urlMaps.list', 'urlMaps.get', 'forwardingRules.list', 'forwardingRules.get', 'globalForwardingRules.list', 'globalForwardingRules.get')",
			},

			// Approximately 5000 read requests per second: https://cloud.google.com/storage/docs/request-rate#auto-scaling
			// Tables: gcp_storage_bucket, gcp_storage_object
			{
				Name:       "gcp_storage",
				FillRate:   5000,
				BucketSize: 5000,
				Scope:      []string{"connection", "service", "action"},
				Where:      "service = 'storage' and action in ('buckets.list', 'buckets.get', 'buckets.getIamPolicy', 'objects.list', 'objects.get', 'objects.getIamPolicy')",
			},

			// Secret Manager requests per minute per project: 600
			// Doc: https://cloud.google.com/secret-manager/quotas#request-rate-quotas
			// Tables: gcp_secret_manager_secret
			{
				Name:       "gcp_secret_manager_secret",
				FillRate:   10,
				BucketSize: 600,
				Scope:      []string{"connection", "service", "action"},
				Where:      "service = 'secretmanager' and action in ('secrets.list', 'secrets.get')",
			},

			// Cloud Logging read requests per minute per project: 60
			// Doc: https://cloud.google.com/logging/quotas#api-limits
			// To view your quotas: https://console.cloud.google.com/apis/dashboard (select API → Quotas)
			// Table: gcp_logging_log_entry, gcp_logging_bucket, gcp_logging_metric, gcp_logging_exclusion, gcp_logging_sink
			{
				Name:           "gcp_logging",
				FillRate:       1,
				BucketSize:     60,
				MaxConcurrency: 1,
				Scope:          []string{"connection", "service", "action"},
				Where:          "service = 'logging' and action in ('logEntries.list', 'logEntries.get', 'buckets.list', 'buckets.get', 'logMetrics.list', 'logMetrics.get', 'exclusions.list', 'exclusions.get', 'sinks.list', 'sinks.get')",
			},
			// Pub/Sub Admin API requests per minute per project: 6,000 (≈100 ops/s)
			// Doc: https://cloud.google.com/pubsub/quotas
			// Table: gcp_pubsub_topic, gcp_pubsub_subscription
			{
				Name:       "gcp_pubsub",
				FillRate:   100,
				BucketSize: 6000,
				Scope:      []string{"connection", "service", "action"},
				Where:      "service = 'pubsub' and action in ('topics.list', 'topics.get', 'topics.getIamPolicy', 'subscriptions.list', 'subscriptions.get', 'subscriptions.getIamPolicy', 'snapshots.list', 'snapshots.get', 'snapshots.getIamPolicy')",
			},

			// IAM v1 API read requests per project per minute: 6,000 (e.g., getting an allow policy)
			// Doc: https://cloud.google.com/iam/quotas#quotas
			// Table: gcp_iam_role, gcp_service_account
			{
				Name:       "gcp_iam",
				FillRate:   100,
				BucketSize: 6000,
				Scope:      []string{"connection", "service", "action"},
				Where:      "service = 'iam' and action in ('roles.list', 'roles.get', 'serviceAccounts.list', 'serviceAccounts.get', 'serviceAccounts.getIamPolicy', 'serviceAccountKeys.list', 'serviceAccountKeys.get')",
			},

			// Cloud Functions API read (list/get) requests per minute per user: 60
			// Doc: https://cloud.google.com/functions/quotas
			{
				Name:       "gcp_cloudfunctions_function",
				FillRate:   1,
				BucketSize: 60,
				Scope:      []string{"connection", "service", "action"},
				Where:      "service = 'cloudfunctions' and action in ('functions.list', 'functions.get')",
			},

			// Cloud DNS API calls per 100 seconds per project per region: 60,000
			// Doc: https://cloud.google.com/dns/quotas#quotas
			// Login to the console to see the quota: https://console.cloud.google.com/apis/api/dns.googleapis.com/quotas?inv=1&invt=Ab4ulA
			// Tables: gcp_dns_record_set, gcp_dns_policy, gcp_dns_managed_zone
			{
				Name:       "gcp_dns",
				FillRate:   1000,
				BucketSize: 60000,
				Scope:      []string{"connection", "service", "action", "location"},
				Where:      "service = 'dns' and action in ('managedZones.list', 'managedZones.get', 'policies.list', 'policies.get', 'resourceRecordSets.list', 'resourceRecordSets.get')",
			},
		},
		ConnectionKeyColumns: []plugin.ConnectionKeyColumn{
			{
				Name:    "project",
				Hydrate: getProject,
			},
		},
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
		},
		TableMap: map[string]*plugin.Table{
			"gcp_alloydb_cluster":                                     tableGcpAlloyDBCluster(ctx),
			"gcp_alloydb_instance":                                    tableGcpAlloyDBInstance(ctx),
			"gcp_apikeys_key":                                         tableGcpApiKeysKey(ctx),
			"gcp_app_engine_application":                              tableGcpAppEngineApplication(ctx),
			"gcp_artifact_registry_repository":                        tableGcpArtifactRegistryRepository(ctx),
			"gcp_audit_policy":                                        tableGcpAuditPolicy(ctx),
			"gcp_organization_audit_policy":                           tableGcpOrganizationAuditPolicy(ctx),
			"gcp_bigquery_dataset":                                    tableGcpBigQueryDataset(ctx),
			"gcp_bigquery_job":                                        tableGcpBigQueryJob(ctx),
			"gcp_bigquery_table":                                      tableGcpBigqueryTable(ctx),
			"gcp_bigtable_instance":                                   tableGcpBigtableInstance(ctx),
			"gcp_billing_account":                                     tableGcpBillingAccount(ctx),
			"gcp_billing_budget":                                      tableGcpBillingBudget(ctx),
			"gcp_cloud_asset":                                         tableGcpCloudAsset(ctx),
			"gcp_cloud_identity_group":                                tableGcpCloudIdentityGroup(ctx),
			"gcp_cloud_identity_group_membership":                     tableGcpCloudIdentityGroupMembership(ctx),
			"gcp_cloudfunctions_function":                             tableGcpCloudfunctionFunction(ctx),
			"gcp_cloud_run_job":                                       tableGcpCloudRunJob(ctx),
			"gcp_cloud_run_service":                                   tableGcpCloudRunService(ctx),
			"gcp_composer_environment":                                tableGcpComposerEnvironment(ctx),
			"gcp_compute_address":                                     tableGcpComputeAddress(ctx),
			"gcp_compute_autoscaler":                                  tableGcpComputeAutoscaler(ctx),
			"gcp_compute_backend_bucket":                              tableGcpComputeBackendBucket(ctx),
			"gcp_compute_backend_service":                             tableGcpComputeBackendService(ctx),
			"gcp_compute_disk":                                        tableGcpComputeDisk(ctx),
			"gcp_compute_disk_metric_read_ops":                        tableGcpComputeDiskMetricReadOps(ctx),
			"gcp_compute_disk_metric_read_ops_daily":                  tableGcpComputeDiskMetricReadOpsDaily(ctx),
			"gcp_compute_disk_metric_read_ops_hourly":                 tableGcpComputeDiskMetricReadOpsHourly(ctx),
			"gcp_compute_disk_metric_write_ops":                       tableGcpComputeDiskMetricWriteOps(ctx),
			"gcp_compute_disk_metric_write_ops_daily":                 tableGcpComputeDiskMetricWriteOpsDaily(ctx),
			"gcp_compute_disk_metric_write_ops_hourly":                tableGcpComputeDiskMetricWriteOpsHourly(ctx),
			"gcp_compute_firewall":                                    tableGcpComputeFirewall(ctx),
			"gcp_compute_forwarding_rule":                             tableGcpComputeForwardingRule(ctx),
			"gcp_compute_global_address":                              tableGcpComputeGlobalAddress(ctx),
			"gcp_compute_global_forwarding_rule":                      tableGcpComputeGlobalForwardingRule(ctx),
			"gcp_compute_ha_vpn_gateway":                              tableGcpComputeHaVpnGateway(ctx),
			"gcp_compute_image":                                       tableGcpComputeImage(ctx),
			"gcp_compute_instance":                                    tableGcpComputeInstance(ctx),
			"gcp_compute_instance_group":                              tableGcpComputeInstanceGroup(ctx),
			"gcp_compute_instance_group_manager":                      tableGcpComputeInstanceGroupManager(ctx),
			"gcp_compute_instance_metric_cpu_utilization":             tableGcpComputeInstanceMetricCpuUtilization(ctx),
			"gcp_compute_instance_metric_cpu_utilization_daily":       tableGcpComputeInstanceMetricCpuUtilizationDaily(ctx),
			"gcp_compute_instance_metric_cpu_utilization_hourly":      tableGcpComputeInstanceMetricCpuUtilizationHourly(ctx),
			"gcp_compute_instance_template":                           tableGcpComputeInstanceTemplate(ctx),
			"gcp_compute_machine_image":                               tableGcpComputeMachineImage(ctx),
			"gcp_compute_machine_type":                                tableGcpComputeMachineType(ctx),
			"gcp_compute_network":                                     tableGcpComputeNetwork(ctx),
			"gcp_compute_node_group":                                  tableGcpComputeNodeGroup(ctx),
			"gcp_compute_node_template":                               tableGcpComputeNodeTemplate(ctx),
			"gcp_compute_project_metadata":                            tableGcpComputeProjectMetadata(ctx),
			"gcp_compute_region":                                      tableGcpComputeRegion(ctx),
			"gcp_compute_resource_policy":                             tableGcpComputeResourcePolicy(ctx),
			"gcp_compute_router":                                      tableGcpComputeRouter(ctx),
			"gcp_compute_snapshot":                                    tableGcpComputeSnapshot(ctx),
			"gcp_compute_ssl_policy":                                  tableGcpComputeSslPolicy(ctx),
			"gcp_compute_security_policy":                             tableGcpComputeSecurityPolicy(ctx),
			"gcp_compute_subnetwork":                                  tableGcpComputeSubnetwork(ctx),
			"gcp_compute_target_https_proxy":                          tableGcpComputeTargetHttpsProxy(ctx),
			"gcp_compute_target_pool":                                 tableGcpComputeTargetPool(ctx),
			"gcp_compute_target_ssl_proxy":                            tableGcpComputeTargetSslProxy(ctx),
			"gcp_compute_target_vpn_gateway":                          tableGcpComputeTargetVpnGateway(ctx),
			"gcp_compute_tpu":                                         tableGcpComputeTpu(ctx),
			"gcp_compute_url_map":                                     tableGcpComputeURLMap(ctx),
			"gcp_compute_vpn_tunnel":                                  tableGcpComputeVpnTunnel(ctx),
			"gcp_compute_zone":                                        tableGcpComputeZone(ctx),
			"gcp_dataplex_asset":                                      tableGcpDataplexAsset(ctx),
			"gcp_dataplex_lake":                                       tableGcpDataplexLake(ctx),
			"gcp_dataplex_task":                                       tableGcpDataplexTask(ctx),
			"gcp_dataplex_zone":                                       tableGcpDataplexZone(ctx),
			"gcp_dataproc_cluster":                                    tableGcpDataprocCluster(ctx),
			"gcp_dataproc_metastore_service":                          tableGcpDataprocMetastoreService(ctx),
			"gcp_dns_managed_zone":                                    tableGcpDnsManagedZone(ctx),
			"gcp_dns_policy":                                          tableDnsPolicy(ctx),
			"gcp_dns_record_set":                                      tableDnsRecordSet(ctx),
			"gcp_firestore_database":                                  tableGcpFirestoreDatabase(ctx),
			"gcp_iam_policy":                                          tableGcpIAMPolicy(ctx),
			"gcp_iam_role":                                            tableGcpIamRole(ctx),
			"gcp_kms_key":                                             tableGcpKmsKey(ctx),
			"gcp_kms_key_ring":                                        tableGcpKmsKeyRing(ctx),
			"gcp_kms_key_version":                                     tableGcpKmsKeyVersion(ctx),
			"gcp_kubernetes_cluster":                                  tableGcpKubernetesCluster(ctx),
			"gcp_kubernetes_node_pool":                                tableGcpKubernetesNodePool(ctx),
			"gcp_logging_bucket":                                      tableGcpLoggingBucket(ctx),
			"gcp_logging_exclusion":                                   tableGcpLoggingExclusion(ctx),
			"gcp_logging_log_entry":                                   tableGcpLoggingLogEntry(ctx),
			"gcp_logging_metric":                                      tableGcpLoggingMetric(ctx),
			"gcp_logging_sink":                                        tableGcpLoggingSink(ctx),
			"gcp_monitoring_alert_policy":                             tableGcpMonitoringAlert(ctx),
			"gcp_monitoring_group":                                    tableGcpMonitoringGroup(ctx),
			"gcp_monitoring_notification_channel":                     tableGcpMonitoringNotificationChannel(ctx),
			"gcp_organization":                                        tableGcpOrganization(ctx),
			"gcp_organization_project":                                tableGcpOrganizationProject(ctx),
			"gcp_project":                                             tableGcpProject(ctx),
			"gcp_project_organization_policy":                         tableGcpProjectOrganizationPolicy(ctx),
			"gcp_project_service":                                     tableGcpProjectService(ctx),
			"gcp_pubsub_snapshot":                                     tableGcpPubSubSnapshot(ctx),
			"gcp_pubsub_subscription":                                 tableGcpPubSubSubscription(ctx),
			"gcp_pubsub_topic":                                        tableGcpPubSubTopic(ctx),
			"gcp_redis_cluster":                                       tableGcpRedisCluster(ctx),
			"gcp_redis_instance":                                      tableGcpRedisInstance(ctx),
			"gcp_secret_manager_secret":                               tableGcpSecretManagerSecret(ctx),
			"gcp_service_account":                                     tableGcpServiceAccount(ctx),
			"gcp_service_account_key":                                 tableGcpServiceAccountKey(ctx),
			"gcp_sql_backup":                                          tableGcpSQLBackup(ctx),
			"gcp_sql_database":                                        tableGcpSQLDatabase(ctx),
			"gcp_sql_database_instance":                               tableGcpSQLDatabaseInstance(ctx),
			"gcp_sql_database_instance_metric_connections":            tableGcpSQLDatabaseInstanceMetricConnections(ctx),
			"gcp_sql_database_instance_metric_connections_daily":      tableGcpSQLDatabaseInstanceMetricConnectionsDaily(ctx),
			"gcp_sql_database_instance_metric_connections_hourly":     tableGcpSQLDatabaseInstanceMetricConnectionsHourly(ctx),
			"gcp_sql_database_instance_metric_cpu_utilization":        tableGcpSQLDatabaseInstanceMetricCpuUtilization(ctx),
			"gcp_sql_database_instance_metric_cpu_utilization_daily":  tableGcpSQLDatabaseInstanceMetricCpuUtilizationDaily(ctx),
			"gcp_sql_database_instance_metric_cpu_utilization_hourly": tableGcpSQLDatabaseInstanceMetricCpuUtilizationHourly(ctx),
			"gcp_storage_bucket":                                      tableGcpStorageBucket(ctx),
			"gcp_storage_object":                                      tableGcpStorageObject(ctx),
			"gcp_tag_binding":                                         tableGcpTagBinding(ctx),
			"gcp_tpu_vm":                                              tableGcpTpuVM(ctx),
			"gcp_vertex_ai_endpoint":                                  tableGcpVertexAIEndpoint(ctx),
			"gcp_vertex_ai_notebook_runtime_template":                 tableGcpVertexAINotebookRuntimeTemplate(ctx),
			"gcp_vertex_ai_model":                                     tableGcpVertexAIModel(ctx),
			"gcp_vpc_access_connector":                                tableGcpVPCAccessConnector(ctx),
			"gcp_workstations_workstation_cluster":                    tableGcpWorkstationsWorkstationCluster(ctx),
			"gcp_workstations_workstation":                            tableGcpWorkstationsWorkstation(ctx),
			/*
				https://github.com/turbot/steampipe/issues/108
				"gcp_compute_route":                   tableGcpComputeRoute(ctx),
			*/

		},
	}

	return p
}
