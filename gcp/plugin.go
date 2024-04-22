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
			"gcp_apikeys_key":                                         tableGcpApiKeysKey(ctx),
			"gcp_app_engine_application":                              tableGcpAppEngineApplication(ctx),
			"gcp_artifact_registry_repository":                        tableGcpArtifactRegistryRepository(ctx),
			"gcp_audit_policy":                                        tableGcpAuditPolicy(ctx),
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
			"gcp_cloud_run_service":                                   tableGcpCloudRunService(ctx),
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
			"gcp_compute_subnetwork":                                  tableGcpComputeSubnetwork(ctx),
			"gcp_compute_target_https_proxy":                          tableGcpComputeTargetHttpsProxy(ctx),
			"gcp_compute_target_pool":                                 tableGcpComputeTargetPool(ctx),
			"gcp_compute_target_ssl_proxy":                            tableGcpComputeTargetSslProxy(ctx),
			"gcp_compute_target_vpn_gateway":                          tableGcpComputeTargetVpnGateway(ctx),
			"gcp_compute_url_map":                                     tableGcpComputeURLMap(ctx),
			"gcp_compute_vpn_tunnel":                                  tableGcpComputeVpnTunnel(ctx),
			"gcp_compute_zone":                                        tableGcpComputeZone(ctx),
			"gcp_dataproc_cluster":                                    tableGcpDataprocCluster(ctx),
			"gcp_dns_managed_zone":                                    tableGcpDnsManagedZone(ctx),
			"gcp_dns_policy":                                          tableDnsPolicy(ctx),
			"gcp_dns_record_set":                                      tableDnsRecordSet(ctx),
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
			"gcp_project":                                             tableGcpProject(ctx),
			"gcp_project_organization_policy":                         tableGcpProjectOrganizationPolicy(ctx),
			"gcp_project_service":                                     tableGcpProjectService(ctx),
			"gcp_pubsub_snapshot":                                     tableGcpPubSubSnapshot(ctx),
			"gcp_pubsub_subscription":                                 tableGcpPubSubSubscription(ctx),
			"gcp_pubsub_topic":                                        tableGcpPubSubTopic(ctx),
			"gcp_redis_instance":                                      tableGcpRedisInstance(ctx),
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
			"gcp_vertex_ai_endpoint":                                  tableGcpVertexAIEndpoint(ctx),
			"gcp_vertex_ai_model":                                     tableGcpVertexAIModel(ctx),
			/*
				https://github.com/turbot/steampipe/issues/108
				"gcp_compute_route":                   tableGcpComputeRoute(ctx),
			*/

		},
	}

	return p
}
