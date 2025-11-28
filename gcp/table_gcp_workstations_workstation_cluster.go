package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"google.golang.org/api/workstations/v1"
)

//// TABLE DEFINITION

func tableGcpWorkstationsWorkstationCluster(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_workstations_workstation_cluster",
		Description: "GCP Workstations Workstation Cluster",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    workstationsWorkstationCluster,
			Tags:       map[string]string{"service": "workstations", "action": "workstations.workstationClusters.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listWorkstationsWorkstationClusters,
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "location",
					Require: plugin.Optional,
				},
			},
			Tags: map[string]string{"service": "workstations", "action": "workstations.workstationClusters.list"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The full resource name of the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "Human-readable name for this cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "uid",
				Description: "A system-assigned unique identifier for this cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "Output only. Time when this cluster was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CreateTime").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "update_time",
				Description: "Output only. Time when this cluster was most recently updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("UpdateTime").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "delete_time",
				Description: "Output only. Time when this cluster was soft-deleted.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("DeleteTime").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "etag",
				Description: "Checksum computed by the server. May be sent on update and delete requests to ensure that the client has an up-to-date value before proceeding.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "reconciling",
				Description: "Output only. Indicates whether this cluster is currently being updated to match its intended state.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     workstationsClusterSelfLink,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "network",
				Description: "The network to which the cluster is connected.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "subnetwork",
				Description: "The subnetwork to which the cluster is connected.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "private_cluster_config",
				Description: "Configuration for private cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "conditions",
				Description: "Output only. Status conditions describing the current state of the cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "labels",
				Description: "Client-specified labels.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Labels"),
			},
			{
				Name:        "annotations",
				Description: "Client-specified annotations.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(gcpWorkstationsClusterTitle),
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
				Hydrate:     workstationsClusterTurbotData,
				Transform:   transform.FromField("Akas"),
			},

			// GCP standard columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Hydrate:     workstationsClusterTurbotData,
				Transform:   transform.FromField("Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Hydrate:     workstationsClusterTurbotData,
				Transform:   transform.FromField("Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listWorkstationsWorkstationClusters(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := WorkstationsService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_workstations_workstation_cluster.listWorkstationsWorkstationClusters", "service_error", err)
		return nil, err
	}

	// Get project details
	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	// Get optional filter values from query
	location := d.EqualsQualString("location")

	// If location not specified, default to "-" for all locations
	if location == "" {
		location = "-"
	}

	// Max limit is set as per documentation
	pageSize := types.Int64(500)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *pageSize {
			pageSize = limit
		}
	}

	// Construct parent path
	parent := "projects/" + project + "/locations/" + location

	resp := service.Projects.Locations.WorkstationClusters.List(parent).PageSize(*pageSize)
	if err := resp.Pages(ctx, func(page *workstations.ListWorkstationClustersResponse) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, cluster := range page.WorkstationClusters {
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
		plugin.Logger(ctx).Error("gcp_workstations_workstation_cluster.listWorkstationsWorkstationClusters", "api_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func workstationsWorkstationCluster(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := WorkstationsService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_workstations_workstation_cluster.workstationsWorkstationCluster", "service_error", err)
		return nil, err
	}

	name := d.EqualsQuals["name"].GetStringValue()
	if name == "" {
		return nil, nil
	}

	// If the name does not include the full path, fall back to project/location qualifiers if present
	if !strings.Contains(name, "/") {
		projectId, err := getProject(ctx, d, h)
		if err != nil {
			return nil, err
		}
		project := projectId.(string)

		location := d.EqualsQuals["location"].GetStringValue()

		if location == "" {
			return nil, nil
		}

		name = "projects/" + project + "/locations/" + location + "/workstationClusters/" + name
	}

	resp, err := service.Projects.Locations.WorkstationClusters.Get(name).Do()
	if err != nil {
		plugin.Logger(ctx).Error("gcp_workstations_workstation_cluster.workstationsWorkstationCluster", "api_error", err)
		return nil, err
	}

	return resp, nil
}

//// TRANSFORM FUNCTIONS

func workstationsClusterSelfLink(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(*workstations.WorkstationCluster)
	selfLink := "https://workstations.googleapis.com/v1/" + data.Name
	return selfLink, nil
}

func gcpWorkstationsClusterTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*workstations.WorkstationCluster)

	if data.DisplayName != "" {
		return data.DisplayName, nil
	}

	// Extract name from the full resource path
	parts := strings.Split(data.Name, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1], nil
	}

	return data.Name, nil
}

func workstationsClusterTurbotData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cluster := h.Item.(*workstations.WorkstationCluster)

	// Resource name format:
	// projects/{project}/locations/{location}/workstationClusters/{cluster}
	parts := strings.Split(cluster.Name, "/")

	var project, location string
	if len(parts) >= 6 {
		project = parts[1]
		location = parts[3]
	}

	// Prefer the matrix location (if set) the same way other tables do
	matrixLocation := d.EqualsQualString(matrixKeyLocation)
	if matrixLocation != "" {
		location = matrixLocation
	}

	turbotData := map[string]interface{}{
		"Project":  project,
		"Location": location,
		"Akas":     []string{"gcp://workstations.googleapis.com/" + cluster.Name},
	}

	return turbotData, nil
}
