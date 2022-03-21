package gcp

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
	"google.golang.org/api/cloudresourcemanager/v1"
)

//// TABLE DEFINITION

func tableGcpOrganization(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_organization",
		Description: "GCP Organization",
		List: &plugin.ListConfig{
			Hydrate:           listGCPOrganizations,
			ShouldIgnoreError: isIgnorableError([]string{"403"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "display_name",
				Description: "A human-readable string that refers to the Organization in the GCP Console UI. This string is set by the server and cannot be changed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "The resource name of the organization.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "organization_id",
				Description: "An unique, system generated ID for organization.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Name").Transform(lastPathElement),
			},
			{
				Name:        "lifecycle_state",
				Description: "The organization's current lifecycle state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "Timestamp when the Organization was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "directory_customer_id",
				Description: "The G Suite customer id used in the Directory API.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Owner.DirectoryCustomerId"),
			},

			// Steampipe standard columns
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
				Transform:   transform.From(getOrganizationAka),
			},
		},
	}
}

//// LIST FUNCTION

func listGCPOrganizations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listGCPOrganizations")

	// Create Service Connection
	service, err := CloudResourceManagerService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Max limit isn't mentioned in the documentation
	// Default limit is set as 1000
	rb := &cloudresourcemanager.SearchOrganizationsRequest{
		PageSize: *types.Int64(1000),
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < rb.PageSize {
			rb.PageSize = *limit
		}
	}

	resp := service.Organizations.Search(rb)
	if err := resp.Pages(ctx, func(page *cloudresourcemanager.SearchOrganizationsResponse) error {
		for _, organization := range page.Organizations {
			d.StreamListItem(ctx, organization)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				page.NextPageToken = ""
				return nil
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

//// TRANSFORM FUNCTIONS

func getOrganizationAka(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getOrganizationAka")

	data := d.HydrateItem.(*cloudresourcemanager.Organization)

	// Build resource aka
	akas := []string{"gcp://cloudresourcemanager.googleapis.com/organizations/" + data.DisplayName}

	return akas, nil
}
