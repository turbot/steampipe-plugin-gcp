/*
Package gcp implements a steampipe plugin for gcp.

This plugin provides data that Steampipe uses to present foreign
tables that represent GCP resources.
*/
package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

const pluginName = "steampipe-plugin-gcp"

// Plugin creates this (gcp) plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             pluginName,
		DefaultTransform: transform.FromCamel(),
		DefaultGetConfig: &plugin.GetConfig{
			ShouldIgnoreError: isNotFoundError([]string{"404", "400"}),
		},
		TableMap: map[string]*plugin.Table{
			"gcp_audit_policy":                   tableGcpAuditPolicy(ctx),
			"gcp_cloudfunctions_function":        tableGcpCloudfunctionFunction(ctx),
			"gcp_compute_address":                tableGcpComputeAddress(ctx),
			"gcp_compute_firewall":               tableGcpComputeFirewall(ctx),
			"gcp_compute_forwarding_rule":        tableGcpComputeForwardingRule(ctx),
			"gcp_compute_global_address":         tableGcpComputeGlobalAddress(ctx),
			"gcp_compute_global_forwarding_rule": tableGcpComputeGlobalForwardingRule(ctx),
			"gcp_compute_disk":                   tableGcpComputeDisk(ctx), "gcp_compute_instance": tableGcpComputeInstance(ctx),
			"gcp_compute_network":                 tableGcpComputeNetwork(ctx),
			"gcp_compute_router":                  tableGcpComputeRouter(ctx),
			"gcp_compute_snapshot":                tableGcpComputeSnapshot(ctx),
			"gcp_iam_policy":                      tableGcpIAMPolicy(ctx),
			"gcp_iam_role":                        tableGcpIamRole(ctx),
			"gcp_logging_exclusion":               tableGcpLoggingExclusion(ctx),
			"gcp_logging_metric":                  tableGcpLoggingMetric(ctx),
			"gcp_logging_sink":                    tableGcpLoggingSink(ctx),
			"gcp_monitoring_group":                tableGcpMonitoringGroup(ctx),
			"gcp_monitoring_notification_channel": tableGcpMonitoringNotificationChannel(ctx),
			"gcp_project_service":                 tableGcpProjectService(ctx),
			"gcp_pubsub_snapshot":                 tableGcpPubSubSnapshot(ctx),
			"gcp_pubsub_subscription":             tableGcpPubSubSubscription(ctx),
			"gcp_pubsub_topic":                    tableGcpPubSubTopic(ctx),
			"gcp_service_account":                 tableGcpServiceAccount(ctx),
			"gcp_service_account_key":             tableGcpServiceAccountKey(ctx),
			"gcp_storage_bucket":                  tableGcpStorageBucket(ctx),

			/*
				https://github.com/turbot/steampipe/issues/108
				https://github.com/turbot/steampipe/issues/126

				"gcp_compute_image":                   tableGcpComputeImage(ctx),
				"gcp_compute_route":                   tableGcpComputeRoute(ctx),
				"gcp_compute_vpn_tunnel":              tableGcpComputeVpnTunnel(ctx),
			*/

		},
	}

	return p
}
