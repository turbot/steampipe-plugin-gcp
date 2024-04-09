package gcp

import (
	"context"
	"strconv"
	"strings"
	"sync"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/cloudkms/v1"
)

//// TABLE DEFINITION

func tableGcpKmsKeyVersion(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_kms_key_version",
		Description: "GCP KMS Key Version",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"key_name", "key_ring_name", "location", "crypto_key_version"}),
			Hydrate:    getKeyVersionDetail,
			Tags:       map[string]string{"service": "cloudkms", "action": "cryptoKeyVersions.get"},
		},
		List: &plugin.ListConfig{
			Hydrate:       listKeyVersionDetails,
			ParentHydrate: listKeyRingDetails,
			Tags:          map[string]string{"service": "cloudkms", "action": "cryptoKeyVersions.list"},
		},
		GetMatrixItemFunc: BuildLocationList,
		Columns: []*plugin.Column{
			{
				Name:        "key_name",
				Description: "The resource name for the CryptoKeyVersion.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(kmsKeyVersionTurbotData, "KeyName"),
			},
			{
				Name:        "crypto_key_version",
				Description: "The CryptoKeyVersion of the resource.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromP(kmsKeyVersionTurbotData, "CryptoKeyVersion"),
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
			{
				Name:        "attestation",
				Description: "Statement that was generated and signed by the HSM at key creation time.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(kmsKeyVersionTurbotData, "KeyVersionName"),
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
	plugin.Logger(ctx).Trace("listKeyVersionDetails")

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
	resp := service.Projects.Locations.KeyRings.CryptoKeys.List(keyRing.Name).PageSize(*pageSize)

	var wg sync.WaitGroup
	errorCh := make(chan error, int(*pageSize))
	if err := resp.Pages(ctx, func(page *cloudkms.ListCryptoKeysResponse) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, key := range page.CryptoKeys {
			wg.Add(1)
			go getCryptoKeyVersionDetailsAsync(ctx, d, h, key, pageSize, service, errorCh, &wg)
		}
		wg.Wait()

		// NOTE: close channel before ranging over results
		close(errorCh)
		for err := range errorCh {

			// return the first error
			return err
		}
		return nil

	}); err != nil {
		return nil, err
	}
	return nil, nil
}

func getCryptoKeyVersionDetailsAsync(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, key *cloudkms.CryptoKey, pageSize *int64, service *cloudkms.Service, errorCh chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	err := getCryptoKeyVersionDetails(ctx, d, h, key, pageSize, service)
	if err != nil {
		errorCh <- err
	}
}

func getCryptoKeyVersionDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, key *cloudkms.CryptoKey, pageSize *int64, service *cloudkms.Service) error {
	resp := service.Projects.Locations.KeyRings.CryptoKeys.CryptoKeyVersions.List(key.Name).PageSize(*pageSize)

	err := resp.Pages(ctx, func(page *cloudkms.ListCryptoKeyVersionsResponse) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, keyVersion := range page.CryptoKeyVersions {
			d.StreamListItem(ctx, keyVersion)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				page.NextPageToken = ""
				return nil
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

//// HYDRATE FUNCTIONS

func getKeyVersionDetail(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getKeyVersionDetail")

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

	name := d.EqualsQuals["key_name"].GetStringValue()
	location := d.EqualsQuals["location"].GetStringValue()
	ringName := d.EqualsQuals["key_ring_name"].GetStringValue()
	version := d.EqualsQuals["crypto_key_version"].GetInt64Value()
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
	key_name := strings.Split(key.Name, "/")[7]
	crypto_key_version := strings.Split(key.Name, "/")[9]

	turbotData := map[string]interface{}{
		"Project":          project,
		"Location":         location,
		"KeyRing":          key_ring_name,
		"KeyName":          key_name,
		"KeyVersionName":   key_name + "/" + crypto_key_version,
		"CryptoKeyVersion": crypto_key_version,
		"Akas":             []string{"gcp://cloudkms.googleapis.com/" + key.Name},
	}

	return turbotData[param], nil
}

func kmsKeyVersionSelfLink(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*cloudkms.CryptoKeyVersion)
	selfLink := "https://cloudkms.googleapis.com/v1/" + data.Name

	return selfLink, nil
}
