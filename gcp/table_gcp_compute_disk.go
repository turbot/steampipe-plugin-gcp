package gcp

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	"google.golang.org/api/compute/v1"
)

func tableGcpComputeDisk(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_disk",
		Description: "GCP Compute Disk",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeDisk,
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeDisk,
		},
		Columns: []*plugin.Column{
			// commonly used columns
			{
				Name:        "name",
				Description: "Name of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique identifier for the resource. This identifier is defined by the server.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "creation_timestamp",
				Description: "Timestamp when the disk was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "size_gb",
				Description: "Size, in GB, of the persistent disk.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "status",
				Description: "The status of disk creation. CREATING: Disk is provisioning. RESTORING: Source data is being copied into the disk. FAILED: Disk creation failed. READY: Disk is ready for use. DELETING: Disk is deleting.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "An optional description of this resource. Provide this property when you create the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "disk_encryption_kms_key",
				Description: "The name of the encryption key that is used to encrypt storage data.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DiskEncryptionKey.KmsKeyName"),
			},
			{
				Name:        "disk_encryption_kms_key_service_account",
				Description: "The service account being used for the encryption request for the given KMS key. If absent, the Compute Engine default service account is used.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DiskEncryptionKey.KmsKeyServiceAccount"),
			},
			{
				Name:        "kind",
				Description: "Type of the resource. Always compute#disk for disks.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_attach_timestamp",
				Description: "Timestamp when the disk was last attached.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromGo().NullIfZero(),
			},
			{
				Name:        "last_detach_timestamp",
				Description: "Timestamp when the disk was last detached.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromGo().NullIfZero(),
			},
			{
				Name:        "physical_block_size_bytes",
				Description: "Physical block size of the persistent disk, in bytes. If not present in a request, a default value is used.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "self_link",
				Description: "Server-defined fully-qualified URL for this resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_disk",
				Description: "The source disk used to create this disk. You can provide this as a partial or full URL to the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_disk_id",
				Description: "The unique ID of the disk used to create this disk. This value identifies the exact disk that was used to create this persistent disk.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_image",
				Description: "The source image used to create this disk. If the source image is deleted, this field will not be set.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_image_id",
				Description: "The ID value of the image used to create this disk. This value identifies the exact image that was used to create this persistent disk.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_snapshot",
				Description: "The source snapshot used to create this disk.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_snapshot_id",
				Description: "The unique ID of the snapshot used to create this disk. This value identifies the exact snapshot that was used to create this persistent disk.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_image_encryption_key",
				Description: "The customer-supplied encryption key of the source image. Required if the source image is protected by a customer-supplied encryption key.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_snapshot_encryption_key",
				Description: "The customer-supplied encryption key of the source snapshot. Required if the source snapshot is protected by a customer-supplied encryption key.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "URL of the disk type resource describing which disk type to use to create the disk. Provide this when creating the disk. For example: projects/project/zones/zone/diskTypes/pd-standard  or pd-ssd",
				Type:        proto.ColumnType_STRING,
			},
			// type_name is a simpler view of the type, without the full path
			{
				Name:        "type_name",
				Description: "Type of the disk. For example: pd-standard or pd-ssd",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Type").Transform(lastPathElement),
			},
			{
				Name:        "location_type",
				Description: "Loation type where the disk resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(diskLocation, "Type"),
			},
			{
				Name:        "region",
				Description: "URL of the region where the disk resides. Only applicable for regional resources.",
				Type:        proto.ColumnType_STRING,
			},
			// region_name is a simpler view of the region, without the full path
			{
				Name:        "region_name",
				Description: "Name of the region where the disk resides. Only applicable for regional resources.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Region").Transform(lastPathElement),
			},
			{
				Name:        "zone",
				Description: "URL of the zone where the disk resides.",
				Type:        proto.ColumnType_STRING,
			},
			// zone_name is a simpler view of the zone, without the full path
			{
				Name:        "zone_name",
				Description: "The zone name in which the disk resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Zone").Transform(lastPathElement),
			},
			{
				Name:        "guest_os_features",
				Description: "A list of features to enable on the guest operating system. Applicable only for bootable images.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "labels",
				Description: "A map of labels assigned to bucket",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "licenses",
				Description: "A list of publicly visible licenses.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "license_codes",
				Description: "Integer license codes indicating which licenses are attached to this disk.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "replica_zones",
				Description: "URLs of the zones where the disk should be replicated to. Only applicable for regional resources.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "resource_policies",
				Description: "Resource policies applied to this disk for automatic snapshot creations.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "users",
				Description: "Links to the users of the disk (attached instances) in form: projects/project/zones/zone/instances/instance",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "iam_policy",
				Description: "An Identity and Access Management (IAM) policy, which specifies access controls for Google Cloud resources. A `Policy` is a collection of `bindings`. A `binding` binds one or more `members` to a single `role`. Members can be user accounts, service accounts, Google groups, and domains (such as G Suite). A `role` is a named list of permissions; each `role` can be an IAM predefined role or a user-created custom role. For some types of Google Cloud resources, a `binding` can also specify a `condition`, which is a logical expression that allows access to a resource only if the expression evaluates to `true`.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getComputeDiskIamPolicy,
				Transform:   transform.FromValue(),
			},

			// standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
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
				Transform:   transform.From(diskAka),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(diskLocation, "Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromConstant(activeProject()),
			},
		},
	}
}

//// LIST FUNCTIONS

func listComputeDisk(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeDisk")

	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	project := activeProject()
	resp := service.Disks.AggregatedList(project)
	if err := resp.Pages(ctx, func(page *compute.DiskAggregatedList) error {
		for _, item := range page.Items {
			for _, disk := range item.Disks {
				d.StreamListItem(ctx, disk)
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getComputeDisk(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getComputeDisk")
	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	var disk compute.Disk
	project := activeProject()
	name := d.KeyColumnQuals["name"].GetStringValue()

	resp := service.Disks.AggregatedList(project).Filter("name=" + name)
	if err := resp.Pages(ctx, func(page *compute.DiskAggregatedList) error {
		for _, item := range page.Items {
			for _, i := range item.Disks {
				disk = *i
			}
		}
		return nil
	},
	); err != nil {
		return nil, err
	}

	return &disk, nil
}

func getComputeDiskIamPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	disk := h.Item.(*compute.Disk)

	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	var resp *compute.Policy
	project := activeProject()
	zoneName := getLastPathElement(types.SafeString(disk.Zone))

	// disk can be regional or zonal
	if zoneName == "" {
		regionName := getLastPathElement(types.SafeString(disk.Region))
		// regional disk get iam policy
		resp, err = service.RegionDisks.GetIamPolicy(project, regionName, disk.Name).Do()
		if err != nil {
			return nil, err
		}
		return resp, nil
	}

	// zonal disk get iam policy
	resp, err = service.Disks.GetIamPolicy(project, zoneName, disk.Name).Do()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

//// TRANSFORM FUNCTIONS

func diskAka(_ context.Context, d *transform.TransformData) (interface{}, error) {
	i := d.HydrateItem.(*compute.Disk)

	zoneName := getLastPathElement(types.SafeString(i.Zone))
	regionName := getLastPathElement(types.SafeString(i.Region))
	diskName := types.SafeString(i.Name)

	akas := []string{"gcp://compute.googleapis.com/projects/" + activeProject() + "/zones/" + zoneName + "/disks/" + diskName}

	if zoneName == "" {
		akas = []string{"gcp://compute.googleapis.com/projects/" + activeProject() + "/regions/" + regionName + "/disks/" + diskName}
	}

	return akas, nil
}

func diskLocation(_ context.Context, d *transform.TransformData) (interface{}, error) {
	i := d.HydrateItem.(*compute.Disk)
	param := d.Param.(string)

	zoneName := getLastPathElement(types.SafeString(i.Zone))
	regionName := getLastPathElement(types.SafeString(i.Region))

	locationData := map[string]string{
		"Type":     "ZONAL",
		"Location": zoneName,
	}

	if zoneName == "" {
		locationData["Type"] = "REGIONAL"
		locationData["Location"] = regionName
	}

	return locationData[param], nil
}
