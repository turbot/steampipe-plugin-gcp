package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"google.golang.org/api/artifactregistry/v1"
)

//// TABLE DEFINITION

func tableGcpArtifactRegistryRepository(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_artifact_registry_repository",
		Description: "GCP Artifact Registry Repository",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "location"}),
			Hydrate:    getArtifactRegistryRepository,
		},
		List: &plugin.ListConfig{
			Hydrate: listArtifactRegistryRepositories,
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "location",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildArtifactRegistryLocationList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the repository.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(artifactRegistryRepositoryData, "Title"),
			},
			{
				Name:        "cleanup_policy_dry_run",
				Description: "If true, the cleanup pipeline is prevented from deleting versions in this repository.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "create_time",
				Description: "The time when the repository was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "The user-provided description of the repository.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "format",
				Description: "The format of packages that are stored in the repository.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kms_key_name",
				Description: "The Cloud KMS resource name of the customer managed encryption key that's used to encrypt the contents of the Repository.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "mode",
				Description: "The mode of the repository.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "satisfies_pzs",
				Description: "If set, the repository satisfies physical zone separation.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "size_bytes",
				Description: "The size, in bytes, of all artifact storage in this repository.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "update_time",
				Description: "The time when the repository was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},

			// JSON field
			{
				Name:        "cleanup_policies",
				Description: "Cleanup policies for this repository.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "docker_config",
				Description: "Docker repository config contains repository level configuration for the repositories of docker type.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "maven_config",
				Description: "Maven repository config contains repository level configuration for the repositories of maven type.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "remote_repository_config",
				Description: "Configuration specific for a Remote Repository.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "sbom_config",
				Description: "Config and state for sbom generation for resources within this Repository.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "virtual_repository_config",
				Description: "Configuration specific for a Virtual Repository.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "self_link",
				Description: "An URL that can be used to access the resource again.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     artifactRegistryRepositorySelfLink,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "labels",
				Description: "A set of labels associated with this repository.",
				Type:        proto.ColumnType_JSON,
			},

			// standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(artifactRegistryRepositoryData, "Title"),
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
				Transform:   transform.FromP(artifactRegistryRepositoryData, "Akas"),
			},

			// Standard GCP columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(artifactRegistryRepositoryData, "Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(artifactRegistryRepositoryData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listArtifactRegistryRepositories(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString("location")

	var location string
	matrixLocation := d.EqualsQualString(matrixKeyLocation)
	// Since, when the service API is disabled, matrixLocation value will be nil
	if matrixLocation != "" {
		location = matrixLocation
	}

	// Minimize API call for given location
	if region != "" && region != location {
		return nil, nil
	}

	// Create Service Connection
	service, err := ArtifactRegistryService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_artifact_registry_repository.listArtifactRegistryRepositories", "service_error", err)
		return nil, err
	}

	// Max limit isn't mentioned in the documentation
	// Default limit is set as 1000
	pageSize := types.Int64(1000)
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

	data := "projects/" + project + "/locations/" + location

	resp := service.Projects.Locations.Repositories.List(data).PageSize(*pageSize)
	if err := resp.Pages(ctx, func(page *artifactregistry.ListRepositoriesResponse) error {
		for _, repo := range page.Repositories {
			d.StreamListItem(ctx, repo)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				page.NextPageToken = ""
				return nil
			}
		}
		return nil
	}); err != nil {
		plugin.Logger(ctx).Error("gcp_artifact_registry_repository.listArtifactRegistryRepositories", "api_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getArtifactRegistryRepository(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := ArtifactRegistryService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_artifact_registry_repository.getArtifactRegistryRepository", "service_error", err)
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

	// check if name or location is empty
	if name == "" || location == "" {
		return nil, nil
	}

	resp, err := service.Projects.Locations.Repositories.Get("projects/" + project + "/locations/" + location + "/repositories/" + name).Do()
	if err != nil {
		plugin.Logger(ctx).Error("gcp_artifact_registry_repository.getArtifactRegistryRepository", "api_error", err)
		return nil, err
	}

	return resp, nil
}

func artifactRegistryRepositorySelfLink(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(*artifactregistry.Repository)

	var location string
	matrixLocation := d.EqualsQualString(matrixKeyLocation)
	// Since, when the service API is disabled, matrixLocation value will be nil
	if matrixLocation != "" {
		location = matrixLocation
	}

	projectID := strings.Split(data.Name, "/")[1]
	name := strings.Split(data.Name, "/")[5]

	selfLink := "https://artifactregistry.googleapis.com/v1/projects/" + projectID + "/locations/" + location + "/repositories/" + name

	return selfLink, nil
}

//// TRANSFORM FUNCTIONS

func artifactRegistryRepositoryData(ctx context.Context, h *transform.TransformData) (interface{}, error) {
	data := h.HydrateItem.(*artifactregistry.Repository)
	param := h.Param.(string)

	projectID := strings.Split(data.Name, "/")[1]
	name := strings.Split(data.Name, "/")[5]
	location := strings.Split(data.Name, "/")[3]

	turbotData := map[string]interface{}{
		"Project":  projectID,
		"Title":    name,
		"Location": location,
		"Akas":     []string{"gcp://artifactregistry.googleapis.com/projects/" + projectID + "/repositories/" + name},
	}

	return turbotData[param], nil
}
