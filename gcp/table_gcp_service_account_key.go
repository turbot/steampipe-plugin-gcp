package gcp

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	"google.golang.org/api/iam/v1"
)

//// TABLE DEFINITION

func tableGcpServiceAccountKey(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_service_account_key",
		Description: "GCP Service Account Key",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			ItemFromKey:       serviceAccountKeyNameFromKey,
			Hydrate:           getGcpServiceAccountKey,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
		},
		List: &plugin.ListConfig{
			ParentHydrate: listGcpServiceAccounts,
			Hydrate:       listGcpServiceAccountKeys,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the service account key",
			},
			{
				Name:        "service_account_id",
				Type:        proto.ColumnType_STRING,
				Description: "Service account in which the key is located",
				Transform:   transform.FromP(getGcpServiceAccountKeyTurbotData, "ServiceAccountName"),
			},
			{
				Name:        "key_type",
				Description: "The type of the service account key",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "key_algorithm",
				Description: "Specifies the algorithm (and possibly key size) for the key",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "key_origin",
				Description: "Specifies the origin of the key",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "private_key_data",
				Description: "Specifies the private key data, which allows the assertion of the service account identity",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "private_key_type",
				Description: "Specifies the output format for the private key",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "public_key_data",
				Description: "Specifies the public key data",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "valid_after_time",
				Description: "Specifies the timestamp, after which the key can be used",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "valid_before_time",
				Description: "Specifies the timestamp, after which the key gets invalid",
				Type:        proto.ColumnType_STRING,
			},

			// standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(getGcpServiceAccountKeyTurbotData, "Title"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(getGcpServiceAccountKeyTurbotData, "Akas"),
			},

			// standard gcp columns
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromConstant(activeProject()),
			},
		},
	}
}

//// ITEM FROM KEY

func serviceAccountKeyNameFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	name := quals["name"].GetStringValue()
	item := &iam.ServiceAccountKey{
		Name: name,
	}
	return item, nil
}

//// FETCH FUNCTIONS

func listGcpServiceAccountKeys(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Fetch Service Account details
	serviceAccount := h.Item.(*iam.ServiceAccount)

	service, err := iam.NewService(ctx)
	if err != nil {
		return nil, err
	}

	result, err := service.Projects.ServiceAccounts.Keys.List(serviceAccount.Name).Do()
	for _, serviceAccountKey := range result.Keys {
		d.StreamLeafListItem(ctx, serviceAccountKey)
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getGcpServiceAccountKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	serviceAccountKey := h.Item.(*iam.ServiceAccountKey)

	service, err := iam.NewService(ctx)
	if err != nil {
		return nil, err
	}

	op, err := service.Projects.ServiceAccounts.Keys.Get(serviceAccountKey.Name).Do()
	if err != nil {
		return nil, err
	}

	return op, nil
}

/// TRANSFORM FUNCTIONS

func getGcpServiceAccountKeyTurbotData(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getGcpServiceAccountKeyTurbotData")
	serviceAccountKey := d.HydrateItem.(*iam.ServiceAccountKey)

	splitName := strings.Split(serviceAccountKey.Name, "/keys/")
	akas := []string{"gcp://iam.googleapis.com/" + serviceAccountKey.Name}

	if d.Param.(string) == "Title" {
		return splitName[1], nil
	} else if d.Param.(string) == "ServiceAccountName" {
		return splitName[0], nil
	} else {
		return akas, nil
	}
}
