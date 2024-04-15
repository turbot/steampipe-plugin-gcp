package gcp

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/container/v1"
)

//// TABLE DEFINITION

func tableGcpKubernetesNodePool(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_kubernetes_node_pool",
		Description: "GCP Kubernetes Node Pool",
		List: &plugin.ListConfig{
			Hydrate:       listKubernetesNodePools,
			ParentHydrate: listKubernetesClusters,
			Tags:          map[string]string{"service": "container", "action": "nodePools.list"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "location", "cluster_name"}),
			Hydrate:    getKubernetesNodePool,
			Tags:       map[string]string{"service": "container", "action": "nodePools.get"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the node pool.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "location_type",
				Description: "Location type of the cluster. Possible values are: 'REGIONAL', 'ZONAL'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(gcpKubernetesNodePoolLocationType),
			},
			{
				Name:        "status",
				Description: "The status of the nodes in this pool instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_name",
				Type:        proto.ColumnType_STRING,
				Description: "Cluster in which the Node pool is located.",
				Transform:   transform.FromP(kubernetesNodePoolTurbotData, "ClusterName"),
			},
			{
				Name:        "initial_node_count",
				Description: "The initial node count for the pool.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "version",
				Description: "The version of the Kubernetes of this node.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "pod_ipv4_cidr_size",
				Description: "The pod CIDR block size per node in this node pool.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "autoscaling",
				Description: "Autoscaler configuration for this node pool.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "conditions",
				Description: "Which conditions caused the current node pool state.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "config",
				Description: "The node configuration of the pool.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "instance_group_urls",
				Description: "The resource URLs of the managed instance groups associated with this node pool.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "locations",
				Description: "The list of Google Compute Engine zones.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "management",
				Description: "Node management configuration for this node pool.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "max_pods_constraint",
				Description: "The constraint on the maximum number of pods that can be run simultaneously on a node in the node pool.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "upgrade_settings",
				Description: "Upgrade settings control disruption and speed of the upgrade.",
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
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(kubernetesNodePoolTurbotData, "Akas"),
			},

			// GCP standard columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(kubernetesNodePoolTurbotData, "Location"),
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

//// LIST FUNCTION

func listKubernetesNodePools(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listKubernetesNodePools")

	// Get the details of Cluster
	cluster := h.Item.(*container.Cluster)

	// Create Service Connection
	service, err := ContainerService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	resp, err := service.Projects.Locations.Clusters.NodePools.List("projects/" + project + "/locations/" + cluster.Zone + "/clusters/" + cluster.Name).Do()
	if err != nil {
		return nil, err
	}

	for _, nodePool := range resp.NodePools {
		d.StreamLeafListItem(ctx, nodePool)
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getKubernetesNodePool(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getKubernetesNodePool")

	// Create Service Connection
	service, err := ContainerService(ctx, d)
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
	location := d.EqualsQuals["location"].GetStringValue()
	clusterName := d.EqualsQuals["cluster_name"].GetStringValue()

	// Empty check
	if name == "" || location == "" || clusterName == "" {
		return nil, nil
	}
	parent := "projects/" + project + "/locations/" + location + "/clusters/" + clusterName + "/nodePools/" + name

	resp, err := service.Projects.Locations.Clusters.NodePools.Get(parent).Do()
	// apply rate limiting
	d.WaitForListRateLimit(ctx)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

//// TRANSFORM FUNCTIONS

func kubernetesNodePoolTurbotData(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("kubernetesNodePoolTurbotData")
	nodePool := d.HydrateItem.(*container.NodePool)

	splitName := strings.Split(nodePool.SelfLink, "/")
	akas := []string{strings.Replace(nodePool.SelfLink, "https://", "gcp://", 1)}

	result := map[string]interface{}{
		"ClusterName": splitName[9],
		"Location":    splitName[7],
		"Akas":        akas,
	}
	return result[d.Param.(string)], nil
}

func gcpKubernetesNodePoolLocationType(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("gcpKubernetesNodePoolLocationType")
	nodePool := d.HydrateItem.(*container.NodePool)

	splitName := strings.Split(nodePool.SelfLink, "/")

	if splitName[6] == "locations" {
		return "REGIONAL", nil
	}
	return "ZONAL", nil
}
