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

func tableGcpWorkstationsWorkstation(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_workstations_workstation",
		Description: "GCP Workstations workstation",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getGcpWorkstationsWorkstation,
			Tags:       map[string]string{"service": "workstations", "action": "workstations.workstations.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listGcpWorkstationsWorkstations,
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
				Func: getGcpWorkstationsWorkstationIamPolicy,
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
				Hydrate:     gcpWorkstationsWorkstationTurbotData,
				Transform:   transform.FromField("Cluster"),
			},
			{
				Name:        "config",
				Description: "The workstation configuration associated with this workstation.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     gcpWorkstationsWorkstationTurbotData,
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
				Hydrate:     gcpWorkstationsWorkstationSelfLink,
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
				Hydrate:     getGcpWorkstationsWorkstationIamPolicy,
				Transform:   transform.FromValue(),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(gcpWorkstationsWorkstationTitle),
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
				Hydrate:     gcpWorkstationsWorkstationTurbotData,
				Transform:   transform.FromField("Akas"),
			},

			// GCP standard columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Hydrate:     gcpWorkstationsWorkstationTurbotData,
				Transform:   transform.FromField("Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Hydrate:     gcpWorkstationsWorkstationTurbotData,
				Transform:   transform.FromField("Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listGcpWorkstationsWorkstations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := WorkstationsService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_workstations_workstation.listGcpWorkstationsWorkstations", "service_error", err)
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
	cluster := d.EqualsQualString("cluster")
	config := d.EqualsQualString("config")

	// If location not specified, default to "-" for all locations
	if location == "" {
		location = "-"
	}

	// If cluster not specified, default to "-" for all clusters
	if cluster == "" {
		cluster = "-"
	}

	// If config not specified, default to "-" for all configs
	if config == "" {
		config = "-"
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
	parent := "projects/" + project + "/locations/" + location + "/workstationClusters/" + cluster + "/workstationConfigs/" + config

	resp := service.Projects.Locations.WorkstationClusters.WorkstationConfigs.Workstations.List(parent).PageSize(*pageSize)
	if err := resp.Pages(ctx, func(page *workstations.ListWorkstationsResponse) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, workstation := range page.Workstations {
			d.StreamListItem(ctx, workstation)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				page.NextPageToken = ""
				return nil
			}
		}
		return nil
	}); err != nil {
		plugin.Logger(ctx).Error("gcp_workstations_workstation.listGcpWorkstationsWorkstations", "api_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getGcpWorkstationsWorkstation(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := WorkstationsService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_workstations_workstation.getGcpWorkstationsWorkstation", "service_error", err)
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
		plugin.Logger(ctx).Error("gcp_workstations_workstation.getGcpWorkstationsWorkstation", "api_error", err)
		return nil, err
	}

	return resp, nil
}

func getGcpWorkstationsWorkstationIamPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(*workstations.Workstation)

	// Create Service Connection
	service, err := WorkstationsService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_workstations_workstation.getGcpWorkstationsWorkstationIamPolicy", "service_error", err)
		return nil, err
	}

	resp, err := service.Projects.Locations.WorkstationClusters.WorkstationConfigs.Workstations.GetIamPolicy(data.Name).Do()
	if err != nil {
		plugin.Logger(ctx).Error("gcp_workstations_workstation.getGcpWorkstationsWorkstationIamPolicy", "api_error", err)
		return nil, err
	}

	return resp, nil
}

//// TRANSFORM FUNCTIONS

func gcpWorkstationsWorkstationSelfLink(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(*workstations.Workstation)
	selfLink := "https://workstations.googleapis.com/v1/" + data.Name
	return selfLink, nil
}

func gcpWorkstationsWorkstationTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
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

func gcpWorkstationsWorkstationTurbotData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
