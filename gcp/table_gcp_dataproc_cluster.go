package gcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"google.golang.org/api/dataproc/v1"
)

func tableGcpDataprocCluster(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_dataproc_cluster",
		Description: "GCP Dataproc Cluster",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("cluster_name"),
			Hydrate:    getDataprocCluster,
			Tags:       map[string]string{"service": "dataproc", "action": "clusters.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listDataprocClusters,
			KeyColumns: plugin.KeyColumnSlice{
				// String columns
				{Name: "cluster_name", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "state", Require: plugin.Optional, Operators: []string{"="}},
			},
			Tags: map[string]string{"service": "dataproc", "action": "clusters.list"},
		},
		GetMatrixItemFunc: BuildComputeLocationList,
		Columns: []*plugin.Column{
			// commonly used columns
			{
				Name:        "cluster_name",
				Description: "The cluster name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_uuid",
				Description: "A cluster UUID (Unique Universal Identifier). Dataproc generates this value when it creates the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The cluster's state.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Status.State"),
			},
			{
				Name:        "config",
				Description: "The cluster config.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "labels",
				Description: "The labels to associate with this cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "metrics",
				Description: "Contains cluster daemon metrics such as HDFS and YARN stats.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "status",
				Description: "Cluster status.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "status_history",
				Description: "The previous cluster status.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     dataprocClusterSelfLink,
				Transform:   transform.FromValue(),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterName"),
			},
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Labels"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     gcpDataprocClusterTurbotData,
				Transform:   transform.FromField("Akas"),
			},

			// Standard GCP columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Hydrate:     gcpDataprocClusterTurbotData,
				Transform:   transform.FromField("Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProjectId"),
			},
		},
	}
}

//// LIST FUNCTION

func listDataprocClusters(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var location string
	matrixLocation := d.EqualsQualString(matrixKeyLocation)
	// Since, when the service API is disabled, matrixLocation value will be nil
	if matrixLocation != "" {
		location = matrixLocation
	}

	// Create Service Connection
	service, err := DataprocService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_dataproc_cluster.listDataprocClusters", "connection_error", err)
		return nil, err
	}

	var filters []string
	if d.EqualsQualString("cluster_name") != "" {
		filters = append(filters, fmt.Sprint("clusterName = ", d.EqualsQualString("cluster_name")))
	}

	if d.EqualsQualString("state") != "" {
		filters = append(filters, fmt.Sprint("status.state = ", d.EqualsQualString("state")))
	}

	filterString := ""
	if len(filters) > 0 {
		filterString = strings.Join(filters, "AND")
	}

	// Max limit is set as per documentation
	pageSize := types.Int64(500)
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

	resp := service.Projects.Regions.Clusters.List(project, location).PageSize(*pageSize).Filter(filterString)
	if err := resp.Pages(ctx, func(page *dataproc.ListClustersResponse) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

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
		plugin.Logger(ctx).Error("gcp_dataproc_cluster.listDataprocClusters", "api_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getDataprocCluster(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var location string
	matrixLocation := d.EqualsQualString(matrixKeyLocation)
	// Since, when the service API is disabled, matrixLocation value will be nil
	if matrixLocation != "" {
		location = matrixLocation
	}

	clusterName := d.EqualsQuals["cluster_name"].GetStringValue()

	if len(clusterName) < 1 {
		return nil, nil
	}

	// Create Service Connection
	service, err := DataprocService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_dataproc_cluster.getDataprocCluster", "connection_error", err)
		return nil, err
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	resp, err := service.Projects.Regions.Clusters.Get(project, location, clusterName).Do()
	if err != nil {
		plugin.Logger(ctx).Error("gcp_dataproc_cluster.getDataprocCluster", "api_error", err)
		return nil, err
	}

	return resp, nil
}

//// TRANSFORM FUNCTION

func gcpDataprocClusterTurbotData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cluster := h.Item.(*dataproc.Cluster)

	project := cluster.ProjectId
	var location string
	matrixLocation := d.EqualsQualString(matrixKeyLocation)
	// Since, when the service API is disabled, matrixLocation value will be nil
	if matrixLocation != "" {
		location = matrixLocation
	}

	turbotData := map[string]interface{}{
		"Project":  project,
		"Location": location,
		"Akas":     []string{"gcp://dataproc.googleapis.com/projects/" + project + "/regions/" + location + "/clusters/" + cluster.ClusterName},
	}

	return turbotData, nil
}

func dataprocClusterSelfLink(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(*dataproc.Cluster)

	var location string
	matrixLocation := d.EqualsQualString(matrixKeyLocation)
	// Since, when the service API is disabled, matrixLocation value will be nil
	if matrixLocation != "" {
		location = matrixLocation
	}

	selfLink := "https://dataproc.googleapis.com/v1/projects/" + data.ProjectId + "/regions/" + location + "/clusters/" + data.ClusterName

	return selfLink, nil
}
