package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/secretmanager/v1"
)

//// TABLE DEFINITION

func tableGcpSecretManagerSecret(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_secret_manager_secret",
		Description: "GCP Secret Manager Secret",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getGcpSecretManagerSecret,
		},
		List: &plugin.ListConfig{
			Hydrate: listGcpSecretManagerSecrets,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The resource name of the secret.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(lastPathElement),
			},
			{
				Name:        "create_time",
				Description: "The time at which the secret was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "etag",
				Description: "Etag of the currently stored Secret.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ttl",
				Description: "The TTL of the secret.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TTL"),
			},
			{
				Name:        "expire_time",
				Description: "The expiration time of the secret.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ExpireTime").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "replication",
				Description: "The replication policy of the secret.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "labels",
				Description: "The labels assigned to the secret.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "annotations",
				Description: "Custom metadata about the secret. Annotations are distinct from various forms of labels.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "topics",
				Description: "A list of up to 10 Pub/Sub topics to which messages are published when control plane operations are called on the secret or its versions.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "version_aliases",
				Description: "Mapping from version alias to version name.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard steampipe columns
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
				Transform:   transform.From(secretManagerSecretNameToAkas),
			},

			// Standard GCP columns
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

func listGcpSecretManagerSecrets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := SecretManagerService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_secret_manager_secret.listGcpSecretManagerSecrets", "service_error", err)
		return nil, err
	}

	// Max limit is set as per documentation
	pageSize := int64(100)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil && *limit < pageSize {
		pageSize = *limit
	}

	// Get project details
	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	resp := service.Projects.Secrets.List("projects/" + project).PageSize(pageSize)
	if err := resp.Pages(
		ctx,
		func(page *secretmanager.ListSecretsResponse) error {
			for _, secret := range page.Secrets {
				d.StreamListItem(ctx, secret)

				// Check if context has been cancelled or if the limit has been hit (if specified)
				if d.RowsRemaining(ctx) == 0 {
					return nil
				}
			}
			return nil
		},
	); err != nil {
		plugin.Logger(ctx).Error("gcp_secret_manager_secret.listGcpSecretManagerSecrets", "api_error", err)
		return nil, err
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getGcpSecretManagerSecret(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := SecretManagerService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gcp_secret_manager_secret.getGcpSecretManagerSecret", "service_error", err)
		return nil, err
	}

	// Get project details
	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	name := d.EqualsQuals["name"].GetStringValue()

	// If the name has been passed as an empty string, API does not return any error
	if len(name) < 1 {
		return nil, nil
	}
	secretName := "projects/" + project + "/secrets/" + name

	op, err := service.Projects.Secrets.Get(secretName).Do()
	if err != nil {
		plugin.Logger(ctx).Error("gcp_secret_manager_secret.getGcpSecretManagerSecret", "api_error", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func secretManagerSecretNameToAkas(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	secret := d.HydrateItem.(*secretmanager.Secret)

	// Get data for turbot defined properties
	akas := []string{"gcp://secretmanager.googleapis.com/" + secret.Name}

	return akas, nil
}
