package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

//// TABLE DEFINITION

func tableGcpComputeImage(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_image",
		Description: "GCP Compute Image",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "source_project"}),
			Hydrate:    getComputeImage,
			Tags:       map[string]string{"service": "compute", "action": "images.get"},
		},
		List: &plugin.ListConfig{
			ParentHydrate:     listComputeImageProjects,
			Hydrate:           listImagesForProject,
			ShouldIgnoreError: isIgnorableError([]string{"404"}),
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "deprecation_state", Require: plugin.Optional},
				{Name: "family", Require: plugin.Optional},
				{Name: "source_project", Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
				{Name: "source_type", Require: plugin.Optional},
			},
			Tags: map[string]string{"service": "compute", "action": "images.list"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getComputeImageIamPolicy,
				Tags: map[string]string{"service": "compute", "action": "images.getIamPolicy"},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "A friendly name that identifies the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique identifier for the resource.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "kind",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A user-specified, human-readable description of the image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "The server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "deprecated",
				Description: "An object comtaining the detailed deprecation status associated with this image.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "deprecation_state",
				Description: "The deprecation state associated with this image.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Deprecated.State"),
			},
			{
				Name:        "archive_size_bytes",
				Description: "Size of the image tar.gz archive stored in Google Cloud Storage (in bytes).",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "creation_timestamp",
				Description: "The creation timestamp of the resource.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "disk_size_gb",
				Description: "Size of the image when restored onto a persistent disk (in GB).",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "family",
				Description: "The name of the image family to which this image belongs.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "label_fingerprint",
				Description: "A fingerprint for the labels being applied to this image, which is essentially a hash of the labels used for optimistic locking.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_disk",
				Description: "The URL of the source disk used to create this image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_disk_id",
				Description: "The ID value of the disk used to create this image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_image",
				Description: "The URL of the source image used to create this image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_image_id",
				Description: "The ID value of the image used to create this image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_project",
				Description: "The project in which the image is defined.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(computeImageSelfLinkToTurbotData, "Project"),
			},
			{
				Name:        "source_snapshot",
				Description: "The ID value of the snapshot used to create this image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_snapshot_id",
				Description: "The ID value of the snapshot used to create this image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_type",
				Description: "The type of the image used to create this disk.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "image_encryption_key",
				Description: "The customer-supplied encryption key of the image.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "guest_os_features",
				Description: "A list of features to enable on the guest operating system.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "iam_policy",
				Description: "An Identity and Access Management (IAM) policy, which specifies access controls for Google Cloud resources. A `Policy` is a collection of `bindings`. A `binding` binds one or more `members` to a single `role`. Members can be user accounts, service accounts, Google groups, and domains (such as G Suite). A `role` is a named list of permissions; each `role` can be an IAM predefined role or a user-created custom role. For some types of Google Cloud resources, a `binding` can also specify a `condition`, which is a logical expression that allows access to a resource only if the expression evaluates to `true`.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getComputeImageIamPolicy,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "licenses",
				Description: "A list of applicable license URI.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "raw_disk",
				Description: "A set of parameters of the raw disk image.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "source_disk_encryption_key",
				Description: "The customer-supplied encryption key of the source disk.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "source_image_encryption_key",
				Description: "The customer-supplied encryption key of the source image.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "source_snapshot_encryption_key",
				Description: "The customer-supplied encryption key of the source snapshot.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "storage_locations",
				Description: "A list of Cloud Storage bucket storage location of the image (regional or multi-regional).",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "labels",
				Description: "A set of labels to apply to this image.",
				Type:        proto.ColumnType_JSON,
			},

			// standard steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Labels"),
			},
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
				Transform:   transform.FromP(computeImageSelfLinkToTurbotData, "Akas"),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromConstant("global"),
			},
			{
				Name:        "project",
				Description: "The gcp project queried.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getProject,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeImageProjects(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeImageProjects")

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	qualProjects := []string{}

	if d.EqualsQuals["source_project"] != nil {
		value := d.EqualsQuals["source_project"]
		if value.GetStringValue() != "" {
			qualProjects = []string{value.GetStringValue()}
		} else if value.GetListValue() != nil {
			qualProjects = getListValues(value.GetListValue())
		}
	}

	// List images only for requested projects
	if len(qualProjects) > 0 {
		for _, projectName := range qualProjects {
			d.StreamListItem(ctx, projectName)
		}
		return nil, nil
	}

	// List of projects in which standard images resides
	projectList := []string{
		"centos-cloud",
		"cos-cloud",
		"debian-cloud",
		"fedora-coreos-cloud",
		"rhel-cloud",
		"rhel-sap-cloud",
		"suse-cloud",
		"suse-sap-cloud",
		"ubuntu-os-cloud",
		"windows-cloud",
		"windows-sql-cloud",
		project,
	}

	for _, projectName := range projectList {
		d.StreamListItem(ctx, projectName)

		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

func listImagesForProject(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	projectName := h.Item.(string)

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	filterQuals := []filterQualMap{
		{"family", "family", "string"},
		{"deprecation_state", "deprecated.state", "string"},
		{"source_type", "sourceType", "string"},
	}

	filters := buildQueryFilter(filterQuals, d.EqualsQuals)
	filterString := ""
	if len(filters) > 0 {
		filterString = strings.Join(filters, " ")
	}

	// Max limit is set as per documentation
	// https://pkg.go.dev/google.golang.org/api@v0.48.0/compute/v1?utm_source=gopls#ImagesListCall.MaxResults
	pageSize := types.Int64(500)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *pageSize {
			pageSize = limit
		}
	}

	// resp := service.Images.List(project).Filter("deprecated.state!=\"DEPRECATED\"")
	resp := service.Images.List(projectName).MaxResults(*pageSize).Filter(filterString)
	if err := resp.Pages(ctx, func(page *compute.ImageList) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, image := range page.Items {
			d.StreamListItem(ctx, image)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				page.NextPageToken = ""
				break
			}
		}
		return nil
	}); err != nil {
		plugin.Logger(ctx).Error("listImagesForProject", "list error", err)
		// Handle project not found error
		if err.(*googleapi.Error).Code == 404 {
			return nil, nil
		}
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getComputeImage(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	name := d.EqualsQuals["name"].GetStringValue()
	project := d.EqualsQuals["source_project"].GetStringValue()

	// Error: pq: rpc error: code = Unknown desc = json: invalid use of ,string struct tag,
	// trying to unmarshal "projects/project/global/images/" into uint64
	if len(name) < 1 {
		return nil, nil
	}

	req, err := service.Images.Get(project, name).Do()
	if err != nil {
		return nil, err
	}

	return req, nil
}

func getComputeImageIamPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

	image := h.Item.(*compute.Image)
	splittedTitle := strings.Split(image.SelfLink, "/")
	imageProject := types.SafeString(splittedTitle[6])

	// If image project is not same as the project where api is called
	// do not make GetIamPolicy call as we might not have access to other project
	if !strings.EqualFold(imageProject, project) {
		return nil, nil
	}

	resp, err := service.Images.GetIamPolicy(project, image.Name).Do()
	if err != nil {
		return err, nil
	}

	return resp, nil
}

//// TRANSFORM FUNCTIONS

func computeImageSelfLinkToTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	image := d.HydrateItem.(*compute.Image)
	param := d.Param.(string)

	// get the resource title
	project := strings.Split(image.SelfLink, "/")[6]

	turbotData := map[string]interface{}{
		"Project": project,
		"Akas":    []string{"gcp://compute.googleapis.com/projects/" + project + "/global/images/" + image.Name},
	}

	return turbotData[param], nil
}
