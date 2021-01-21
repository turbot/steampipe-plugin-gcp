# Table: gcp_pubsub_snapshot

Snapshots are used in for seek operations, which allow you to manage message acknowledgments in bulk.

## Examples

### Basic info

```sql
select
  name,
  split_part(topic, '/', 4) as topic_name,
  expire_time,
  tags
from
  gcp_pubsub_snapshot;
```


### Find pubsub snapshots with policies that grant public access

```sql
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