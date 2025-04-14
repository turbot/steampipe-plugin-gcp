---
title: "Steampipe Table: gcp_pubsub_subscription - Query Google Cloud Pub/Sub Subscriptions using SQL"
description: "Allows users to query Google Cloud Pub/Sub Subscriptions, specifically providing details on subscription configuration, including topics, acknowledgement deadlines, and retention policies."
folder: "Pub/Sub"
---

# Table: gcp_pubsub_subscription - Query Google Cloud Pub/Sub Subscriptions using SQL

Google Cloud Pub/Sub is a messaging service that allows you to send and receive messages between independent applications. Subscriptions in Pub/Sub represent a pipeline from a topic to a receiving entity. They allow the receiving entity to receive messages from a topic, ensuring reliable delivery of the messages.

## Table Usage Guide

The `gcp_pubsub_subscription` table provides insights into Pub/Sub subscriptions within Google Cloud Platform. As a developer or system administrator, explore subscription-specific details through this table, including associated topics, acknowledgement deadlines, and message retention policies. Utilize it to monitor the configuration and status of your Pub/Sub subscriptions, ensuring reliable message delivery between your applications.

## Examples

### List of pubsub subscriptions which are not configured with dead letter topic
Determine the areas in which pubsub subscriptions are not configured with a dead letter topic, allowing you to pinpoint potential issues in the message delivery process.

```sql+postgres
select
  name,
  topic_name,
  dead_letter_policy_topic
from
  gcp_pubsub_subscription
where
  dead_letter_policy_topic is null;
```

```sql+sqlite
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
Analyze the settings to understand the configuration of your message subscriptions, including message retention duration and delivery attempts. This can help you optimize your message delivery and retention processes for better resource management and efficiency.

```sql+postgres
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

```sql+sqlite
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