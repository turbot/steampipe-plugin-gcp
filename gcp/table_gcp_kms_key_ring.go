package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/cloudkms/v1"
)

//// TABLE DEFINITION

func tableGcpKmsKeyRing(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_kms_key_ring",
		Description: "GCP KMS Key Ring",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "location"}),
			Hydrate:    getKeyRingDetail,
			Tags:       map[string]string{"service": "cloudkms", "action": "keyRings.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listKeyRingDetails,
			Tags:    map[string]string{"service": "cloudkms", "action": "keyRings.list"},
		},
		GetMatrixItemFunc: BuildLocationList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The resource name for the KeyRing.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(lastPathElement),
			},
			{
				Name:        "create_time",
				Description: "The time at which this KeyRing was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "iam_policy",
				Description: "An Identity and Access Management (IAM) policy, which specifies access controls for Google Cloud resources. A `Policy` is a collection of `bindings`. A `binding` binds one or more `members` to a single `role`. Members can be user accounts, service accounts, Google groups, and domains (such as G Suite). A `role` is a named list of permissions; each `role` can be an IAM predefined role or a user-created custom role. For some types of Google Cloud resources, a `binding` can also specify a `condition`, which is a logical expression that allows access to a resource only if the expression evaluates to `true`.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getKmsKeyRingIamPolicy,
				Transform:   transform.FromValue(),
			},

			// Standard steampipe columns
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
				Hydrate:     gcpKmsKeyRingTurbotData,
				Transform:   transform.FromField("Akas"),
			},

			// GCP standard columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Hydrate:     gcpKmsKeyRingTurbotData,
				Transform:   transform.FromField("Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Hydrate:     gcpKmsKeyRingTurbotData,
				Transform:   transform.FromField("Project"),
			},
		},
	}
}

// // LIST FUNCTION
func listKeyRingDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listKeyRingDetails")

	var location string
	matrixLocation := d.EqualsQualString(matrixKeyLocation)
	// Since, when the service API is disabled, matrixLocation value will be nil
	if matrixLocation != "" {
		location = matrixLocation
	}

	// Create Service Connection
	service, err := KMSService(ctx, d)
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

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	resp := service.Projects.Locations.KeyRings.List("projects/" + project + "/locations/" + location).PageSize(*pageSize)
	if err := resp.Pages(ctx, func(page *cloudkms.ListKeyRingsResponse) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, ring := range page.KeyRings {
			d.StreamListItem(ctx, ring)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
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

//// HYDRATE FUNCTIONS

func getKeyRingDetail(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getKeyRingDetail")

	var location string
	matrixLocation := d.EqualsQualString(matrixKeyLocation)
	// Since, when the service API is disabled, matrixLocation value will be nil
	if matrixLocation != "" {
		location = matrixLocation
	}
	// Create Service Connection
	service, err := KMSService(ctx, d)
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
	loc := d.EqualsQuals["location"].GetStringValue()

	// to prevent duplicate value
	if location != loc {
		return nil, nil
	}
	resp, err := service.Projects.Locations.KeyRings.Get("projects/" + project + "/locations/" + loc + "/keyRings/" + name).Do()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func getKmsKeyRingIamPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getKmsKeyRingIamPolicy")

	// Create Service Connection
	service, err := KMSService(ctx, d)
	if err != nil {
		return nil, err
	}

	keyRing := h.Item.(*cloudkms.KeyRing)

	resp, err := service.Projects.Locations.KeyRings.GetIamPolicy(keyRing.Name).Do()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

//// TRANSFORM FUNCTIONS

func gcpKmsKeyRingTurbotData(_ context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	key := h.Item.(*cloudkms.KeyRing)

	data := strings.Split(key.Name, "/")

	turbotData := map[string]interface{}{
		"Project":  data[1],
		"Location": data[3],
		"Akas":     []string{"gcp://cloudkms.googleapis.com/" + key.Name},
	}

	return turbotData, nil
}
