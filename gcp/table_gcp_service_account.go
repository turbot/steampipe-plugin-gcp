package gcp

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"google.golang.org/api/iam/v1"
)

//// TABLE DEFINITION

func tableGcpServiceAccount(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_service_account",
		Description: "GCP Service Account",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getGcpServiceAccount,
			Tags:       map[string]string{"service": "iam", "action": "serviceAccounts.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listGcpServiceAccounts,
			Tags:    map[string]string{"service": "iam", "action": "serviceAccounts.list"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getServiceAccountIamPolicy,
				Tags: map[string]string{"service": "iam", "action": "serviceAccounts.getIamPolicy"},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The resource name of the service account",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(lastPathElement),
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
				Name:        "disabled",
				Description: "Specifies whether the service is account is disabled, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "description",
				Description: "A user-specified, human-readable description of the service account. The maximum length is 256 UTF-8 bytes.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "oauth2_client_id",
				Description: "The OAuth 2.0 client ID for the service account.",
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
				Transform:   transform.FromField("ProjectId"),
			},
		},
	}
}

//// LIST FUNCTION

func listGcpServiceAccounts(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := IAMService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Max limit is set as per documentation
	// https://pkg.go.dev/google.golang.org/api@v0.48.0/iam/v1?utm_source=gopls#ProjectsServiceAccountsListCall.PageSize
	pageSize := types.Int64(100)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *pageSize {
			pageSize = limit
		}
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	resp := service.Projects.ServiceAccounts.List("projects/" + project).PageSize(*pageSize)
	if err := resp.Pages(
		ctx,
		func(page *iam.ListServiceAccountsResponse) error {
			// apply rate limiting
			d.WaitForListRateLimit(ctx)

			for _, account := range page.Accounts {
				d.StreamListItem(ctx, account)

				// Check if context has been cancelled or if the limit has been hit (if specified)
				// if there is a limit, it will return the number of rows required to reach this limit
				if d.RowsRemaining(ctx) == 0 {
					page.NextPageToken = ""
					return nil
				}
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
	plugin.Logger(ctx).Trace("getGcpServiceAccount")
	// Create Service Connection
	service, err := IAMService(ctx, d)
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

	// If the name has been passed as empty string, API does not returns any error
	if len(name) < 1 {
		return nil, nil
	}
	accountName := "projects/" + project + "/serviceAccounts/" + name

	op, err := service.Projects.ServiceAccounts.Get(accountName).Do()
	if err != nil {
		plugin.Logger(ctx).Debug("getGcpServiceAccount__", "ERROR", err)
		return nil, err
	}

	return op, nil
}

func getServiceAccountIamPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	account := h.Item.(*iam.ServiceAccount)
	plugin.Logger(ctx).Trace("getServiceAccountIamPolicy")

	// Create Service Connection
	service, err := IAMService(ctx, d)
	if err != nil {
		return nil, err
	}

	op, err := service.Projects.ServiceAccounts.GetIamPolicy(account.Name).Do()
	if err != nil {
		plugin.Logger(ctx).Debug("getServiceAccountIamPolicy__", "Error", err)
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
