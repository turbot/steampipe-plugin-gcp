---
title: "Steampipe Table: gcp_pubsub_topic - Query GCP PubSub Topics using SQL"
description: "Allows users to query GCP PubSub Topics, specifically returning details such as topic name, project ID, and subscription information, providing insights into topic configurations and subscription details."
folder: "Pub/Sub"
---

# Table: gcp_pubsub_topic - Query GCP PubSub Topics using SQL

Google Cloud Pub/Sub is a scalable, durable event ingestion and delivery system that serves as a foundation for real-time analytics and event-driven computing systems. Pub/Sub offers at-least-once message delivery and real-time streaming through a simple and consistent API. It provides strong security and authentication, ensuring that your data is safe and only accessible to authorized services and users.

## Table Usage Guide

The `gcp_pubsub_topic` table provides insights into PubSub Topics within Google Cloud Platform (GCP). As a DevOps engineer, explore topic-specific details through this table, including topic name, project ID, and subscription information. Utilize it to uncover information about topics, such as their configurations, the number of subscriptions, and other associated metadata.

## Examples

### List of pubsub topics which are not encrypted
Discover the segments that have unencrypted pubsub topics in your Google Cloud Platform. This is particularly useful for identifying potential security risks and ensuring all your data is adequately protected.

```sql+postgres
select
  name,
  kms_key_name
from
  gcp_pubsub_topic
where
  kms_key_name = '';
```

```sql+sqlite
select
  name,
  kms_key_name
from
  gcp_pubsub_topic
where
  kms_key_name is null;
```

### List of regions which are allowed in message storage policy for each topic
Determine the areas in which message storage policies are permitted for each topic to manage and streamline your data storage strategy effectively.

```sql+postgres
select
  name,
  jsonb_array_elements_text(
    message_storage_policy_allowed_persistence_regions
  )
from
  gcp_pubsub_topic;
```

```sql+sqlite
select
  name,
  json_each.value
from
  gcp_pubsub_topic,
  json_each(message_storage_policy_allowed_persistence_regions);
```

### Find topics with policies that grant public access
This query allows you to pinpoint specific topics that have policies granting public access. This can be useful for identifying potential security risks and ensuring that sensitive information is adequately protected.

```sql+postgres
select
  name,
  split_part(s ->> 'role', '/', 2) as role,
  entity
from
  gcp_pubsub_topic,
  jsonb_array_elements(iam_policy -> 'bindings') as s,
  jsonb_array_elements_text(s -> 'members') as entity
where
  entity = 'allUsers'
  or entity = 'allAuthenticatedUsers';
```

```sql+sqlite
select
  g.name,
  substr(
    json_extract(s.value, '$.role'),
    instr(json_extract(s.value, '$.role'), '/') + 1
  ) as role,
  e.value as entity
from
  gcp_pubsub_topic g,
  json_each(json_extract(g.iam_policy, '$.bindings')) as s,
  json_each(json_extract(s.value, '$.members')) as e
where
  e.value = 'allUsers'
  or e.value = 'allAuthenticatedUsers';
```