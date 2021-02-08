package gcp

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"google.golang.org/api/bigtableadmin/v2"
)

//// TABLE DEFINITION

func tableGcpBigtableInstance(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_bigtable_instance",
		Description: "GCP Bigtable Instance",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getBigtableInstance,
		},
		List: &plugin.ListConfig{
			Hydrate: listBigtableInstances,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "A friendly name that identifies the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(bigtableInstanceToTurbotData, "Name"),
			},
			{
				Name:        "display_name",
				Description: "The descriptive name for this instance as it appears in UIs. Can be changed at any time, but should be kept globally unique to avoid conflicts.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_type",
				Description: "Specifies the type of the instance. Defaults to `PRODUCTION`.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Type"),
			},
			{
				Name:        "state",
				Description: "Specifies the current state of the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "iam_policy",
				Description: "An Identity and Access Management (IAM) policy, which specifies access controls for Google Cloud resources. A `Policy` is a collection of `bindings`. A `binding` binds one or more `members` to a single `role`. Members can be user accounts, service accounts, Google groups, and domains (such as G Suite). A `role` is a named list of permissions; each `role` can be an IAM predefined role or a user-created custom role. For some types of Google Cloud resources, a `binding` can also specify a `condition`, which is a logical expression that allows access to a resource only if the expression evaluates to `true`.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBigtableInstanceIamPolicy,
				Transform:   transform.FromValue(),
			},

			// standard steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Labels"),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(bigtableInstanceToTurbotData, "Title"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(bigtableInstanceToTurbotData, "Akas"),
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
				Transform:   transform.FromConstant(activeProject()),
			},
		},
	}
}

//// LIST FUNCTION

func listBigtableInstances(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listBigtableInstances")
	service, err := bigtableadmin.NewService(ctx)
	if err != nil {
		return nil, err
	}

	project := activeProject()
	resp := service.Projects.Instances.List("projects/" + project)
	if err := resp.Pages(ctx, func(page *bigtableadmin.ListInstancesResponse) error {
		for _, instance := range page.Instances {
			d.StreamListItem(ctx, instance)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getBigtableInstance(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	service, err := bigtableadmin.NewService(ctx)
	if err != nil {
		return nil, err
	}

	name := d.KeyColumnQuals["name"].GetStringValue()
	project := activeProject()

	resp, err := service.Projects.Instances.Get("projects/" + project + "/instances/" + name).Do()
	if err != nil {
		return nil, err
	}

	// If the name filed kept as empty, API does not return any errors
	if len(resp.Name) < 1 {
		return nil, nil
	}

	return resp, nil
}

func getBigtableInstanceIamPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	service, err := bigtableadmin.NewService(ctx)
	if err != nil {
		return nil, err
	}

	instance := h.Item.(*bigtableadmin.Instance)
	getIamPolicyRequest := bigtableadmin.GetIamPolicyRequest{}

	req, err := service.Projects.Instances.GetIamPolicy(instance.Name, &getIamPolicyRequest).Do()
	if err != nil {
		// Return nil, if the resource not present
		result := isNotFoundError([]string{"404"})
		if result != nil {
			return nil, nil
		}
		return nil, err
	}

	return req, nil
}

//// TRANSFORM FUNCTIONS

func bigtableInstanceToTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	instance := d.HydrateItem.(*bigtableadmin.Instance)
	param := d.Param.(string)

	// get the resource title
	splittedTitle := strings.Split(instance.Name, "/")

	turbotData := map[string]interface{}{
		"Name":  splittedTitle[len(splittedTitle)-1],
		"Title": splittedTitle[len(splittedTitle)-1],
		"Akas":  []string{"gcp://bigtableadmin.googleapis.com/" + instance.Name},
	}

	return turbotData[param], nil
}
