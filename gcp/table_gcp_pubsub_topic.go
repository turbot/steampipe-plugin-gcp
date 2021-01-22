package gcp

import (
	"context"
	"os"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"google.golang.org/api/pubsub/v1"
)

func tableGcpPubSubTopic(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_pubsub_topic",
		Description: "GCP Pub/Sub Topic",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getPubSubTopic,
		},
		List: &plugin.ListConfig{
			Hydrate: listPubSubTopics,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the topic.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(topicNameToTurbotData, "Name"),
			},
			{
				Name:        "kms_key_name",
				Description: "The resource name of the Cloud KMS CryptoKey to be used to protect access to messages published on this topic.",
				Type:        proto.ColumnType_STRING,
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
				Name:        "tags",
				Description: "A map of tags for the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Labels"),
			},
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(topicNameToTurbotData, "Title"),
			},
			{
				Name:        "akas",
				Description: "Array of globally unique identifier strings (also known as) for the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(topicNameToTurbotData, "Akas"),
			},
			{
				Name:        "project",
				Description: "The Google Project in which the resource is located",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(topicNameToTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listPubSubTopics(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	service, err := pubsub.NewService(ctx)
	if err != nil {
		return nil, err
	}

	project := os.Getenv("GCP_PROJECT")
	resp := service.Projects.Topics.List("projects/" + project)
	if err := resp.Pages(ctx, func(page *pubsub.ListTopicsResponse) error {
		for _, topic := range page.Topics {
			d.StreamListItem(ctx, topic)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getPubSubTopic(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	service, err := pubsub.NewService(ctx)
	if err != nil {
		return nil, err
	}

	name := d.KeyColumnQuals["name"].GetStringValue()
	project := os.Getenv("GCP_PROJECT")

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
	service, err := pubsub.NewService(ctx)
	if err != nil {
		return nil, err
	}

	topic := h.Item.(*pubsub.Topic)

	req, err := service.Projects.Topics.GetIamPolicy(topic.Name).Do()
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

func topicNameToTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	topic := d.HydrateItem.(*pubsub.Topic)
	param := d.Param.(string)

	// get the resource title
	splittedTitle := strings.Split(topic.Name, "/")

	turbotData := map[string]interface{}{
		"Project": splittedTitle[1],
		"Name":    splittedTitle[len(splittedTitle)-1],
		"Title":   splittedTitle[len(splittedTitle)-1],
		"Akas":    []string{"gcp://pubsub.googleapis.com/" + topic.Name},
	}

	return turbotData[param], nil
}
