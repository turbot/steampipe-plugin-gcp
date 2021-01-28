package gcp

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/storage/v1"
)

func tableGcpStorageBucket(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_storage_bucket",
		Description: "GCP Storage Bucket",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getGcpStorageBucket,
		},
		List: &plugin.ListConfig{
			Hydrate: listGcpStorageBuckets,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the bucket.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The ID of the bucket. For buckets, the id and name properties are the same.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "The kind of item this is. For buckets, this is always storage#bucket.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The creation time of the bucket in RFC 3339 format.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "location_type",
				Description: "The type of the bucket location.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "storage_class",
				Description: "The bucket's default storage class, used whenever no storageClass is specified for a newly-created object. This defines how objects in the bucket are stored and determines the SLA and the cost of storage. Values include MULTI_REGIONAL, REGIONAL, STANDARD, NEARLINE, COLDLINE, ARCHIVE, and DURABLE_REDUCED_AVAILABILITY. If this value is not specified when the bucket is created, it will default to STANDARD.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "billing_requester_pays",
				Description: "When set to true, Requester Pays is enabled for this bucket.",
				Default:     false,
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Billing.RequesterPays"),
			},
			{
				Name:        "default_event_based_hold",
				Description: "The default value for event-based hold on newly created objects in this bucket. Event-based hold is a way to retain objects indefinitely until an event occurs, signified by the hold's release. After being released, such objects will be subject to bucket-level retention (if any).",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "default_kms_key_name",
				Description: "A Cloud KMS key that will be used to encrypt objects inserted into this bucket, if no encryption method is specified.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Encryption.DefaultKmsKeyName"),
			},

			{
				Name:        "etag",
				Description: "HTTP 1.1 Entity tag for the bucket.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "iam_configuration_bucket_policy_only_enabled",
				Description: "The bucket's uniform bucket-level access configuration. The feature was formerly known as Bucket Policy Only. For backward compatibility, this field will be populated with identical information as the uniformBucketLevelAccess field.",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
				Transform:   transform.FromField("IamConfiguration.BucketPolicyOnly.Enabled"),
			},
			{
				Name:        "iam_configuration_public_access_prevention",
				Description: "The bucket's Public Access Prevention configuration. Currently, 'unspecified' and 'enforced' are supported.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("IamConfiguration.PublicAccessPrevention"),
			},
			{
				Name:        "iam_configuration_uniform_bucket_level_access_enabled",
				Description: "The bucket's uniform bucket-level access configuration.",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
				Transform:   transform.FromField("IamConfiguration.UniformBucketLevelAccess.Enabled"),
			},
			{
				Name:        "labels",
				Description: "Labels that apply to this bucket.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "log_bucket",
				Description: "The destination bucket where the current bucket's logs should be placed.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Logging.LogBucket"),
			},
			{
				Name:        "log_object_prefix",
				Description: "A prefix for log object names.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Logging.LogObjectPrefix"),
			},
			{
				Name:        "metageneration",
				Description: "The metadata generation of this bucket.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "owner_entity",
				Description: "The entity, in the form project-owner-projectId. This is always the project team's owner group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Owner.Entity"),
			},
			{
				Name:        "owner_entity_id",
				Description: "The ID for the entity.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Owner.EntityId"),
			},
			{
				Name:        "project_number",
				Description: "The project number of the project the bucket belongs to.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "versioning_enabled",
				Description: "While set to true, versioning is fully enabled for this bucket.",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
				Transform:   transform.FromField("Versioning.Enabled"),
			},
			{
				Name:        "website_main_page_suffix",
				Description: "If the requested object path is missing, the service will ensure the path has a trailing '/', append this suffix, and attempt to retrieve the resulting object. This allows the creation of index.html objects to represent directory pages.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Website.MainPageSuffix"),
			},
			{
				Name:        "website_not_found_page",
				Description: "If the requested object path is missing, and any mainPageSuffix object is missing, if applicable, the service will return the named object from this bucket as the content for a 404 Not Found result.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Website.NotFoundPage"),
			},
			{
				Name:        "self_link",
				Description: "The URI of this bucket.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "updated",
				Description: "The modification time of the bucket.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "acl",
				Description: "An access-control list",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGcpStorageBucketACLs,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "default_object_acl",
				Description: "Lists of object access control entries",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGcpStorageBucketDefaultACLs,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "cors",
				Description: "The bucket's Cross-Origin Resource Sharing (CORS) configuration.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "iam_policy",
				Description: "An Identity and Access Management (IAM) policy, which specifies access controls for Google Cloud resources. A `Policy` is a collection of `bindings`. A `binding` binds one or more `members` to a single `role`. Members can be user accounts, service accounts, Google groups, and domains (such as G Suite). A `role` is a named list of permissions; each `role` can be an IAM predefined role or a user-created custom role. For some types of Google Cloud resources, a `binding` can also specify a `condition`, which is a logical expression that allows access to a resource only if the expression evaluates to `true`.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGcpStorageBucketIAMPolicy,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "lifecycle_rules",
				Description: "The bucket's lifecycle configuration. See lifecycle management for more information.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Lifecycle.Rule"),
			},
			{
				Name:        "retention_policy",
				Description: "The bucket's retention policy. The retention policy enforces a minimum retention time for all objects contained in the bucket, based on their creation time. Any attempt to overwrite or delete objects younger than the retention period will result in a PERMISSION_DENIED error.",
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
				Transform:   transform.From(bucketAka),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
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

//// LIST FUNCTION

func listGcpStorageBuckets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	project := activeProject()

	service, err := storage.NewService(ctx)
	if err != nil {
		return nil, err
	}

	resp := service.Buckets.List(project)
	if err := resp.Pages(ctx, func(page *storage.Buckets) error {
		for _, bucket := range page.Items {
			d.StreamListItem(ctx, bucket)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, err
}

func getGcpStorageBucket(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// project := activeProject()
	name := d.KeyColumnQuals["name"].GetStringValue()

	service, err := storage.NewService(ctx)
	if err != nil {
		return nil, err
	}

	req, err := service.Buckets.Get(name).Do()
	if err != nil {
		plugin.Logger(ctx).Trace("getGcpStorageBucket", "Error", err)
		return nil, err
	}

	return req, nil
}

//// HYDRATE FUNCTIONS

func getGcpStorageBucketIAMPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getGcpStorageBucketIAMPolicy")
	bucket := h.Item.(*storage.Bucket)

	// Create Session
	service, err := storage.NewService(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := service.Buckets.GetIamPolicy(bucket.Name).Do()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func getGcpStorageBucketACLs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getGcpStorageBucketACLs")
	bucket := h.Item.(*storage.Bucket)

	// Create Session
	service, err := storage.NewService(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := service.BucketAccessControls.List(bucket.Name).Do()
	if err != nil {
		gerr, _ := err.(*googleapi.Error)

		// It should not error out if the bucket has uniform bucket-level accesss
		// googleapi: Error 400: Cannot get legacy ACL for a bucket that has uniform bucket-level access.
		// Read more at https://cloud.google.com/storage/docs/uniform-bucket-level-access, invalid
		if gerr.Code == 400 && strings.HasPrefix(gerr.Message, "Cannot get legacy ACL for a bucket that has uniform bucket-level access") {
			return nil, nil
		}
		return nil, err
	}

	return resp, nil
}

func getGcpStorageBucketDefaultACLs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getGcpStorageBucketDefaultACLs")
	bucket := h.Item.(*storage.Bucket)

	// Create Session
	service, err := storage.NewService(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := service.DefaultObjectAccessControls.List(bucket.Name).Do()
	if err != nil {
		gerr, _ := err.(*googleapi.Error)

		// It should not error out if the bucket has uniform bucket-level accesss
		// googleapi: Error 400: Cannot get legacy ACL for a bucket that has uniform bucket-level access.
		// Read more at https://cloud.google.com/storage/docs/uniform-bucket-level-access, invalid
		if gerr.Code == 400 && strings.HasPrefix(gerr.Message, "Cannot get legacy ACL for a bucket that has uniform bucket-level access") {
			return nil, nil
		}
		return nil, err
	}

	return resp, nil
}

//// TRANSFORM FUNCTIONS

func bucketAka(_ context.Context, d *transform.TransformData) (interface{}, error) {
	bucket := d.HydrateItem.(*storage.Bucket)
	project := activeProject()

	akas := []string{"gcp://storage.googleapis.com/projects/" + project + "/buckets/" + bucket.Name}
	return akas, nil
}
