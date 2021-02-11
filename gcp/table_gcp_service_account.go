package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"google.golang.org/api/iam/v1"
)

//// TABLE DEFINITION

func tableGcpServiceAccount(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_service_account",
		Description: "GCP Service Account",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			Hydrate:           getGcpServiceAccount,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listGcpServiceAccounts,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The resource name of the service account",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "email",
				Description: "The email address of the service account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "unique_id",
				Description: "The unique, stable numeric ID for the service account. Each service account retains its unique ID even if you delete the service account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "A user-specified, human-readable name for the service account. The maximum length is 100 UTF-8 bytes. Optional",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A user-specified, human-readable description of the service account. The maximum length is 256 UTF-8 bytes.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "iam_policy",
				Description: "An Identity and Access Management (IAM) policy, which specifies access controls for Google Cloud resources. A `Policy` is a collection of `bindings`. A `binding` binds one or more `members` to a single `role`. Members can be user accounts, service accounts, Google groups, and domains (such as G Suite). A `role` is a named list of permissions; each `role` can be an IAM predefined role or a user-created custom role. For some types of Google Cloud resources, a `binding` can also specify a `condition`, which is a logical expression that allows access to a resource only if the expression evaluates to `true`.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getServiceAccountIamPolicy,
				Transform:   transform.FromValue(),
			},

			// standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Default:     transform.FromField("Name"),
				Transform:   transform.FromField("DisplayName"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(serviceAccountNameToAkas),
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
				Transform:   transform.FromConstant(projectName),
			},
		},
	}
}

//// LIST FUNCTION

func listGcpServiceAccounts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	service, err := iam.NewService(ctx)
	if err != nil {
		return nil, err
	}

	project := projectName
	resp := service.Projects.ServiceAccounts.List("projects/" + project)
	if err := resp.Pages(
		ctx,
		func(page *iam.ListServiceAccountsResponse) error {
			for _, account := range page.Accounts {
				d.StreamListItem(ctx, account)
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getGcpServiceAccount(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getGcpServiceAccount")
	service, err := iam.NewService(ctx)
	if err != nil {
		return nil, err
	}

	name := d.KeyColumnQuals["name"].GetStringValue()

	op, err := service.Projects.ServiceAccounts.Get(name).Do()
	if err != nil {
		logger.Debug("getGcpServiceAccount__", "ERROR", err)
		return nil, err
	}

	return op, nil
}

func getServiceAccountIamPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	account := h.Item.(*iam.ServiceAccount)
	logger := plugin.Logger(ctx)
	logger.Trace("getServiceAccountIamPolicy")

	service, err := iam.NewService(ctx)
	if err != nil {
		return nil, err
	}

	op, err := service.Projects.ServiceAccounts.GetIamPolicy(account.Name).Do()
	if err != nil {
		logger.Debug("getServiceAccountIamPolicy__", "Error", err)
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func serviceAccountNameToAkas(_ context.Context, d *transform.TransformData) (interface{}, error) {
	serviceAccount := d.HydrateItem.(*iam.ServiceAccount)

	// Get data for turbot defined properties
	akas := []string{"gcp://iam.googleapis.com/" + serviceAccount.Name}

	return akas, nil
}
