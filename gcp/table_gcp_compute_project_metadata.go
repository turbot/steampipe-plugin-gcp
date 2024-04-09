package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGcpComputeProjectMetadata(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_project_metadata",
		Description: "GCP Compute Project Metadata",
		List: &plugin.ListConfig{
			Hydrate: listComputeProjectMetadata,
			Tags:    map[string]string{"service": "compute", "action": "projects.get"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The ID of the project.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique identifier for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "default_service_account",
				Description: "Default service account used by VMs running in this project.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "An optional textual description of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_timestamp",
				Description: "Creation timestamp in RFC3339 text format.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "default_network_tier",
				Description: "This signifies the default network tier used for configuring resources of the project and can only take the following values: PREMIUM, STANDARD.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "xpn_project_status",
				Description: "The role this project has in a shared VPC configuration.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "common_instance_metadata",
				Description: "Metadata key/value pairs available to all instances contained in this project.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "enabled_features",
				Description: "Restricted features enabled for use on this project.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "quotas",
				Description: "Quotas assigned to this project.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "usage_export_location",
				Description: "The naming prefix for daily usage reports and the Google Cloud Storage bucket where they are stored.",
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
				Hydrate:     getComputeProjectMetadataTurbotData,
			},

			// Standard GCP columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromConstant("global"),
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

func listComputeProjectMetadata(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)
	plugin.Logger(ctx).Trace("listComputeProjectMetadata", "GCP_PROJECT: ", project)

	resp, err := service.Projects.Get(project).Do()
	if err != nil {
		return nil, err
	}
	d.StreamListItem(ctx, resp)

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getComputeProjectMetadataTurbotData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	// Build resource aka
	akas := []string{"gcp://cloudresourcemanager.googleapis.com/projects/" + project}

	// Mapping all turbot defined properties
	turbotData := map[string]interface{}{
		"Akas": akas,
	}

	return turbotData, nil
}
