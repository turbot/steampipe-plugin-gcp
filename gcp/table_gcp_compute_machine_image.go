package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/compute/v1"
)

//// TABLE DEFINITION

func tableGcpComputeMachineImage(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_machine_image",
		Description: "GCP Compute Machine Image",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeMachineImage,
			Tags:       map[string]string{"service": "compute", "action": "machineImages.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeMachineImages,
			Tags:    map[string]string{"service": "compute", "action": "machineImages.list"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "Name of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "A unique identifier for this machine image. The server defines this identifier.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "self_link",
				Description: "The URL for this machine image. The server defines this URL.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_timestamp",
				Description: "The creation timestamp for this machine image in RFC3339 text format.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromGo().NullIfZero(),
			},
			{
				Name:        "description",
				Description: "An optional description of this resource. Provide this property when you create the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "guest_flush",
				Description: "Whether to attempt an application consistent machine image by informing the OS to prepare for the snapshot process.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "kind",
				Description: "The resource type, which is always compute#machineImage for machine image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_instance",
				Description: "The source instance used to create the machine image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the machine image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "total_storage_bytes",
				Description: "Total size of the storage used by the machine image.",
				Type:        proto.ColumnType_INT,
			},

			// JSON columns
			{
				Name:        "instance_properties",
				Description: "Properties of source instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "machine_image_encryption_key",
				Description: "Encrypts the machine image using a customer-supplied encryption key. After you encrypt a machine image using a customer-supplied key, you must provide the same key if you use the machine image later.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "saved_disks",
				Description: "An array of Machine Image specific properties for disks attached to the source instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "source_disk_encryption_keys",
				Description: "The customer-supplied encryption key of the disks attached to the source instance. Required if the source disk is protected by a customer-supplied encryption key.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "storage_locations",
				Description: "The regional or multi-regional Cloud Storage bucket location where the machine image is stored.",
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
				Transform:   transform.FromP(machineImageTurbotData, "Akas"),
			},

			// GCP standard columns
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(machineImageTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeMachineImages(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_compute_machine_image.listComputeMachineImages", "connection_error", err)
		return nil, err
	}

	// Max limit is set as per documentation
	// https://cloud.google.com/compute/docs/reference/rest/v1/machineImages/list#query-parameters
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
		plugin.Logger(ctx).Error("gcp_compute_machine_image.listComputeMachineImages.getProjectCached", "cached_function", err)
		return nil, err
	}
	project := projectId.(string)

	resp := service.MachineImages.List(project).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *compute.MachineImageList) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, machineImage := range page.Items {
			d.StreamListItem(ctx, machineImage)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				page.NextPageToken = ""
				return nil
			}
		}
		return nil
	}); err != nil {
		plugin.Logger(ctx).Error("gcp_compute_machine_image.listComputeMachineImages", "api_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getComputeMachineImage(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_compute_machine_image.getComputeMachineImage", "connection_error", err)
		return nil, err
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_compute_machine_image.getComputeMachineImage.getProjectCached", "cache_error", err)
		return nil, err
	}
	project := projectId.(string)
	machineImageName := d.EqualsQualString("name")

	// Return nil, if no input provided
	if machineImageName == "" {
		return nil, nil
	}

	resp, err := service.MachineImages.Get(project, machineImageName).Do()
	if err != nil {
		plugin.Logger(ctx).Error("gcp_compute_machine_image.getComputeMachineImage", "api_error", err)
		return nil, err
	}
	return resp, nil
}

//// TRANSFORM FUNCTIONS

func machineImageTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*compute.MachineImage)
	param := d.Param.(string)

	project := strings.Split(data.SelfLink, "/")[6]

	turbotData := map[string]interface{}{
		"Project": project,
		"Akas":    []string{"gcp://compute.googleapis.com/projects/" + project + "/machineImages/" + data.Name},
	}

	return turbotData[param], nil
}
