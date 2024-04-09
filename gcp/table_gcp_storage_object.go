package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/storage/v1"
)

func tableGcpStorageObject(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_storage_object",
		Description: "GCP Storage Object",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"bucket", "name"}),
			Hydrate:    getStorageObject,
			Tags:       map[string]string{"service": "storage", "action": "objects.get"},
		},
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "bucket", Require: plugin.Required, CacheMatch: "exact"},
				{Name: "prefix", Require: plugin.Optional},
			},
			Hydrate: listStorageObjects,
			Tags:    map[string]string{"service": "storage", "action": "objects.list"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getStorageObjectIAMPolicy,
				Tags: map[string]string{"service": "storage", "action": "objects.getIamPolicy"},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the object.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "prefix",
				Description: "The prefix of the key of the object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("prefix"),
			},
			{
				Name:        "id",
				Description: "The ID of the object, including the bucket name, object name, and generation number.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "bucket",
				Description: "The name of the bucket containing this object.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cache_control",
				Description: "Cache-Control directive for the object data.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "component_count",
				Description: "Number of underlying components that make up this object.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "content_disposition",
				Description: "Content-Disposition of the object data.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "content_encoding",
				Description: "Content-Encoding of the object data.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "content_language",
				Description: "Content-Language of the object data.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "content_type",
				Description: "Content-Type of the object data.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "crc32c",
				Description: "CRC32c checksum, as described in RFC 4960, Appendix B; encoded using base64 in big-endian byte order.",
				Transform:   transform.FromField("Crc32c"),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "custom_time",
				Description: "A timestamp in RFC 3339 format specified by the user for an object",
				Transform:   transform.FromGo().NullIfZero(),
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "etag",
				Description: "Entity tag for the object",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "event_based_hold",
				Description: "Whether or not the object is under event-based hold.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "generation",
				Description: "The content generation of this object.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "kind",
				Description: "The kind of item this is.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kms_key_name",
				Description: "Cloud KMS Key used to encrypt this object, if the object is encrypted by such a key.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "md5_hash",
				Description: "MD5 hash of the data; encoded using base64",
				Transform:   transform.FromField("Md5Hash"),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "media_link",
				Description: "Media download link",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "metadata",
				Description: "User-provided metadata, in key/value pairs.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "metageneration",
				Description: "The version of the metadata for this object at this generation.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "retention_expiration_time",
				Description: "A server-determined value that specifies the earliest time that the object's retention period expires.",
				Transform:   transform.FromGo().NullIfZero(),
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "self_link",
				Description: "The link to this object.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "size",
				Description: "Content-Length of the data in bytes.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "storage_class",
				Description: "Storage class of the object.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "temporary_hold",
				Description: "Whether or not the object is under temporary hold.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "time_created",
				Description: "The creation time of the object.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "time_deleted",
				Description: "The deletion time of the object.",
				Transform:   transform.FromGo().NullIfZero(),
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "time_storage_class_updated",
				Description: "The time at which the object's storage class was last changed.",
				Transform:   transform.FromGo().NullIfZero(),
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "updated",
				Description: "The modification time of the object metadata.",
				Transform:   transform.FromGo().NullIfZero(),
				Type:        proto.ColumnType_TIMESTAMP,
			},

			// JSON fields
			{
				Name:        "acl",
				Description: "Access controls on the object.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "customer_encryption",
				Description: "Metadata of customer-supplied encryption key, if the object is encrypted by such a key",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "owner",
				Description: "The owner of the object. This will always be the uploader of the object.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "iam_policy",
				Description: "An Identity and Access Management (IAM) policy, which specifies access controls for Google Cloud resources. A `Policy` is a collection of `bindings`. A `binding` binds one or more `members` to a single `role`. Members can be user accounts, service accounts, Google groups, and domains (such as G Suite). A `role` is a named list of permissions; each `role` can be an IAM predefined role or a user-created custom role. For some types of Google Cloud resources, a `binding` can also specify a `condition`, which is a logical expression that allows access to a resource only if the expression evaluates to `true`.",
				Hydrate:     getStorageObjectIAMPolicy,
				Transform:   transform.FromValue(),
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
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getObjectAka,
				Transform:   transform.FromValue(),
			},

			// standard GCP columns
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

func listStorageObjects(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	bucket := d.EqualsQualString("bucket")
	prefix := d.EqualsQualString("prefix")

	// The bucket name should not be empty
	if bucket == "" {
		return nil, nil
	}

	service, err := StorageService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Trace("gcp_storage_object.listStorageObjects", "connection_error", err)
		return nil, err
	}

	// Max limit isn't mentioned in the documentation
	// Default limit is set as 1000
	maxResults := types.Int64(1000)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *maxResults {
			maxResults = limit
		}
	}

	resp := service.Objects.List(bucket).Prefix(prefix).Projection("full").MaxResults(*maxResults)
	if err := resp.Pages(ctx, func(page *storage.Objects) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, object := range page.Items {
			d.StreamListItem(ctx, object)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				page.NextPageToken = ""
				break
			}
		}
		return nil
	}); err != nil {
		plugin.Logger(ctx).Trace("gcp_storage_object.listStorageObjects", "api_error", err)
		return nil, err
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getStorageObject(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	bucket := d.EqualsQuals["bucket"].GetStringValue()
	name := d.EqualsQuals["name"].GetStringValue()

	// Return nil, if input parameters are empty
	if bucket == "" || name == "" {
		return nil, nil
	}

	service, err := StorageService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Trace("gcp_storage_object.getStorageObject", "connection_error", err)
		return nil, err
	}

	req, err := service.Objects.Get(bucket, name).Do()
	if err != nil {
		plugin.Logger(ctx).Trace("gcp_storage_object.getStorageObject", "api_error", err)
		return nil, err
	}

	return req, nil
}

func getStorageObjectIAMPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	object := h.Item.(*storage.Object)

	// Create Session
	service, err := StorageService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Trace("gcp_storage_object.getStorageObjectIAMPolicy", "connection_error", err)
		return nil, err
	}

	resp, err := service.Objects.GetIamPolicy(object.Bucket, object.Name).Do()
	if err != nil {

		// Return nil, if uniform bucket-level access is enabled
		if strings.Contains(err.(*googleapi.Error).Message, "Object policies are disabled for bucket") {
			return nil, nil
		}
		plugin.Logger(ctx).Trace("gcp_storage_object.getStorageObjectIAMPolicy", "api_error", err)
		return nil, err
	}

	return resp, nil
}

func getObjectAka(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	object := h.Item.(*storage.Object)

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Trace("gcp_storage_object.getObjectAka", "cache_error", err)
		return nil, err
	}
	project := projectId.(string)

	akas := []string{"gcp://storage.googleapis.com/projects/" + project + "/buckets/" + object.Bucket + "/objects/" + object.Name}
	return akas, nil
}
