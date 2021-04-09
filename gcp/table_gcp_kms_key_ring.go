package gcp

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	"google.golang.org/api/cloudkms/v1"
)

func tableGcpKmsKeyRing(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_kms_key_ring",
		Description: "GCP Kms Key Ring",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getKeyRingDetail,
		},
		List: &plugin.ListConfig{
			Hydrate: listKeyRingDetails,
		},
		GetMatrixItem: BuildLocationList,
		Columns: []*plugin.Column{
			// commonly used columns
			{
				Name:        "name",
				Description: "The resource name for the KeyRing.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "The time at which this KeyRing was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},

			// standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpKmsKeyRingTurbotData, "Title"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(gcpKmsKeyRingTurbotData, "Akas"),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpKmsKeyRingTurbotData, "Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(gcpKmsKeyRingTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTIONS
func listKeyRingDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	location := plugin.GetMatrixItem(ctx)[matrixKeyLocation].(string)
	plugin.Logger(ctx).Trace("listKeyRingDetails")

	// Create Service Connection
	service, err := KMSService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}
	project := projectData.Project
	resp := service.Projects.Locations.KeyRings.List("projects/" + project + "/locations/" + location)
	if err := resp.Pages(ctx, func(page *cloudkms.ListKeyRingsResponse) error {
		for _, ring := range page.KeyRings {
			d.StreamListItem(ctx, ring)
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

	// Create Service Connection
	service, err := KMSService(ctx, d)
	if err != nil {
		return nil, err
	}

	name := d.KeyColumnQuals["name"].GetStringValue()
	resp, err := service.Projects.Locations.KeyRings.Get(name).Do()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

//// TRANSFORM FUNCTIONS

func gcpKmsKeyRingTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	key := d.HydrateItem.(*cloudkms.KeyRing)
	param := d.Param.(string)

	project := strings.Split(key.Name, "/")[1]
	location := strings.Split(key.Name, "/")[3]
	title := strings.Split(key.Name, "/")[5]

	turbotData := map[string]interface{}{
		"Project":  project,
		"Location": location,
		"Title":    title,
		"Akas":     []string{"gcp://cloudkms.googleapis.com/" + key.Name},
	}

	return turbotData[param], nil
}
