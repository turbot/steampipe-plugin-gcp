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
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getIamRole,
		},
		List: &plugin.ListConfig{
			Hydrate:           listIamRoles,
			ShouldIgnoreError: isIgnorableError([]string{"403"}),
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

type roleInfo struct {
	Role *iam.Role
	Type bool
}

//// FETCH FUNCTIONS

func listIamRoles(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := IAMService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}
	project := projectData.Project

	// List all the project roles
	customRoles := service.Projects.Roles.List("projects/" + project).View("FULL")
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
	managedRole := service.Roles.List().View("FULL")
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

	var name string
	if h.Item != nil {
		name = h.Item.(*roleInfo).Role.Name
	} else {
		quals := d.KeyColumnQuals
		name = quals["name"].GetStringValue()
	}

	var op *iam.Role

	// Create Service Connection
	service, err := IAMService(ctx, d)
	if err != nil {
		return nil, err
	}

	gcpManaged := true

	// Checking whether the role is predefined or project role
	if strings.HasPrefix(name, "projects/") {
		gcpManaged = false
		op, err = service.Projects.Roles.Get(name).Do()
	} else {
		op, err = service.Roles.Get(name).Do()
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
