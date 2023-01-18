package gcp

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
	"google.golang.org/api/apikeys/v2"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableGcpApiKey(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_api_key",
		Description: "GCP API Key",
		// Get: &plugin.GetConfig{
		// 	KeyColumns: plugin.SingleColumn("name"),
		// 	Hydrate:    getIamRole,
		// },
		List: &plugin.ListConfig{
			Hydrate:           listApiKeys,
			ShouldIgnoreError: isIgnorableError([]string{"404"}),
			// KeyColumns: plugin.KeyColumnSlice{
			// 	{Name: "is_gcp_managed", Require: plugin.Optional, Operators: []string{"<>", "="}},
			// },
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The resource name of the key.",
			},
			{
				Name:        "uid",
				Type:        proto.ColumnType_STRING,
				Description: "Unique id in UUID4 format.",
			},
			{
				Name:        "create_time",
				Description: "A timestamp identifying the time this key was originally created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "delete_time",
				Description: "A timestamp when this key was deleted.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Type:        proto.ColumnType_STRING,
				Description: "Human-readable display name of this key that you can modify.",
			},
			{
				Name:        "etag",
				Description: "A checksum computed by the server based on the current value of the Key resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "key_string",
				Description: "An encrypted and signed value held by this key.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "update_time",
				Description: "A timestamp identifying the time this key was last updated.",
				Type:        proto.ColumnType_STRING,
			},

			// JSON columns
			{
				Name:        "annotations",
				Description: "Annotations is an unstructured key-value map stored with a policy that may be set by external tools to store and retrieve arbitrary metadata.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "restrictions",
				Description: "The restrictions on the key.",
				Type:        proto.ColumnType_STRING,
			},

			// standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
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

//// FETCH FUNCTIONS

func listApiKeys(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listApiKeys")

	// Create Service Connection
	service, err := APIService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Max limit isn't mentioned in the documentation
	// Default limit is set as 1000
	pageSize := types.Int64(1000)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *pageSize {
			pageSize = limit
		}
	}

	// Get project details
	getProjectCached := plugin.HydrateFunc(getProject).WithCache()
	projectId, err := getProjectCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	resp := service.Projects.Locations.Keys.List(project).PageSize(*pageSize)

	if err := resp.Pages(
		ctx,
		func(page *apikeys.V2ListKeysResponse) error {
			for _, item := range page.Keys {
				d.StreamListItem(ctx, item)

				// Check if context has been cancelled or if the limit has been hit (if specified)
				// if there is a limit, it will return the number of rows required to reach this limit
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					page.NextPageToken = ""
					return nil
				}
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

// func getIamRole(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 	plugin.Logger(ctx).Trace("getIamRole")

// 	var name string
// 	if h.Item != nil {
// 		name = h.Item.(*roleInfo).Role.Name
// 	} else {
// 		quals := d.KeyColumnQuals
// 		name = quals["name"].GetStringValue()
// 	}

// 	var op *iam.Role

// 	// Create Service Connection
// 	service, err := IAMService(ctx, d)
// 	if err != nil {
// 		return nil, err
// 	}

// 	gcpManaged := true

// 	// Checking whether the role is predefined or project role
// 	if strings.HasPrefix(name, "projects/") {
// 		gcpManaged = false
// 		op, err = service.Projects.Roles.Get(name).Do()
// 	} else {
// 		op, err = service.Roles.Get(name).Do()
// 	}

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &roleInfo{op, gcpManaged}, nil
// }

// /// TRANSFORM FUNCTIONS

// func gcpIAMRoleTurbotData(ctx context.Context, d *transform.TransformData) (interface{}, error) {
// 	plugin.Logger(ctx).Trace("googleIAMRoleTurbotData")
// 	roleData := d.HydrateItem.(*roleInfo)
// 	akas := []string{"gcp://iam.googleapis.com/" + roleData.Role.Name}

// 	if d.Param.(string) == "RoleID" {
// 		return getLastPathElement(types.SafeString(roleData.Role.Name)), nil
// 	}

// 	return akas, nil
// }
