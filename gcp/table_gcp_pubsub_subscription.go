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

func tableGcpPubSubSubscription(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_pubsub_subscription",
		Description: "GCP Pub/Sub Subscription",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getPubSubSubscription,
			Tags:       map[string]string{"service": "pubsub", "action": "subscriptions.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listPubSubSubscription,
			Tags:    map[string]string{"service": "pubsub", "action": "subscriptions.list"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getPubSubSubscriptionIamPolicy,
				Tags: map[string]string{"service": "pubsub", "action": "subscriptions.getIamPolicy"},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the subscription.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(lastPathElement),
			},
			{
				Name:        "topic",
				Description: "The name of the topic from which this subscription is receiving messages.",
				Type:        proto.ColumnType_STRING,
			},
			// topic_name is a simpler view of the topic, without the full path
			{
				Name:        "topic_name",
				Description: "The name of the topic from which this subscription is receiving messages.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Topic").Transform(lastPathElement),
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(pubsubSubscriptionSelfLink),
			},
			{
				Name:        "filter",
				Description: "An expression written in the Pub/Sub [filter language](https://cloud.google.com/pubsub/docs/filtering). If non-empty, then only `PubsubMessage`s whose `attributes` field matches the filter are delivered on this subscription. If empty, then no messages are filtered out.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ack_deadline_seconds",
				Description: "The approximate amount of time (on a best-effort basis) Pub/Sub waits for the subscriber to acknowledge receipt before resending the message.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "message_retention_duration",
				Description: "How long to retain unacknowledged messages in the subscription's backlog, from the moment a message is published. If `retain_acked_messages` is true, then this also configures the retention of acknowledged messages, and thus configures how far back in time a `Seek` can be done. Defaults to 7 days. Cannot be more than 7 days or less than 10 minutes.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "retain_acked_messages",
				Description: "Indicates whether to retain acknowledged messages. If true, then messages are not expunged from the subscription's backlog, even if they are acknowledged, until they fall out of the `message_retention_duration` window.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "dead_letter_policy_topic",
				Description: "The name of the topic to which dead letter messages should be published.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DeadLetterPolicy.DeadLetterTopic"),
			},
			{
				Name:        "dead_letter_policy_max_delivery_attempts",
				Description: "The maximum number of delivery attempts for any message. The value must be between 5 and 100.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("DeadLetterPolicy.MaxDeliveryAttempts"),
			},
			{
				Name:        "enable_message_ordering",
				Description: "If true, messages published with the same `ordering_key` in `PubsubMessage` will be delivered to the subscribers in the order in which they are received by the Pub/Sub system. Otherwise, they may be delivered in any order.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "expiration_policy_ttl",
				Description: "Specifies the \"time-to-live\" duration for an associated resource. The resource expires if it is not active for a period of `ttl`. The definition of \"activity\" depends on the type of the associated resource. The minimum and maximum allowed values for `ttl` depend on the type of the associated resource, as well. If `ttl` is not set, the associated resource never expires.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ExpirationPolicy.Ttl"),
			},
			{
				Name:        "push_config_endpoint",
				Description: "A URL locating the endpoint to which messages should be pushed. For example, a Webhook endpoint might use `https://example.com/push`",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PushConfig.PushEndpoint"),
			},
			{
				Name:        "push_config_attributes",
				Description: "Endpoint configuration attributes that can be used to control different aspects of the message delivery. The only currently supported attribute is \"x-goog-version\". This attribute indicates the version of the data expected by the endpoint. This controls the shape of the pushed message (i.e., its fields and metadata).",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PushConfig.Attributes"),
			},
			{
				Name:        "push_config_oidc_token_service_account_email",
				Description: "Service account email to be used for generating the OIDC token.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PushConfig.OidcToken.ServiceAccountEmail"),
			},
			{
				Name:        "push_config_oidc_token_audience",
				Description: "Audience to be used when generating OIDC token. The audience claim identifies the recipients that the JWT is intended for. The audience value is a single case-sensitive string.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PushConfig.OidcToken.Audience"),
			},
			{
				Name:        "retry_policy_maximum_backoff",
				Description: "The maximum delay between consecutive deliveries of a given message. Value should be between 0 and 600 seconds. Defaults to 600 seconds.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RetryPolicy.MaximumBackoff"),
			},
			{
				Name:        "retry_policy_minimum_backoff",
				Description: "The minimum delay between consecutive deliveries of a given message. Value should be between 0 and 600 seconds. Defaults to 10 seconds.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RetryPolicy.MinimumBackoff"),
			},
			{
				Name:        "iam_policy",
				Description: "An Identity and Access Management (IAM) policy, which specifies access controls for Google Cloud resources. A `Policy` is a collection of `bindings`. A `binding` binds one or more `members` to a single `role`. Members can be user accounts, service accounts, Google groups, and domains (such as G Suite). A `role` is a named list of permissions; each `role` can be an IAM predefined role or a user-created custom role. For some types of Google Cloud resources, a `binding` can also specify a `condition`, which is a logical expression that allows access to a resource only if the expression evaluates to `true`.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getPubSubSubscriptionIamPolicy,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "labels",
				Description: "A set of labels attached with the subscription.",
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
				Transform:   transform.FromP(subscriptionNameToTurbotData, "Akas"),
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
				Transform:   transform.FromP(subscriptionNameToTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listPubSubSubscription(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

	resp := service.Projects.Subscriptions.List("projects/" + project).PageSize(*pageSize)
	if err := resp.Pages(ctx, func(page *pubsub.ListSubscriptionsResponse) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, sub := range page.Subscriptions {
			d.StreamListItem(ctx, sub)

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

func getPubSubSubscription(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getPubSubSubscription")

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

	req, err := service.Projects.Subscriptions.Get("projects/" + project + "/subscriptions/" + name).Do()
	if err != nil {
		return nil, err
	}

	// Return nil, if the response contains empty data
	if len(req.Name) < 1 {
		return nil, nil
	}

	return req, nil
}

func getPubSubSubscriptionIamPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getPubSubSubscriptionIamPolicy")

	// Create Service Connection
	service, err := PubsubService(ctx, d)
	if err != nil {
		return nil, err
	}

	resource := h.Item.(*pubsub.Subscription)
	req, err := service.Projects.Subscriptions.GetIamPolicy(resource.Name).Do()
	if err != nil {
		return nil, err
	}

	return req, nil
}

//// TRANSFORM FUNCTIONS

func subscriptionNameToTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	subscription := d.HydrateItem.(*pubsub.Subscription)
	param := d.Param.(string)

	splittedTitle := strings.Split(subscription.Name, "/")

	turbotData := map[string]interface{}{
		"Project": splittedTitle[1],
		"Akas":    []string{"gcp://pubsub.googleapis.com/" + subscription.Name},
	}

	return turbotData[param], nil
}

func pubsubSubscriptionSelfLink(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*pubsub.Subscription)
	selfLink := "https://pubsub.googleapis.com/v1/" + data.Name

	return selfLink, nil
}
