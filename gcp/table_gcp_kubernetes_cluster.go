package gcp

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	"google.golang.org/api/container/v1"
)

//// TABLE DEFINITION

func tableGcpKubernetesCluster(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_kubernetes_cluster",
		Description: "GCP Kubernetes Cluster",
		List: &plugin.ListConfig{
			Hydrate:           listKubernetesClusters,
			ShouldIgnoreError: isIgnorableError([]string{"403"}),
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "location"}),
			Hydrate:    getKubernetesCluster,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of this cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "location_type",
				Description: "Location type of the cluster i.e REGIONAL/ZONAL.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(gcpKubernetesClusterLocationType),
			},
			{
				Name:        "status",
				Description: "The current status of this cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "An optional description of this cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "The time the cluster was created, in RFC3339 text format.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "cluster_ipv4_cidr",
				Description: "The IP address range of the container pods in this cluster, in CIDR notation.",
				Type:        proto.ColumnType_CIDR,
			},
			{
				Name:        "shielded_nodes_enabled",
				Description: "Denotes whether Shielded Nodes features are enabled on all nodes in this cluster.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ShieldedNodes.Enabled"),
			},
			{
				Name:        "database_encryption_key_name",
				Description: "Name of CloudKMS key to use for the encryption.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseEncryption.KeyName"),
			},
			{
				Name:        "database_encryption_state",
				Description: "Denotes the state of etcd encryption.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DatabaseEncryption.State"),
			},
			{
				Name:        "max_pods_per_node",
				Description: "Constraint enforced on the max num of pods per node.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("DefaultMaxPodsConstraint.MaxPodsPerNode"),
			},
			{
				Name:        "legacy_abac_enabled",
				Description: "Configuration for the legacy ABAC authorization mode.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("LegacyAbac.Enabled"),
			},
			{
				Name:        "current_master_version",
				Description: "The current software version of the master endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "current_node_count",
				Description: "The number of nodes currently in the cluster.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "current_node_version",
				Description: "The current version of the node software components.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status_message",
				Description: "Additional information about the current status of this cluster, if available.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "subnetwork",
				Description: "The name of the Google Compute Engine subnetwork to which the cluster is connected.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "enable_kubernetes_alpha",
				Description: "Indicates whether kubernetes alpha features are enabled on this cluster.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "enable_tpu",
				Description: "Enable the ability to use Cloud TPUs in this cluster.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "endpoint",
				Description: "The IP address of this cluster's master endpoint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "expire_time",
				Description: "The time the cluster will be automatically deleted.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ExpireTime").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "initial_cluster_version",
				Description: "The initial Kubernetes version for this cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "initial_node_count",
				Description: "The number of nodes to create in this cluster.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "label_fingerprint",
				Description: "The fingerprint of the set of labels for this cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "logging_service",
				Description: "The logging service the cluster should use to write logs.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "monitoring_service",
				Description: "The monitoring service the cluster should use to write metrics.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "network",
				Description: "The name of the Google Compute Engine network to which the cluster is connected.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "node_ipv4_cidr_size",
				Description: "The size of the address space on each node for hosting containers.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "services_ipv4_cidr",
				Description: "The IP address range of the Kubernetes services in this cluster, in CIDR notation.",
				Type:        proto.ColumnType_CIDR,
			},
			{
				Name:        "tpu_ipv4_cidr_block",
				Description: "The IP address range of the Cloud TPUs in this cluster, in CIDR notation.",
				Type:        proto.ColumnType_CIDR,
			},
			{
				Name:        "zone",
				Description: "The name of the Google Compute Engine zone in which the cluster resides.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "addons_config",
				Description: "Configurations for the various addons available to run in the cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "authenticator_groups_config",
				Description: "Configuration controlling RBAC group membership information.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "autoscaling",
				Description: "Cluster-level autoscaling configuration.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "conditions",
				Description: "Which conditions caused the current cluster state.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "binary_authorization",
				Description: "Configuration for Binary Authorization.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "instance_group_urls",
				Description: "List of urls for instance groups.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "ip_allocation_policy",
				Description: "Configuration for cluster IP allocation.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "locations",
				Description: "The list of Google Compute Engine zones in which the cluster's nodes should be located.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "maintenance_policy",
				Description: "Configure the maintenance policy for this cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "master_auth",
				Description: "The authentication information for accessing the master endpoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "master_authorized_networks_config",
				Description: "The configuration options for master authorized networks feature.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "network_config",
				Description: "Configuration for cluster networking.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "network_policy",
				Description: "Configuration options for the NetworkPolicy feature.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "node_config",
				Description: "Parameters used in creating the cluster's nodes.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "node_pools",
				Description: "The node pools associated with this cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "notification_config",
				Description: "Notification configuration of the cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "private_cluster_config",
				Description: "Configuration for private cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "release_channel",
				Description: "Release channel configuration.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "resource_labels",
				Description: "The resource labels for the cluster to use to annotate any related Google Compute Engine resources.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "resource_usage_export_config",
				Description: "Configuration for exporting resource usages.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "vertical_pod_autoscaling",
				Description: "Cluster-level Vertical Pod Autoscaling configuration.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "workload_identity_config",
				Description: "Configuration for the use of Kubernetes Service Accounts in GCP IAM policies.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ResourceLabels"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(gcpKubernetesClusterTurbotData, "Akas"),
			},

			// GCP standard columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpKubernetesClusterTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listKubernetesClusters(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listKubernetesClusters")

	// Create Service Connection
	service, err := ContainerService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details
	getProjectCached := plugin.HydrateFunc(getProject).WithCache()
	projectId, err := getProjectCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	resp, err := service.Projects.Locations.Clusters.List("projects/" + project + "/locations/-").Do()
	if err != nil {
		return nil, err
	}

	for _, cluster := range resp.Clusters {
		d.StreamListItem(ctx, cluster)
	}

	return nil, nil
}

/// HYDRATE FUNCTIONS

func getKubernetesCluster(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getKubernetesCluster")

	// Create Service Connection
	service, err := ContainerService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details
	getProjectCached := plugin.HydrateFunc(getProject).WithCache()
	projectId, err := getProjectCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	name := d.KeyColumnQuals["name"].GetStringValue()
	location := d.KeyColumnQuals["location"].GetStringValue()

	// Return nil, if no input provided
	if name == "" || location == "" {
		return nil, nil
	}

	resp, err := service.Projects.Locations.Clusters.Get("projects/" + project + "/locations/" + location + "/clusters/" + name).Do()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

//// TRANSFORM FUNCTIONS

func gcpKubernetesClusterTurbotData(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("gcpKubernetesClusterTurbotData")
	cluster := d.HydrateItem.(*container.Cluster)

	splitName := strings.Split(cluster.SelfLink, "/")
	akas := []string{strings.Replace(cluster.SelfLink, "https://", "gcp://", 1)}

	result := map[string]interface{}{
		"ClusterName": splitName[9],
		"Location":    splitName[7],
		"Project":     splitName[5],
		"Akas":        akas,
	}
	return result[d.Param.(string)], nil
}

func gcpKubernetesClusterLocationType(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("gcpKubernetesClusterLocationType")
	cluster := d.HydrateItem.(*container.Cluster)

	splitName := strings.Split(cluster.SelfLink, "/")

	if splitName[6] == "locations" {
		return "REGIONAL", nil
	}
	return "ZONAL", nil
}
