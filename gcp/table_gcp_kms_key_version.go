package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	"google.golang.org/api/cloudkms/v1"
)

//// TABLE DEFINITION

func tableGcpKmsKeyVersion(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_kms_key_version",
		Description: "GCP KMS Key Version",
		// Get: &plugin.GetConfig{
		// 	KeyColumns: plugin.AllColumns([]string{"name", "location", "key_ring_name", "crypto_key_versions"}),
		// 	Hydrate:    getKeyVersionDetail,
		// },
		List: &plugin.ListConfig{
			Hydrate:           listKeyVersionDetails,
			ParentHydrate:     listKeyDetails,
			ShouldIgnoreError: isIgnorableError([]string{"403"}),
			// KeyColumns: plugin.KeyColumnSlice{
			// 	// String columns
			// 	{Name: "purpose", Require: plugin.Optional, Operators: []string{"<>", "="}},
			// 	{Name: "rotation_period", Require: plugin.Optional, Operators: []string{"<>", "="}},
			// },
		},
		GetMatrixItem: BuildLocationList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The resource name for the CryptoKeyVersion.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(kmsKeyVersionTurbotData, "ResourceName"),
			},
			{
				Name:        "crypto_key_versions",
				Description: "The CryptoKeyVersion of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(kmsKeyVersionTurbotData, "CryptoKeyVersions"),
			},
			{
				Name:        "key_ring_name",
				Description: "The resource name for the KeyRing.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(kmsKeyVersionTurbotData, "KeyRing"),
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(kmsKeyVersionSelfLink),
			},
			{
				Name:        "state",
				Description: "The current state of the CryptoKeyVersion.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "The time at which this CryptoKeyVersion was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "algorithm",
				Description: "The CryptoKeyVersionAlgorithm that this CryptoKeyVersion supports.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "attestation",
				Description: "Statement that was generated and signed by the HSM at key creation time.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "destroy_event_time",
				Description: "The time this CryptoKeyVersion's key material was destroyed.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "destroy_time",
				Description: "The time this CryptoKeyVersion's key material is scheduled for destruction.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "external_key_uri",
				Description: "The URI for an external resource that this CryptoKeyVersion represents.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ExternalProtectionLevelOptions.ExternalKeyUri"),
			},
			{
				Name:        "generate_time",
				Description: "The time this CryptoKeyVersion's key material was generated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "import_failure_reason",
				Description: "The root cause of an import failure.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "import_job",
				Description: "The name of the ImportJob used to import this CryptoKeyVersion.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "protection_level",
				Description: "The ProtectionLevel describing how crypto operations are performed with this CryptoKeyVersion.",
				Type:        proto.ColumnType_STRING,
			},
			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(kmsKeyVersionTurbotData, "ResourceName"),
			},
			// {
			// 	Name:        "tags",
			// 	Description: ColumnDescriptionTags,
			// 	Type:        proto.ColumnType_JSON,
			// 	Transform:   transform.FromField("Labels"),
			// },
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(kmsKeyVersionTurbotData, "Akas"),
			},

			// GCP standard columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(kmsKeyVersionTurbotData, "Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(kmsKeyVersionTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTIONS

func listKeyVersionDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listKeyVersionDetails")

	// Create Service Connection
	service, err := KMSService(ctx, d)
	if err != nil {
		return nil, err
	}

	// filterQuals := []filterQualMap{
	// 	{"purpose", "purpose", "string"},
	// 	{"rotation_period", "rotationPeriod", "string"},
	// }

	// filters := buildQueryFilterFromQuals(filterQuals, d.Quals)
	// filterString := ""
	// if len(filters) > 0 {
	// 	filterString = strings.Join(filters, " ")
	// }

	// Max limit isn't mentioned in the documentation
	// Default limit is set as 1000
	pageSize := types.Int64(1000)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *pageSize {
			pageSize = limit
		}
	}

	key := h.Item.(*cloudkms.CryptoKey)

	resp := service.Projects.Locations.KeyRings.CryptoKeys.CryptoKeyVersions.List(key.Name).PageSize(*pageSize)
	// Filter(filterString).

	if err := resp.Pages(ctx, func(page *cloudkms.ListCryptoKeyVersionsResponse) error {
		for _, key := range page.CryptoKeyVersions {
			d.StreamListItem(ctx, key)

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

/* func getKeyVersionDetail(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getKeyVersionDetail")

	// Create Service Connection
	service, err := KMSService(ctx, d)
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
	location := d.KeyColumnQuals["location"].GetStringValue()
	ringName := d.KeyColumnQuals["key_ring_name"].GetStringValue()
	version := d.KeyColumnQuals["crypto_key_versions"].GetInt64Value()

	resp, err := service.Projects.Locations.KeyRings.CryptoKeys.CryptoKeyVersions.
		Get("projects/" + project + "/locations/" + location + "/keyRings/" + ringName + "/cryptoKeys/" + name + "/cryptoKeyVersions/" + strconv.FormatInt(version, 10)).Do()
	if err != nil {
		return nil, err
	}

	return resp, nil
} */

func getKeyVersionPublicKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getKeyVersionPublicKey")

	// Create Service Connection
	service, err := KMSService(ctx, d)
	if err != nil {
		return nil, err
	}
	param := h.Item.(*cloudkms.CryptoKeyVersion).Name

	resp, err := service.Projects.Locations.KeyRings.CryptoKeys.CryptoKeyVersions.GetPublicKey(param).Do()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

//// TRANSFORM FUNCTIONS

func kmsKeyVersionTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	key := d.HydrateItem.(*cloudkms.CryptoKeyVersion)
	param := d.Param.(string)

	project := strings.Split(key.Name, "/")[1]
	location := strings.Split(key.Name, "/")[3]
	key_ring_name := strings.Split(key.Name, "/")[5]
	resource_name := strings.Split(key.Name, "/")[7]
	crypto_key_versions := strings.Split(key.Name, "/")[9]

	turbotData := map[string]interface{}{
		"Project":           project,
		"Location":          location,
		"KeyRing":           key_ring_name,
		"ResourceName":      resource_name,
		"CryptoKeyVersions": crypto_key_versions,
		"Akas":              []string{"gcp://cloudkms.googleapis.com/" + key.Name},
	}

	return turbotData[param], nil
}

func kmsKeyVersionSelfLink(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*cloudkms.CryptoKeyVersion)
	selfLink := "https://cloudkms.googleapis.com/v1/" + data.Name

	return selfLink, nil
}
