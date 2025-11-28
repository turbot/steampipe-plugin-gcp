package gcp

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"google.golang.org/api/workstations/v1"
)

//// TABLE DEFINITION

func tableGcpWorkstationsWorkstation(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_workstations_workstation",
		Description: "GCP Workstations workstation",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getWorkstationsWorkstation,
			Tags:       map[string]string{"service": "workstations", "action": "workstations.workstations.get"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listWorkstationsWorkstationClusters,
			Hydrate:       listWorkstationsWorkstations,
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "location",
					Require: plugin.Optional,
				},
				{
					Name:    "cluster",
					Require: plugin.Optional,
				},
				{
					Name:    "config",
					Require: plugin.Optional,
				},
			},
			Tags: map[string]string{"service": "workstations", "action": "workstations.workstations.list"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getWorkstationsWorkstationIamPolicy,
				Tags: map[string]string{"service": "workstations", "action": "workstations.workstations.getIamPolicy"},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The full resource name of the workstation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster",
				Description: "The workstation cluster containing this workstation.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     workstationsWorkstationTurbotData,
				Transform:   transform.FromField("Cluster"),
			},
			{
				Name:        "config",
				Description: "The workstation configuration associated with this workstation.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     workstationsWorkstationTurbotData,
				Transform:   transform.FromField("Config"),
			},
			{
				Name:        "display_name",
				Description: "Human-readable name for this workstation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "uid",
				Description: "A system-assigned unique identifier for this workstation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "Output only. Current state of the workstation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "host",
				Description: "Output only. Host to which clients can send HTTPS traffic that will be received by the workstation. Authorized traffic will be received to the workstation as HTTP on port 80.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "Output only. Time when this workstation was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CreateTime").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "update_time",
				Description: "Output only. Time when this workstation was most recently updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("UpdateTime").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "start_time",
				Description: "Output only. Time when this workstation was most recently successfully started, regardless of the workstation's initial state.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("StartTime").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "delete_time",
				Description: "Output only. Time when this workstation was soft-deleted.",
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
				Description: "Output only. Indicates whether this workstation is currently being updated to match its intended state.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     workstationsWorkstationSelfLink,
				Transform:   transform.FromValue(),
			},

			// JSON fields
			{
				Name:        "annotations",
				Description: "Client-specified annotations.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "labels",
				Description: "Client-specified labels.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "env",
				Description: "Environment variables passed to the workstation container's entrypoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "iam_policy",
				Description: "An Identity and Access Management (IAM) policy, which specifies access controls for Google Cloud resources.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWorkstationsWorkstationIamPolicy,
				Transform:   transform.FromValue(),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(workstationsWorkstationTitle),
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
				Hydrate:     workstationsWorkstationTurbotData,
				Transform:   transform.FromField("Akas"),
			},

			// GCP standard columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Hydrate:     workstationsWorkstationTurbotData,
				Transform:   transform.FromField("Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Hydrate:     workstationsWorkstationTurbotData,
				Transform:   transform.FromField("Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listWorkstationsWorkstations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the cluster from parent hydrate
	cluster := h.Item.(*workstations.WorkstationCluster)

	// Get optional filter values from query
	locationQual := d.EqualsQualString("location")
	clusterQual := d.EqualsQualString("cluster")
	configQual := d.EqualsQualString("config")

	// Extract cluster name from the full resource path
	// Format: projects/{project}/locations/{location}/workstationClusters/{cluster}
	clusterParts := strings.Split(cluster.Name, "/")
	if len(clusterParts) < 6 {
		return nil, nil
	}

	clusterName := clusterParts[5]
	locationName := clusterParts[3]

	// If location qualifier is provided and doesn't match, skip this cluster
	if locationQual != "" && locationQual != locationName {
		return nil, nil
	}

	// If cluster qualifier is provided and doesn't match, skip this cluster
	if clusterQual != "" && clusterQual != clusterName {
		return nil, nil
	}

	// Create Service Connection
	service, err := WorkstationsService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_workstations_workstation.listWorkstationsWorkstations", "service_error", err)
		return nil, err
	}

	// List configs for this cluster
	// Path format: projects/{project}/locations/{location}/workstationClusters/{cluster}
	// The List() method already knows to list workstationConfigs, so we just pass the cluster path
	configsResp := service.Projects.Locations.WorkstationClusters.WorkstationConfigs.List(cluster.Name)
	if err := configsResp.Pages(ctx, func(configsPage *workstations.ListWorkstationConfigsResponse) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, config := range configsPage.WorkstationConfigs {
			// Extract config name from the full resource path
			configParts := strings.Split(config.Name, "/")
			if len(configParts) < 8 {
				continue
			}
			configName := configParts[7]

			// If config qualifier is provided and doesn't match, skip this config
			if configQual != "" && configQual != configName {
				continue
			}

			// List workstations for this config
			// Path format: projects/{project}/locations/{location}/workstationClusters/{cluster}/workstationConfigs/{config}
			// The List() method already knows to list workstations, so we just pass the config path
			workstationsResp := service.Projects.Locations.WorkstationClusters.WorkstationConfigs.Workstations.List(config.Name)
			if err := workstationsResp.Pages(ctx, func(workstationsPage *workstations.ListWorkstationsResponse) error {
				// apply rate limiting
				d.WaitForListRateLimit(ctx)

				for _, workstation := range workstationsPage.Workstations {
					d.StreamListItem(ctx, workstation)

					// Check if context has been cancelled or if the limit has been hit (if specified)
					if d.RowsRemaining(ctx) == 0 {
						workstationsPage.NextPageToken = ""
						return nil
					}
				}
				return nil
			}); err != nil {
				plugin.Logger(ctx).Error("gcp_workstations_workstation.listWorkstationsWorkstations", "list_workstations_error", err, "config", config.Name)
				// Continue with next config even if this one fails
			}

			// Check if we've hit the limit
			if d.RowsRemaining(ctx) == 0 {
				configsPage.NextPageToken = ""
				return nil
			}
		}
		return nil
	}); err != nil {
		plugin.Logger(ctx).Error("gcp_workstations_workstation.listWorkstationsWorkstations", "list_configs_error", err, "cluster", cluster.Name)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getWorkstationsWorkstation(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := WorkstationsService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_workstations_workstation.getWorkstationsWorkstation", "service_error", err)
		return nil, err
	}

	name := d.EqualsQuals["name"].GetStringValue()
	if name == "" {
		return nil, nil
	}

	// If the name does not include the full path, fall back to project/loc qualifiers if present
	if !strings.Contains(name, "/") {
		projectId, err := getProject(ctx, d, h)
		if err != nil {
			return nil, err
		}
		project := projectId.(string)

		location := d.EqualsQuals["location"].GetStringValue()
		cluster := d.EqualsQuals["cluster"].GetStringValue()
		config := d.EqualsQuals["config"].GetStringValue()

		if location == "" || cluster == "" || config == "" {
			return nil, nil
		}

		name = "projects/" + project + "/locations/" + location + "/workstationClusters/" + cluster + "/workstationConfigs/" + config + "/workstations/" + name
	}

	resp, err := service.Projects.Locations.WorkstationClusters.WorkstationConfigs.Workstations.Get(name).Do()
	if err != nil {
		plugin.Logger(ctx).Error("gcp_workstations_workstation.getWorkstationsWorkstation", "api_error", err)
		return nil, err
	}

	return resp, nil
}

func getWorkstationsWorkstationIamPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(*workstations.Workstation)

	// Create Service Connection
	service, err := WorkstationsService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_workstations_workstation.getWorkstationsWorkstationIamPolicy", "service_error", err)
		return nil, err
	}

	resp, err := service.Projects.Locations.WorkstationClusters.WorkstationConfigs.Workstations.GetIamPolicy(data.Name).Do()
	if err != nil {
		plugin.Logger(ctx).Error("gcp_workstations_workstation.getWorkstationsWorkstationIamPolicy", "api_error", err)
		return nil, err
	}

	return resp, nil
}

//// TRANSFORM FUNCTIONS

func workstationsWorkstationSelfLink(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(*workstations.Workstation)
	selfLink := "https://workstations.googleapis.com/v1/" + data.Name
	return selfLink, nil
}

func workstationsWorkstationTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*workstations.Workstation)

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

func workstationsWorkstationTurbotData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	workstation := h.Item.(*workstations.Workstation)

	parts := strings.Split(workstation.Name, "/")

	var project, location, cluster, config string
	if len(parts) >= 10 {
		project = parts[1]
		location = parts[3]
		cluster = parts[5]
		config = parts[7]
	}

	matrixLocation := d.EqualsQualString(matrixKeyLocation)
	if matrixLocation != "" {
		location = matrixLocation
	}

	turbotData := map[string]interface{}{
		"Project":  project,
		"Location": location,
		"Cluster":  cluster,
		"Config":   config,
		"Akas":     []string{"gcp://workstations.googleapis.com/" + workstation.Name},
	}

	return turbotData, nil
}
