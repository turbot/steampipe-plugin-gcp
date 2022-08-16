package gcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/dataproc/v1"
)

func tableGcpDataprocCluster(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_dataproc_cluster",
		Description: "GCP Dataproc Cluster",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("cluster_name"),
			Hydrate:    getDataprocCluster,
		},
		List: &plugin.ListConfig{
			Hydrate: listDataprocClusters,
			KeyColumns: plugin.KeyColumnSlice{
				// String columns
				{Name: "cluster_name", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "state", Require: plugin.Optional, Operators: []string{"="}},
			},
		},
		GetMatrixItem: BuildComputeLocationList,
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
				Transform:   transform.FromP(gcpDataprocClusterTurbotData, "Akas"),
			},

			// Standard GCP columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpDataprocClusterTurbotData, "Location"),
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
	matrixLocation := plugin.GetMatrixItem(ctx)[matrixKeyLocation]
	// Since, when the service API is disabled, matrixLocation value will be nil
	if matrixLocation != nil {
		location = matrixLocation.(*compute.Region).Name
	}

	// Create Service Connection
	service, err := DataprocService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_dataproc_cluster.listDataprocClusters", "connection_error", err)
		return nil, err
	}

var filters []string
	if d.KeyColumnQualString("cluster_name") != "" {
		filters = append(filters, fmt.Sprint("clusterName = ", d.KeyColumnQualString("cluster_name")))
	}

	if d.KeyColumnQualString("state") != "" {
		filters = append(filters, fmt.Sprint("status.state = ", d.KeyColumnQualString("state")))
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
	getProjectCached := plugin.HydrateFunc(getProject).WithCache()
	projectId, err := getProjectCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	resp := service.Projects.Regions.Clusters.List(project, location).PageSize(*pageSize).Filter(filterString)
	if err := resp.Pages(ctx, func(page *dataproc.ListClustersResponse) error {
		for _, cluster := range page.Clusters {
			d.StreamListItem(ctx, cluster)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
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
	matrixLocation := plugin.GetMatrixItem(ctx)[matrixKeyLocation]
	// Since, when the service API is disabled, matrixLocation value will be nil
	if matrixLocation != nil {
		location = matrixLocation.(*compute.Region).Name
	}

	clusterName := d.KeyColumnQuals["cluster_name"].GetStringValue()

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
	getProjectCached := plugin.HydrateFunc(getProject).WithCache()
	projectId, err := getProjectCached(ctx, d, h)
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

func gcpDataprocClusterTurbotData(ctx context.Context, d *transform.TransformData) (interface{}, error) {

	cluster := d.HydrateItem.(*dataproc.Cluster)
	param := d.Param.(string)

	project := cluster.ProjectId
	var location string
	matrixLocation := plugin.GetMatrixItem(ctx)[matrixKeyLocation]
	// Since, when the service API is disabled, matrixLocation value will be nil
	if matrixLocation != nil {
		location = matrixLocation.(*compute.Region).Name
	}

	turbotData := map[string]interface{}{
		"Project":  project,
		"Location": location,
		"Akas":     []string{"gcp://dataproc.googleapis.com/projects/" + project + "/regions/" + location + "/clusters/" + cluster.ClusterName},
	}

	return turbotData[param], nil
}
