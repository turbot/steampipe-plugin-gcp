package gcp

import (
	"context"
	"strings"
	"sync"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"google.golang.org/api/compute/v1"
)

//// TABLE DEFINITION

func tableGcpComputeImage(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_image",
		Description: "GCP Compute Image",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "project"}),
			Hydrate:    getComputeImage,
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeImages,
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
			// {
			// 	Name:        "iam_policy",
			// 	Description: "An Identity and Access Management (IAM) policy, which specifies access controls for Google Cloud resources. A `Policy` is a collection of `bindings`. A `binding` binds one or more `members` to a single `role`. Members can be user accounts, service accounts, Google groups, and domains (such as G Suite). A `role` is a named list of permissions; each `role` can be an IAM predefined role or a user-created custom role. For some types of Google Cloud resources, a `binding` can also specify a `condition`, which is a logical expression that allows access to a resource only if the expression evaluates to `true`.",
			// 	Type:        proto.ColumnType_JSON,
			// 	Hydrate:     getComputeImageIamPolicy,
			// 	Transform:   transform.FromValue(),
			// },
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
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(computeImageSelfLinkToTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeImages(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeImages")
	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}
	project := activeProject()

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

	// List all the standard images
	var wg sync.WaitGroup
	imageCh := make(chan []*compute.Image, len(projectList))
	errorCh := make(chan error, len(projectList))

	// Iterating all the available projects
	for _, item := range projectList {
		wg.Add(1)
		go getRowDataForImageAsync(ctx, service, item, &wg, imageCh, errorCh)
	}

	// wait until the listing of all images inside the mentioned projects are completed
	wg.Wait()

	// NOTE: close channel before ranging over results
	close(imageCh)
	close(errorCh)

	for err := range errorCh {
		// return the first error
		return nil, err
	}

	for item := range imageCh {
		for _, data := range item {
			d.StreamListItem(ctx, data)
		}
	}

	return nil, nil
}

func getRowDataForImageAsync(ctx context.Context, service *compute.Service, item string, wg *sync.WaitGroup, imageCh chan []*compute.Image, errorCh chan error) {
	defer wg.Done()

	rowData, err := getRowDataForImage(ctx, service, item)
	if err != nil {
		errorCh <- err
	} else if rowData != nil {
		imageCh <- rowData
	}
}

func getRowDataForImage(ctx context.Context, service *compute.Service, project string) ([]*compute.Image, error) {
	var items []*compute.Image
	resp := service.Images.List(project)
	if err := resp.Pages(ctx, func(page *compute.ImageList) error {
		for _, image := range page.Items {
			items = append(items, image)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return items, nil
}

//// HYDRATE FUNCTIONS

func getComputeImage(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	name := d.KeyColumnQuals["name"].GetStringValue()
	project := d.KeyColumnQuals["project"].GetStringValue()
	// project := activeProject()

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

// func getComputeImageIamPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 	service, err := compute.NewService(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	image := h.Item.(*compute.Image)
// 	project := activeProject()

// 	req, err := service.Images.GetIamPolicy(project, image.Name).Do()
// 	if err != nil {
// 		// Return nil, if the resource not present
// 		result := isNotFoundError([]string{"404"})
// 		if result != nil {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}

// 	return req, nil
// }

//// TRANSFORM FUNCTIONS

func computeImageSelfLinkToTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	image := d.HydrateItem.(*compute.Image)
	param := d.Param.(string)

	// get the resource title
	splittedTitle := strings.Split(image.SelfLink, "/")

	turbotData := map[string]interface{}{
		"Project": splittedTitle[6],
		"Akas":    []string{"gcp://compute.googleapis.com/projects/" + splittedTitle[6] + "/global/images/" + image.Name},
	}

	return turbotData[param], nil
}
