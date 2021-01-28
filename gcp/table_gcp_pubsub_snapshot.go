package gcp

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	"google.golang.org/api/pubsub/v1"
)

func tableGcpPubSubSnapshot(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_pubsub_snapshot",
		Description: "GCP Pub/Sub Snapshot",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getPubSubSnapshot,
		},
		List: &plugin.ListConfig{
			Hydrate: listPubSubSnapshot,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the snapshot",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(snapshotNameToTurbotData, "Name"),
			},
			{
				Name:        "topic",
				Description: "The name of the topic from which this snapshot is retaining messages",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "expire_time",
				Description: "The snapshot is guaranteed to exist up until this time. A newly-created snapshot expires no later than 7 days from the time of its creation. Its exact lifetime is determined at creation by the existing backlog in the source subscription. Specifically, the lifetime of the snapshot is `7 days - (age of oldest unacked message in the subscription)`.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "iam_policy",
				Description: "An Identity and Access Management (IAM) policy, which specifies access controls for Google Cloud resources. A `Policy` is a collection of `bindings`. A `binding` binds one or more `members` to a single `role`. Members can be user accounts, service accounts, Google groups, and domains (such as G Suite). A `role` is a named list of permissions; each `role` can be an IAM predefined role or a user-created custom role. For some types of Google Cloud resources, a `binding` can also specify a `condition`, which is a logical expression that allows access to a resource only if the expression evaluates to `true`.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getPubSubSnapshotIamPolicy,
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
				Transform:   transform.FromP(snapshotNameToTurbotData, "Title"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(snapshotNameToTurbotData, "Akas"),
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

//// LIST FUNCTION

func listPubSubSnapshot(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	project := activeProject()

	service, err := pubsub.NewService(ctx)
	if err != nil {
		return nil, err
	}

	resp := service.Projects.Snapshots.List("projects/" + project)
	if err := resp.Pages(ctx, func(page *pubsub.ListSnapshotsResponse) error {
		for _, snapshot := range page.Snapshots {
			d.StreamListItem(ctx, snapshot)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getPubSubSnapshot(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	name := d.KeyColumnQuals["name"].GetStringValue()
	project := activeProject()
	service, err := pubsub.NewService(ctx)
	if err != nil {
		return nil, err
	}

	req, err := service.Projects.Snapshots.Get("projects/" + project + "/snapshots/" + name).Do()
	if err != nil {
		return nil, err
	}
	return req, nil
}

func getPubSubSnapshotIamPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	resource := h.Item.(*pubsub.Snapshot)

	service, err := pubsub.NewService(ctx)
	if err != nil {
		return nil, err
	}

	req, err := service.Projects.Snapshots.GetIamPolicy(resource.Name).Do()
	if err != nil {
		return nil, err
	}

	return req, nil
}

//// TRANSFORM FUNCTIONS

func snapshotNameToTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	snapshot := d.HydrateItem.(*pubsub.Snapshot)
	param := d.Param.(string)

	// get the resource title
	splittedTitle := strings.Split(snapshot.Name, "/")

	turbotData := map[string]interface{}{
		"Project": splittedTitle[1],
		"Name":    splittedTitle[len(splittedTitle)-1],
		"Title":   splittedTitle[len(splittedTitle)-1],
		"Akas":    []string{"gcp://pubsub.googleapis.com/" + snapshot.Name},
	}

	return turbotData[param], nil
}
