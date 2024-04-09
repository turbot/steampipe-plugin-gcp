package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/cloudkms/v1"
)

//// TABLE DEFINITION

func tableGcpKmsKey(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_kms_key",
		Description: "GCP KMS Key",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "location", "key_ring_name"}),
			Hydrate:    getKeyDetail,
			Tags:       map[string]string{"service": "cloudkms", "action": "cryptoKeys.get"},
		},
		List: &plugin.ListConfig{
			Hydrate:       listKeyDetails,
			ParentHydrate: listKeyRingDetails,
			KeyColumns: plugin.KeyColumnSlice{
				// String columns
				{Name: "purpose", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "rotation_period", Require: plugin.Optional, Operators: []string{"<>", "="}},
			},
			Tags: map[string]string{"service": "cloudkms", "action": "cryptoKeys.list"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getKeyIamPolicy,
				Tags: map[string]string{"service": "cloudkms", "action": "cryptoKeys.getIamPolicy"},
			},
		},
		GetMatrixItemFunc: BuildLocationList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The resource name for the CryptoKey.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(lastPathElement),
			},
			{
				Name:        "key_ring_name",
				Description: "The resource name for the KeyRing.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(kmsKeyTurbotData, "KeyRing"),
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(kmsKeySelfLink),
			},
			{
				Name:        "create_time",
				Description: "The time at which this CryptoKey was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "purpose",
				Description: "The immutable purpose of this CryptoKey.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "next_rotation_time",
				Description: "At next rotation time, the Key Management Service will automatically: 1. Create a new version of this CryptoKey. 2.Mark the new version as primary.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("NextRotationTime").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "rotation_period",
				Description: "Next rotation time will be advanced by this period when the service automatically rotates a key.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RotationPeriod").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "iam_policy",
				Description: "An Identity and Access Management (IAM) policy, which specifies access controls for Google Cloud resources. A `Policy` is a collection of `bindings`. A `binding` binds one or more `members` to a single `role`. Members can be user accounts, service accounts, Google groups, and domains (such as G Suite). A `role` is a named list of permissions; each `role` can be an IAM predefined role or a user-created custom role. For some types of Google Cloud resources, a `binding` can also specify a `condition`, which is a logical expression that allows access to a resource only if the expression evaluates to `true`.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getKeyIamPolicy,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "labels",
				Description: "Labels with user-defined metadata.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "primary",
				Description: "A copy of the primary CryptoKeyVersion that will be used by Encrypt when this CryptoKey is given in EncryptRequest.name.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "version_template",
				Description: "A template describing settings for new CryptoKeyVersion instances.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(lastPathElement),
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
				Transform:   transform.FromP(kmsKeyTurbotData, "Akas"),
			},

			// GCP standard columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(kmsKeyTurbotData, "Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(kmsKeyTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTIONS

func listKeyDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listKeyDetails")

	// Create Service Connection
	service, err := KMSService(ctx, d)
	if err != nil {
		return nil, err
	}

	filterQuals := []filterQualMap{
		{"purpose", "purpose", "string"},
		{"rotation_period", "rotationPeriod", "string"},
	}

	filters := buildQueryFilterFromQuals(filterQuals, d.Quals)
	filterString := ""
	if len(filters) > 0 {
		filterString = strings.Join(filters, " ")
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

	keyRing := h.Item.(*cloudkms.KeyRing)

	resp := service.Projects.Locations.KeyRings.CryptoKeys.List(keyRing.Name).Filter(filterString).PageSize(*pageSize)
	if err := resp.Pages(ctx, func(page *cloudkms.ListCryptoKeysResponse) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, key := range page.CryptoKeys {
			d.StreamListItem(ctx, key)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
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

func getKeyDetail(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getKeyDetail")

	// Create Service Connection
	service, err := KMSService(ctx, d)
	if err != nil {
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
	ringName := d.EqualsQuals["key_ring_name"].GetStringValue()

	resp, err := service.Projects.Locations.KeyRings.CryptoKeys.Get("projects/" + project + "/locations/" + location + "/keyRings/" + ringName + "/cryptoKeys/" + name).Do()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func getKeyIamPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getKeyIamPolicy")

	// Create Service Connection
	service, err := KMSService(ctx, d)
	if err != nil {
		return nil, err
	}
	param := h.Item.(*cloudkms.CryptoKey).Name

	resp, err := service.Projects.Locations.KeyRings.CryptoKeys.GetIamPolicy(param).Do()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

//// TRANSFORM FUNCTIONS

func kmsKeyTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	key := d.HydrateItem.(*cloudkms.CryptoKey)
	param := d.Param.(string)

	project := strings.Split(key.Name, "/")[1]
	location := strings.Split(key.Name, "/")[3]
	key_ring_name := strings.Split(key.Name, "/")[5]

	turbotData := map[string]interface{}{
		"Project":  project,
		"Location": location,
		"KeyRing":  key_ring_name,
		"Akas":     []string{"gcp://cloudkms.googleapis.com/" + key.Name},
	}

	return turbotData[param], nil
}

func kmsKeySelfLink(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*cloudkms.CryptoKey)
	selfLink := "https://cloudkms.googleapis.com/v1/" + data.Name

	return selfLink, nil
}
