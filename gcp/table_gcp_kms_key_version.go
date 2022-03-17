package gcp

import (
	"context"
	"strconv"
	"strings"
	"sync"

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
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "key_ring_name", "location", "crypto_key_versions"}),
			Hydrate:    getKeyVersionDetail,
		},
		List: &plugin.ListConfig{
			Hydrate:           listKeyVersionDetails,
			ParentHydrate:     listKeyRingDetails,
			ShouldIgnoreError: isIgnorableError([]string{"403"}),
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
				Type:        proto.ColumnType_INT,
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
				Name:        "create_time",
				Description: "The time at which this CryptoKeyVersion was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "destroy_event_time",
				Description: "The time this CryptoKeyVersion's key material was destroyed.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("DestroyEventTime").NullIfZero(),
			},
			{
				Name:        "destroy_time",
				Description: "The time this CryptoKeyVersion's key material is scheduled for destruction.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("DestroyTime").NullIfZero(),
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
	logger := plugin.Logger(ctx)
	logger.Trace("listKeyVersionDetails")

	// Create Service Connection
	service, err := KMSService(ctx, d)
	if err != nil {
		return nil, err
	}

	pageSize := types.Int64(1000)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *pageSize {
			pageSize = limit
		}
	}

	cryptoKeys, err := listKeyDetailsX(ctx, d, h)

	chunkCryptoKeys := chunkSlice(cryptoKeys, 50)
	var wg sync.WaitGroup
	wg.Add(len(chunkCryptoKeys))
	for i := 1; i <= len(chunkCryptoKeys); i++ {
		cryptoKeysSlice := chunkCryptoKeys[i-1]

		for _, key := range cryptoKeysSlice {
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
		}
		wg.Done()

	}
	wg.Wait()

	return nil, nil
}

func listKeyDetailsX(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) ([]*cloudkms.CryptoKey, error) {

	service, err := KMSService(ctx, d)
	if err != nil {
		return nil, err
	}

	pageSize := types.Int64(1000)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *pageSize {
			pageSize = limit
		}
	}

	keyRing := h.Item.(*cloudkms.KeyRing)

	var cryptoKeys []*cloudkms.CryptoKey
	respKey := service.Projects.Locations.KeyRings.CryptoKeys.List(keyRing.Name).PageSize(*pageSize)
	if errKey := respKey.Pages(ctx, func(page *cloudkms.ListCryptoKeysResponse) error {
		for _, key := range page.CryptoKeys {

			cryptoKeys = append(cryptoKeys, key)

		}
		return nil
	}); errKey != nil {
		return nil, errKey
	}
	return cryptoKeys, nil
}

//// HYDRATE FUNCTIONS

func getKeyVersionDetail(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getKeyVersionDetail")

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
	version := d.KeyColumnQuals["crypto_key_version"].GetInt64Value()
	logger.Trace("projects/" + project + "/locations/" + location + "/keyRings/" + ringName + "/cryptoKeys/" + name + "/cryptoKeyVersions/" + strconv.FormatInt(version, 10))
	resp, err := service.Projects.Locations.KeyRings.CryptoKeys.CryptoKeyVersions.
		Get("projects/" + project + "/locations/" + location + "/keyRings/" + ringName + "/cryptoKeys/" + name + "/cryptoKeyVersions/" + strconv.FormatInt(version, 10)).Do()
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

func chunkSlice(slice []*cloudkms.CryptoKey, chunkSize int) [][]*cloudkms.CryptoKey {
	var chunks [][]*cloudkms.CryptoKey
	for {
		if len(slice) == 0 {
			break
		}

		// necessary check to avoid slicing beyond
		// slice capacity
		if len(slice) < chunkSize {
			chunkSize = len(slice)
		}

		chunks = append(chunks, slice[0:chunkSize])
		slice = slice[chunkSize:]
	}

	return chunks
}
