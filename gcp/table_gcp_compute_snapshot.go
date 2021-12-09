package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	"google.golang.org/api/compute/v1"
)

func tableGcpComputeSnapshot(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_snapshot",
		Description: "GCP Compute Snapshot",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeSnapshot,
		},
		List: &plugin.ListConfig{
			Hydrate:           listComputeSnapshots,
			ShouldIgnoreError: isIgnorableError([]string{"403"}),
		},
		Columns: []*plugin.Column{
			// commonly used columns
			{
				Name:        "name",
				Description: "Name of the resource; provided by the client when the resource is created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_disk",
				Description: "The url of the source disk used to create this snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			// source_disk_name is a simpler view of the type, without the full path
			{
				Name:        "source_disk_name",
				Description: "The name of the source disk used to create this snapshot.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SourceDisk").Transform(lastPathElement),
			},
			{
				Name:        "description",
				Description: "An optional description of this resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "auto_created",
				Description: "Set to true if snapshots are automatically created by applying resource policy on the target disk.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "creation_timestamp",
				Description: "Timestamp when snapshot was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "disk_size_gb",
				Description: "Size of the source disk, specified in GB.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "download_bytes",
				Description: "Number of bytes downloaded to restore a snapshot to a disk.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "label_fingerprint",
				Description: "A fingerprint for the labels being applied to this snapshot, which is essentially a hash of the labels set used for optimistic locking. The fingerprint is initially generated by Compute Engine and changes after every request to modify or update labels.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the snapshot. This can be CREATING, DELETING, FAILED, READY, or UPLOADING.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "storage_bytes",
				Description: "A size of the storage used by the snapshot. As snapshots share storage, this number is expected to change with snapshot creation/deletion.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "storage_bytes_status",
				Description: "An indicator whether storageBytes is in a stable state or it is being adjusted as a result of shared storage reallocation. This status can either be UPDATING, meaning the size of the snapshot is being updated, or UP_TO_DATE, meaning the size of the snapshot is up-to-date.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kms_key_name",
				Description: "The name of the encryption key that is used to encrypt snapshot",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SnapshotEncryptionKey.KmsKeyName"),
			},
			{
				Name:        "kms_key_service_account",
				Description: "The service account being used for the encryption request for the given KMS key. If absent, the Compute Engine default service account is used.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SnapshotEncryptionKey.KmsKeyServiceAccount"),
			},
			{
				Name:        "encryption_key_raw_key",
				Description: "Specifies a 256-bit customer-supplied encryption key, encoded in RFC 4648 base64 to either encrypt or decrypt this resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SnapshotEncryptionKey.RawKey"),
			},
			{
				Name:        "encryption_key_sha256",
				Description: "The RFC 4648 base64 encoded SHA-256 hash of the customer-supplied encryption key that protects this resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SnapshotEncryptionKey.Sha256"),
			},
			{
				Name:        "source_disk_encryption_key",
				Description: "The customer-supplied encryption key of the source disk.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SourceDiskEncryptionKey"),
			},
			{
				Name:        "labels",
				Description: "Labels applied to this snapshot.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Labels"),
			},
			{
				Name:        "licenses",
				Description: "A list of public visible licenses that apply to this snapshot. This can be because the original image had licenses attached (such as a Windows image).",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "storage_locations",
				Description: "Cloud Storage bucket storage location of the snapshot (regional or multi-regional).",
				Type:        proto.ColumnType_JSON,
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
				Transform:   transform.FromP(gcpComputeSnapshotTurbotData, "Akas"),
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
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpComputeSnapshotTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeSnapshots(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listComputeSnapshots")

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

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

	resp := service.Snapshots.List(project).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *compute.SnapshotList) error {
		for _, snapshot := range page.Items {
			d.StreamListItem(ctx, snapshot)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
			page.NextPageToken = ""
				return nil
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getComputeSnapshot(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getComputeSnapshot")

	// Create Service Connection
	service, err := ComputeService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details
	getProjectCached := plugin.HydrateFunc(getProject).WithCache()
	projectId, err := getProjectCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	name := d.KeyColumnQuals["name"].GetStringValue()

	// Error: pq: rpc error: code = Unknown desc = json: invalid use of ,string struct tag,
	// trying to unmarshal "projects/project/global/snapshots/" into uint64
	if len(name) < 1 {
		return nil, nil
	}

	resp, err := service.Snapshots.Get(project, name).Do()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

//// TRANSFORM FUNCTIONS

func gcpComputeSnapshotTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	snapshot := d.HydrateItem.(*compute.Snapshot)
	param := d.Param.(string)

	project := strings.Split(snapshot.SelfLink, "/")[6]

	turbotData := map[string]interface{}{
		"Project": project,
		"Akas":    []string{"gcp://compute.googleapis.com/projects/" + project + "/global/snapshots/" + snapshot.Name},
	}

	return turbotData[param], nil
}
