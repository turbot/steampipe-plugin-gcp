# Table: gcp_pubsub_subscription

The subscription connects the topic to a subscriber application that receives and processes messages published to the topic.

## Examples

### List of pubsub subscriptions which are not configured with dead letter topic

```sql
select
  name,
  topic_name,
  dead_letter_policy_topic
from
  gcp_pubsub_subscription
where
  dead_letter_policy_topic is null;
```


### Message configuration details for the subscriptions

```sql
select
  name,
  topic_name,
  ack_deadline_seconds,
  message_retention_duration,
  retain_acked_messages,
  dead_letter_policy_topic,
  dead_letter_policy_max_delivery_attempts,
  enable_message_ordering
from
  gcp_pubsub_subscription;
```
