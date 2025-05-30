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

func tableGcpPubSubTopic(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_pubsub_topic",
		Description: "GCP Pub/Sub Topic",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getPubSubTopic,
			Tags:       map[string]string{"service": "pubsub", "action": "topics.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listPubSubTopics,
			Tags:    map[string]string{"service": "pubsub", "action": "topics.list"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getPubSubTopicIamPolicy,
				Tags: map[string]string{"service": "pubsub", "action": "topics.getIamPolicy"},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the topic.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(lastPathElement),
			},
			{
				Name:        "kms_key_name",
				Description: "The resource name of the Cloud KMS CryptoKey to be used to protect access to messages published on this topic.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(pubsubTopicSelfLink),
			},
			{
				Name:        "message_storage_policy_allowed_persistence_regions",
				Description: "Policy constraining the set of Google Cloud Platform regions where messages published to the topic may be stored. If not present, then no constraints are in effect.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("MessageStoragePolicy.AllowedPersistenceRegions"),
			},
			{
				Name:        "iam_policy",
				Description: "An Identity and Access Management (IAM) policy, which specifies access controls for Google Cloud resources. A `Policy` is a collection of `bindings`. A `binding` binds one or more `members` to a single `role`. Members can be user accounts, service accounts, Google groups, and domains (such as G Suite). A `role` is a named list of permissions; each `role` can be an IAM predefined role or a user-created custom role. For some types of Google Cloud resources, a `binding` can also specify a `condition`, which is a logical expression that allows access to a resource only if the expression evaluates to `true`.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getPubSubTopicIamPolicy,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "labels",
				Description: "A set of labels attached with the topic.",
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
				Transform:   transform.FromP(topicNameToTurbotData, "Akas"),
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
				Transform:   transform.FromP(topicNameToTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listPubSubTopics(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

	resp := service.Projects.Topics.List("projects/" + project).PageSize(*pageSize)
	if err := resp.Pages(ctx, func(page *pubsub.ListTopicsResponse) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, topic := range page.Topics {
			d.StreamListItem(ctx, topic)

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

func getPubSubTopic(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getPubSubTopic")

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
	req, err := service.Projects.Topics.Get("projects/" + project + "/topics/" + name).Do()
	if err != nil {
		return nil, err
	}

	// Return nil, if the response contains empty data
	if len(req.Name) < 1 {
		return nil, nil
	}

	return req, nil
}

func getPubSubTopicIamPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getPubSubTopicIamPolicy")

	// Create Service Connection
	service, err := PubsubService(ctx, d)
	if err != nil {
		return nil, err
	}

	topic := h.Item.(*pubsub.Topic)

	req, err := service.Projects.Topics.GetIamPolicy(topic.Name).Do()
	if err != nil {
		// Return nil, if the resource not present
		result := isIgnorableError([]string{"404"})
		if result != nil {
			return nil, nil
		}
		return nil, err
	}

	return req, nil
}

//// TRANSFORM FUNCTIONS

func topicNameToTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	topic := d.HydrateItem.(*pubsub.Topic)
	param := d.Param.(string)

	// get the resource title
	splittedTitle := strings.Split(topic.Name, "/")

	turbotData := map[string]interface{}{
		"Project": splittedTitle[1],
		"Akas":    []string{"gcp://pubsub.googleapis.com/" + topic.Name},
	}

	return turbotData[param], nil
}

func pubsubTopicSelfLink(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*pubsub.Topic)
	selfLink := "https://pubsub.googleapis.com/v1/" + data.Name

	return selfLink, nil
}
