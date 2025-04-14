---
title: "Steampipe Table: gcp_pubsub_snapshot - Query Google Cloud Pub/Sub Snapshots using SQL"
description: "Allows users to query Google Cloud Pub/Sub Snapshots, providing detailed information about each snapshot in the Google Cloud project."
folder: "Pub/Sub"
---

# Table: gcp_pubsub_snapshot - Query Google Cloud Pub/Sub Snapshots using SQL

Google Cloud Pub/Sub is a messaging service that allows you to send and receive messages between independent applications. A Pub/Sub Snapshot is a point-in-time capture of the message acknowledgment state of a subscription. Snapshots can be used to seek a subscription to a time in the past or to a different subscription's time in the future.

## Table Usage Guide

The `gcp_pubsub_snapshot` table provides insights into Google Cloud Pub/Sub Snapshots within a Google Cloud project. As a DevOps engineer, explore snapshot-specific details through this table, including the snapshot's name, topic, and expiration time. Utilize it to uncover information about snapshots, such as their associated metadata and the state of message acknowledgment at the time the snapshot was created.

## Examples

### Basic info
Explore snapshots in your Google Cloud Pub/Sub service to identify their names, associated topics, expiration times, and any applied tags. This can help you manage and organize your snapshots more effectively.

```sql+postgres
select
  name,
  topic_name,
  expire_time,
  tags
from
  gcp_pubsub_snapshot;
```

```sql+sqlite
select
  name,
  topic_name,
  expire_time,
  tags
from
  gcp_pubsub_snapshot;
```

### Find pubsub snapshots with policies that grant public access
Determine the areas in which public access is granted to pubsub snapshots. This query is useful in identifying potential security risks by pinpointing which snapshots have policies that allow public access.

```sql+postgres
select
  name,
  split_part(s ->> 'role', '/', 2) as role,
  entity
from
  gcp_pubsub_snapshot,
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
  gcp_pubsub_snapshot g,
  json_each(g.iam_policy, '$.bindings') as s,
  json_each(json_extract(s.value, '$.members')) as e
where
  e.value = 'allUsers'
  or e.value = 'allAuthenticatedUsers';
```