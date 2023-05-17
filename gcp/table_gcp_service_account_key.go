package gcp

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iam/v1"
)

//// TABLE DEFINITION

func tableGcpServiceAccountKey(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_service_account_key",
		Description: "GCP Service Account Key",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "service_account_name"}),
			Hydrate:    getGcpServiceAccountKey,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listGcpServiceAccounts,
			Hydrate:       listGcpServiceAccountKeys,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the service account key.",
				Transform:   transform.FromField("Name").Transform(lastPathElement),
			},
			{
				Name:        "service_account_name",
				Type:        proto.ColumnType_STRING,
				Description: "Service account in which the key is located.",
				Transform:   transform.FromP(getGcpServiceAccountKeyTurbotData, "ServiceAccountName"),
			},
			{
				Name:        "key_type",
				Description: "The type of the service account key.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "key_algorithm",
				Description: "Specifies the algorithm (and possibly key size) for the key.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "key_origin",
				Description: "Specifies the origin of the key.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "public_key_data",
				Description: "Specifies the public key data.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getGcpServiceAccountKey,
			},
			{
				Name:        "valid_after_time",
				Description: "Specifies the timestamp, after which the key can be used.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "valid_before_time",
				Description: "Specifies the timestamp, after which the key gets invalid.",
				Type:        proto.ColumnType_TIMESTAMP,
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
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromConstant("global"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Hydrate:     plugin.HydrateFunc(getProject).WithCache(),
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// FETCH FUNCTIONS

func listGcpServiceAccountKeys(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Fetch Service Account details
	serviceAccount := h.Item.(*iam.ServiceAccount)

	// Create Service Connection
	service, err := IAMService(ctx, d)
	if err != nil {
		return nil, err
	}

	result, err := service.Projects.ServiceAccounts.Keys.List(serviceAccount.Name).Do()
	if err != nil {
		return nil, err
	}
	for _, serviceAccountKey := range result.Keys {
		d.StreamLeafListItem(ctx, serviceAccountKey)
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getGcpServiceAccountKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getGcpServiceAccountKey")

	var name, serviceAccountName string
	if h.Item != nil {
		data := h.Item.(*iam.ServiceAccountKey)
		name = strings.Split(data.Name, "/")[5]
		splitName := strings.Split(data.Name, "/")
		serviceAccountName = splitName[3]
	} else {
		name = d.EqualsQuals["name"].GetStringValue()
		serviceAccountName = d.EqualsQuals["service_account_name"].GetStringValue()
	}

	// Empty check for the input param
	if name == "" || serviceAccountName == "" {
		return nil, nil
	}

	// Create Service Connection
	service, err := IAMService(ctx, d)
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

	keyName := "projects/" + project + "/serviceAccounts/" + serviceAccountName + "/keys/" + name

	// We should pass the supported public key output formats as an query parameter to get the public key data for a service account key as mentioned in the doc https://cloud.google.com/iam/docs/reference/rest/v1/projects.serviceAccounts.keys/get#ServiceAccountPublicKeyType,
	// we can pass three value for the query parameter publicKeyType that are TYPE_RAW_PUBLIC_KEY or TYPE_X509_PEM_FILE or TYPE_NONE.
	// TYPE_NONE has excluded here from the query parameter value because if we pass that then it do not return the public key.
	// We are getting the public key data in raw format here.
	queryParameter := googleapi.QueryParameter("publicKeyType", "TYPE_RAW_PUBLIC_KEY")

	op, err := service.Projects.ServiceAccounts.Keys.Get(keyName).Do(queryParameter)
	if err != nil {
		return nil, err
	}

	return op, nil
}

/// TRANSFORM FUNCTIONS

func getGcpServiceAccountKeyTurbotData(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getGcpServiceAccountKeyTurbotData")
	serviceAccountKey := d.HydrateItem.(*iam.ServiceAccountKey)

	splitName := strings.Split(serviceAccountKey.Name, "/")
	akas := []string{"gcp://iam.googleapis.com/" + serviceAccountKey.Name}

	if d.Param.(string) == "Title" {
		return splitName[5], nil
	} else if d.Param.(string) == "ServiceAccountName" {
		return splitName[3], nil
	} else {
		return akas, nil
	}
}
