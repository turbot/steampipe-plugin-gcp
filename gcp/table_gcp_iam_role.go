package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/iam/v1"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableGcpIamRole(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_iam_role",
		Description: "GCP IAM Role",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getIamRole,
			Tags:       map[string]string{"service": "iam", "action": "roles.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listIamRoles,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "is_gcp_managed", Require: plugin.Optional, Operators: []string{"<>", "="}},
			},
			Tags: map[string]string{"service": "iam", "action": "roles.list"},
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
				Hydrate:     getProject,
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

func listIamRoles(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

	showDeleted := helpers.StringSliceContains(d.QueryContext.Columns, "deleted")
	view := "BASIC"
	if helpers.StringSliceContains(d.QueryContext.Columns, "included_permissions") {
		view = "FULL"
	}

	roleType := "ALL"
	// Handling for non equal quals
	if d.Quals["is_gcp_managed"] != nil {
		for _, q := range d.Quals["is_gcp_managed"].Quals {
			value := q.Value.GetBoolValue()
			roleType = "CUSTOM"
			switch q.Operator {
			case "<>":
				if !value {
					roleType = "GCP"
				}
			case "=":
				if value {
					roleType = "GCP"
				}
			}
		}
	}

	plugin.Logger(ctx).Error("listIamRoles", "roleType", roleType)

	// Max limit is set as per documentation
	// https://pkg.go.dev/google.golang.org/api@v0.48.0/iam/v1#ProjectsRolesListCall.PageSize
	pageSize := types.Int64(1000)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *pageSize {
			pageSize = limit
		}
	}

	if roleType == "ALL" || roleType == "CUSTOM" {
		// List all the custom project roles
		customRoles := service.Projects.Roles.List("projects/" + project).View(view).ShowDeleted(showDeleted).PageSize(*pageSize)
		if err := customRoles.Pages(ctx, func(page *iam.ListRolesResponse) error {
			// apply rate limiting
		d.WaitForListRateLimit(ctx)

			for _, role := range page.Roles {
				d.StreamListItem(ctx, &roleInfo{role, false})

				// Check if context has been cancelled or if the limit has been hit (if specified)
				// if there is a limit, it will return the number of rows required to reach this limit
				if d.RowsRemaining(ctx) == 0 {
					page.NextPageToken = ""
					break
				}
			}
			return nil
		},
		); err != nil {
			return nil, err
		}
	}

	if roleType == "ALL" || roleType == "GCP" {
		// List all the pre-defined roles
		managedRole := service.Roles.List().View(view).ShowDeleted(showDeleted).PageSize(*pageSize)
		if err := managedRole.Pages(ctx, func(page *iam.ListRolesResponse) error {
			// apply rate limiting
			d.WaitForListRateLimit(ctx)

			for _, managedRole := range page.Roles {
				d.StreamListItem(ctx, &roleInfo{managedRole, true})

				// Check if context has been cancelled or if the limit has been hit (if specified)
				// if there is a limit, it will return the number of rows required to reach this limit
				if d.RowsRemaining(ctx) == 0 {
					page.NextPageToken = ""
					break
				}
			}

			return nil
		},
		); err != nil {
			return nil, err
		}
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
		quals := d.EqualsQuals
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
