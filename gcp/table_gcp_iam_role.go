package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	"google.golang.org/api/iam/v1"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableGcpIamRole(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_iam_role",
		Description: "GCP IAM Role",
		Get: &plugin.GetConfig{
			KeyColumns:  plugin.SingleColumn("name"),
			ItemFromKey: roleNameFromKey,
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
				Transform:   transform.FromField("Role.Name"),
			},
			{
				Name:        "deleted",
				Description: "Specifies whether the role is deleted, or not",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Role.Deleted"),
			},
			{
				Name:        "role_id",
				Description: "Contains the resource type",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpIAMRoleTurbotData, "RoleID"),
			},
			{
				Name:        "is_gcp_managed",
				Description: "Specifies whether the role is GCP Managed or Customer Managed.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Type"),
			},
			{
				Name:        "stage",
				Description: "The current launch stage of the role",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Role.Stage"),
			},
			{
				Name:        "description",
				Description: "A human-readable description for the role",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Role.Description"),
			},
			{
				Name:        "etag",
				Description: "An unique read-only string that changes whenever the resource is updated",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Role.Etag"),
			},
			{
				Name:        "included_permissions",
				Description: "The names of the permissions this role grants when bound in an IAM policy",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getIamRole,
				Transform:   transform.FromField("Role.IncludedPermissions"),
			},

			// standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Role.Title"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(gcpIAMRoleTurbotData, "Akas"),
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

type roleInfo struct {
	Role *iam.Role
	Type bool
}

//// ITEM FROM KEY

func roleNameFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	name := quals["name"].GetStringValue()
	item := &roleInfo{
		Role: &iam.Role{
			Name: name,
		},
	}
	return item, nil
}

//// FETCH FUNCTIONS

func listIamRoles(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	service, err := iam.NewService(ctx)
	if err != nil {
		return nil, err
	}

	// TODO :: Need to fetch the details from env
	project := activeProject()

	// List all the project roles
	customRoles := service.Projects.Roles.List("projects/" + project)
	if err := customRoles.Pages(
		ctx,
		func(page *iam.ListRolesResponse) error {
			for _, role := range page.Roles {
				d.StreamListItem(ctx, &roleInfo{role, false})
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
				d.StreamListItem(ctx, &roleInfo{managedRole, true})
			}
			return nil
		},
	); err != nil {
		return nil, err
	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getIamRole(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getIamRole")
	roleData := h.Item.(*roleInfo)
	var op *iam.Role

	service, err := iam.NewService(ctx)
	if err != nil {
		return nil, err
	}

	gcpManaged := true

	// Checking whether the role is predefined or project role
	if strings.HasPrefix(roleData.Role.Name, "projects/") {
		gcpManaged = false
		op, err = service.Projects.Roles.Get(roleData.Role.Name).Do()
	} else {
		op, err = service.Roles.Get(roleData.Role.Name).Do()
	}

	if err != nil {
		return nil, err
	}

	return &roleInfo{op, gcpManaged}, nil
}

/// TRANSFORM FUNCTIONS

func gcpIAMRoleTurbotData(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("googleIAMRoleTurbotData")
	roleData := d.HydrateItem.(*roleInfo)
	akas := []string{"gcp://iam.googleapis.com/" + roleData.Role.Name}

	if d.Param.(string) == "RoleID" {
		return getLastPathElement(types.SafeString(roleData.Role.Name)), nil
	}

	return akas, nil
}
