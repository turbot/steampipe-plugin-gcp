package gcp

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/cloudidentity/v1"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableGcpCloudIdentityGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_cloud_identity_group",
		Description: "GCP Cloud Identity Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getCloudIdentityGroup,
			// Since there are no specific permissions required for get the group, the rate-limiter tag is applied based on the method ID.
			// https://cloud.google.com/identity/docs/reference/rest/v1/groups/get
			Tags: map[string]string{"service": "cloudidentity", "action": "groups.get"},
		},
		List: &plugin.ListConfig{
			Hydrate:           listCloudIdentityGroups,
			ShouldIgnoreError: isIgnorableError([]string{"400"}),
			KeyColumns:        plugin.SingleColumn("parent"),
			// Since there are no specific permissions required for listing groups, the rate-limiter tag is applied based on the method ID.
			// https://cloud.google.com/identity/docs/reference/rest/v1/groups/list
			Tags: map[string]string{"service": "cloudidentity", "action": "groups.list"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the group.",
				Transform:   transform.FromField("Name").Transform(lastPathElement),
			},
			{
				Name:        "create_time",
				Description: "The time when the group was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromGo().NullIfZero(),
			},
			{
				Name:        "description",
				Description: "A human-readable description for the group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "A human-readable name for the group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "parent",
				Description: "The resource name of the entity under which this `Group` resides in the Cloud Identity resource hierarchy.",
				Transform:   transform.FromQual("parent"),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "update_time",
				Description: "The time when the group was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromGo().NullIfZero(),
			},

			// JSON Columns
			{
				Name:        "dynamic_group_metadata",
				Description: "Dynamic group metadata like queries and status.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "group_key",
				Description: "The `EntityKey` of the `Group`.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "labels",
				Description: "The labels that apply to the group.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Labels"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(gcpCloudIdentityGroupTurbotData, "Akas"),
			},

			// GCP standard columns
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

//// LIST FUNCTIONS

func listCloudIdentityGroups(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	parent := d.EqualsQualString("parent")

	// The parent should not be empty
	if parent == "" {
		return nil, nil
	}

	// Max limit is set as per documentation
	// https://pkg.go.dev/google.golang.org/api/cloudidentity/v1#GroupsListCall.PageSize
	pageSize := types.Int64(500)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *pageSize {
			pageSize = limit
		}
	}

	// Create Service Connection
	service, err := CloudIdentityService(ctx, d)
	if err != nil {
		logger.Error("gcp_cloud_identity_group.listCloudIdentityGroups", "connection_error", err)
		return nil, err
	}

	groups := service.Groups.List().View("FULL").PageSize(*pageSize).Parent(parent)
	if err := groups.Pages(ctx, func(page *cloudidentity.ListGroupsResponse) error {
		for _, group := range page.Groups {
			d.StreamListItem(ctx, group)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				page.NextPageToken = ""
				break
			}
		}
		return nil
	}); err != nil {
		logger.Error("gcp_cloud_identity_group.listCloudIdentityGroups", "api_error", err)
		return nil, err
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getCloudIdentityGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	name := d.EqualsQualString("name")

	// Return nil, if no input provided
	if name == "" {
		return nil, nil
	}

	// Create Service Connection
	service, err := CloudIdentityService(ctx, d)
	if err != nil {
		logger.Error("gcp_cloud_identity_group.getCloudIdentityGroup", "connection_error", err)
		return nil, err
	}

	group, err := service.Groups.Get("groups/" + name).Do()
	if err != nil {
		logger.Error("gcp_cloud_identity_group.getCloudIdentityGroup", "api_error", err)
		return nil, err
	}

	return group, nil
}

/// TRANSFORM FUNCTIONS

func gcpCloudIdentityGroupTurbotData(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	group := d.HydrateItem.(*cloudidentity.Group)
	akas := []string{"gcp://cloudidentity.googleapis.com/" + group.Name}

	return akas, nil
}
