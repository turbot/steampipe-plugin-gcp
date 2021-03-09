package gcp

import (
	"context"
	"os"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	"google.golang.org/api/iam/v1"
)

// TABLE DEFINITION
func tableGcpIamRole(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_iam_role",
		Description: "GCP IAM Role",
		Get: &plugin.GetConfig{
			KeyColumns:  plugin.SingleColumn("name"),
			Hydrate:     getIamRole,
		},
		List: &plugin.ListConfig{
			Hydrate: listIamRoles,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the role",
			},
			{
				Name:        "deleted",
				Description: "Specifies whether the role is deleted, or not",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "role_id",
				Description: "Contains the resource type",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpIAMRoleTurbotData, "RoleID"),
			},
			{
				Name:        "stage",
				Description: "The current launch stage of the role",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A human-readable description for the role",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "etag",
				Description: "An unique read-only string that changes whenever the resource is updated",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "included_permissions",
				Description: "The names of the permissions this role grants when bound in an IAM policy",
				Type:        proto.ColumnType_JSON,
			},

			// Standard columns for all tables
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "akas",
				Description: "Array of globally unique identifier strings (also known as) for the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(gcpIAMRoleTurbotData, "TurbotAkas"),
			},
		},
	}
}

// ITEM FROM KEY
func roleNameFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	name := quals["name"].GetStringValue()
	item := &iam.Role{
		Name: name,
	}
	return item, nil
}

// FETCH FUNCTIONS
func listIamRoles(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	service, err := iam.NewService(ctx)
	if err != nil {
		return nil, err
	}

	// TODO :: Need to fetch the details from env
	project := os.Getenv("GCP_PROJECT")

	// List all the project roles
	customRoles := service.Projects.Roles.List("projects/" + project)
	if err := customRoles.Pages(
		ctx,
		func(page *iam.ListRolesResponse) error {
			for _, role := range page.Roles {
				d.StreamListItem(ctx, role)
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	// List all the pre-defined roles
	managedRole := service.Roles.List()
	if err := managedRole.Pages(
		ctx,
		func(page *iam.ListRolesResponse) error {
			for _, managedRole := range page.Roles {
				d.StreamListItem(ctx, managedRole)
			}
			return nil
		},
	); err != nil {
		return nil, err
	}
	return nil, err
}

// HYDRATE FUNCTIONS
func getIamRole(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	role := h.Item.(*iam.Role)
	var op *iam.Role

	service, err := iam.NewService(ctx)
	if err != nil {
		return nil, err
	}

	// Checking whether the role is predefined or project role
	if strings.HasPrefix(role.Name, "projects/") {
		op, err = service.Projects.Roles.Get(role.Name).Do()
	} else {
		op, err = service.Roles.Get(role.Name).Do()
	}

	if err != nil {
		return nil, err
	}

	return op, nil
}

// TRANSFORM FUNCTIONS
func gcpIAMRoleTurbotData(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("googleIAMRoleTurbotData")
	role := d.HydrateItem.(*iam.Role)

	splitName := strings.Split(role.Name, "/")
	akas := []string{"gcp://iam.googleapis.com/" + role.Name}

	if d.Param.(string) == "RoleID" {
		return splitName[len(splitName)-1], nil
	}
	return akas, nil
}