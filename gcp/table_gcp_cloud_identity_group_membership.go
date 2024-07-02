package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/cloudidentity/v1"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableGcpCloudIdentityGroupMembership(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_cloud_identity_group_membership",
		Description: "GCP Cloud Identity Group Membership",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "group_name"}),
			Hydrate:    getCloudIdentityGroupMembership,
			// Since there are no specific permissions required for get the group membership, the rate-limiter tag is applied based on the method ID.
			// https://cloud.google.com/identity/docs/reference/rest/v1/groups.memberships/get
			Tags: map[string]string{"service": "cloudidentity", "action": "groups.memberships.get"},
		},
		List: &plugin.ListConfig{
			Hydrate:           listCloudIdentityGroupMemberships,
			ShouldIgnoreError: isIgnorableError([]string{"400"}),
			KeyColumns:        plugin.AllColumns([]string{"group_name"}),
			// Since there are no specific permissions required for listing the group memberships, the rate-limiter tag is applied based on the method ID.
			// https://cloud.google.com/identity/docs/reference/rest/v1/groups.memberships/list
			Tags: map[string]string{"service": "cloudidentity", "action": "groups.memberships.list"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies the membership.",
				Transform:   transform.FromField("Name").Transform(lastPathElement),
			},
			{
				Name:        "group_name",
				Description: "The group in which the membership belongs to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(gcpCloudIdentityGroupName),
			},
			{
				Name:        "create_time",
				Description: "The time when the membership was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromGo().NullIfZero(),
			},
			{
				Name:        "type",
				Description: "The type of the membership.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "update_time",
				Description: "The time when the membership was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromGo().NullIfZero(),
			},

			// JSON Columns
			{
				Name:        "preferred_member_key",
				Description: "The `EntityKey` of the member.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "roles",
				Description: "The membership roles that apply to the membership.",
				Type:        proto.ColumnType_JSON,
			},

			// standard Steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(lastPathElement),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(gcpCloudIdentityGroupMembershipTurbotData, "Akas"),
			},

			// standard GCP columns
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

func listCloudIdentityGroupMemberships(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	groupName := d.EqualsQualString("group_name")

	// Return nil, if no input provided
	if groupName == "" {
		return nil, nil
	}

	// Max limit is set as per documentation
	// https://pkg.go.dev/google.golang.org/api/cloudidentity/v1#GroupsMembershipsListCall.PageSize
	pageSize := types.Int64(1000)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *pageSize {
			pageSize = limit
		}
	}

	// Create Service Connection
	service, err := CloudIdentityService(ctx, d)
	if err != nil {
		logger.Error("gcp_cloud_identity_group_membership.listCloudIdentityGroupMemberships", "connection_error", err)
		return nil, err
	}

	memberships := service.Groups.Memberships.List("groups/" + groupName).View("FULL").PageSize(*pageSize)
	if err := memberships.Pages(ctx, func(page *cloudidentity.ListMembershipsResponse) error {
		for _, membership := range page.Memberships {
			d.StreamListItem(ctx, membership)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				break
			}
		}
		return nil
	}); err != nil {
		logger.Error("gcp_cloud_identity_group_membership.listCloudIdentityGroupMemberships", "api_error", err)
		return nil, err
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getCloudIdentityGroupMembership(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	name := d.EqualsQualString("name")
	groupName := d.EqualsQualString("group_name")

	// Return nil, if no input provided
	if name == "" || groupName == "" {
		return nil, nil
	}

	// Create Service Connection
	service, err := CloudIdentityService(ctx, d)
	if err != nil {
		logger.Error("gcp_cloud_identity_group_membership.getCloudIdentityGroupMembership", "connection_error", err)
		return nil, err
	}

	membership, err := service.Groups.Memberships.Get("groups/" + groupName + "/memberships/" + name).Do()
	if err != nil {
		logger.Error("gcp_cloud_identity_group_membership.getCloudIdentityGroupMembership", "api_error", err)
		return nil, err
	}

	return membership, nil
}

/// TRANSFORM FUNCTIONS

func gcpCloudIdentityGroupMembershipTurbotData(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	membership := d.HydrateItem.(*cloudidentity.Membership)
	akas := []string{"gcp://cloudidentity.googleapis.com/" + membership.Name}

	return akas, nil
}

func gcpCloudIdentityGroupName(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	membership := d.HydrateItem.(*cloudidentity.Membership)
	groupName := strings.Split(membership.Name, "/")[1]

	return groupName, nil
}
