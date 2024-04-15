package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"google.golang.org/api/pubsub/v1"
)

func tableGcpPubSubSnapshot(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_pubsub_snapshot",
		Description: "GCP Pub/Sub Snapshot",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getPubSubSnapshot,
			Tags:       map[string]string{"service": "pubsub", "action": "snapshots.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listPubSubSnapshots,
			Tags:    map[string]string{"service": "pubsub", "action": "snapshots.list"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getPubSubSnapshotIamPolicy,
				Tags: map[string]string{"service": "pubsub", "action": "snapshots.getIamPolicy"},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the snapshot",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(lastPathElement),
			},
			{
				Name:        "topic",
				Description: "The name of the topic from which this snapshot is retaining messages",
				Type:        proto.ColumnType_STRING,
			},
			// topic_name is a simpler view of the topic, without the full path
			{
				Name:        "topic_name",
				Description: "The short name of the topic from which this snapshot is retaining messages.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Topic").Transform(lastPathElement),
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(pubsubSnapshotSelfLink),
			},
			{
				Name:        "expire_time",
				Description: "The snapshot is guaranteed to exist up until this time. A newly-created snapshot expires no later than 7 days from the time of its creation. Its exact lifetime is determined at creation by the existing backlog in the source subscription. Specifically, the lifetime of the snapshot is `7 days - (age of oldest unacked message in the subscription)`.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "iam_policy",
				Description: "An Identity and Access Management (IAM) policy, which specifies access controls for Google Cloud resources. A `Policy` is a collection of `bindings`. A `binding` binds one or more `members` to a single `role`. Members can be user accounts, service accounts, Google groups, and domains (such as G Suite). A `role` is a named list of permissions; each `role` can be an IAM predefined role or a user-created custom role. For some types of Google Cloud resources, a `binding` can also specify a `condition`, which is a logical expression that allows access to a resource only if the expression evaluates to `true`.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getPubSubSnapshotIamPolicy,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "labels",
				Description: "A set of labels attached with the snapshot.",
				Type:        proto.ColumnType_JSON,
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
				Transform:   transform.FromField("Name").Transform(lastPathElement),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(snapshotNameToTurbotData, "Akas"),
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
				Transform:   transform.FromP(snapshotNameToTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listPubSubSnapshots(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := PubsubService(ctx, d)
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

	resp := service.Projects.Snapshots.List("projects/" + project).PageSize(*pageSize)
	if err := resp.Pages(ctx, func(page *pubsub.ListSnapshotsResponse) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, snapshot := range page.Snapshots {
			d.StreamListItem(ctx, snapshot)

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

	return nil, err
}

//// HYDRATE FUNCTIONS

func getPubSubSnapshot(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getPubSubSnapshot")

	// Create Service Connection
	service, err := PubsubService(ctx, d)
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

	req, err := service.Projects.Snapshots.Get("projects/" + project + "/snapshots/" + name).Do()
	if err != nil {
		return nil, err
	}
	return req, nil
}

func getPubSubSnapshotIamPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getPubSubSnapshotIamPolicy")

	// Create Service Connection
	service, err := PubsubService(ctx, d)
	if err != nil {
		return nil, err
	}

	resource := h.Item.(*pubsub.Snapshot)
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
		"Akas":    []string{"gcp://pubsub.googleapis.com/" + snapshot.Name},
	}

	return turbotData[param], nil
}

func pubsubSnapshotSelfLink(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*pubsub.Snapshot)
	selfLink := "https://pubsub.googleapis.com/v1/" + data.Name

	return selfLink, nil
}
