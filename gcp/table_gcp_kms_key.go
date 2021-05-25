package gcp

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	"google.golang.org/api/cloudkms/v1"
)

func tableGcpKmsKey(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_kms_key",
		Description: "GCP Kms Key",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getKeyDetail,
		},
		List: &plugin.ListConfig{
			Hydrate:       listKeyDetails,
			ParentHydrate: listKeyRingDetails,
		},
		GetMatrixItem: BuildLocationList,
		Columns: []*plugin.Column{
			// commonly used columns
			{
				Name:        "name",
				Description: "The resource name for this CryptoKey.",
				Type:        proto.ColumnType_STRING,
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
				Description: "At next_rotation_time, the Key Management Service will automatically: 1. Create a new version of this CryptoKey. 2.Mark the new version as primary.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "rotation_period",
				Description: "RotationPeriod: next_rotation_time will be advanced by this period when the service automatically rotates a key.",
				Type:        proto.ColumnType_STRING,
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

			// standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpKmsKeyTurbotData, "Title"),
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
				Transform:   transform.FromP(gcpKmsKeyTurbotData, "Akas"),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpKmsKeyTurbotData, "Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpKmsKeyTurbotData, "Project"),
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
	keyRing := h.Item.(*cloudkms.KeyRing)
	resp := service.Projects.Locations.KeyRings.CryptoKeys.List(keyRing.Name)
	if err := resp.Pages(ctx, func(page *cloudkms.ListCryptoKeysResponse) error {
		for _, key := range page.CryptoKeys {
			d.StreamListItem(ctx, key)
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
	name := d.KeyColumnQuals["name"].GetStringValue()
	resp, err := service.Projects.Locations.KeyRings.CryptoKeys.Get(name).Do()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

//// TRANSFORM FUNCTIONS

func gcpKmsKeyTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	key := d.HydrateItem.(*cloudkms.CryptoKey)
	param := d.Param.(string)

	project := strings.Split(key.Name, "/")[1]
	location := strings.Split(key.Name, "/")[3]
	title := strings.Split(key.Name, "/")[5]

	turbotData := map[string]interface{}{
		"Project":  project,
		"Location": location,
		"Title":    title,
		"Akas":     []string{"gcp://cloudkms.googleapis.com/" + key.Name},
	}

	return turbotData[param], nil
}
