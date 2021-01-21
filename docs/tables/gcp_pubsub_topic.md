# Table: gcp_pubsub_topic

Pub/Sub is an asynchronous messaging service that decouples services that produce events from services that process events.

## Examples

### List of pubsub topics which are not encrypted

```sql
select
  name,
  kms_key_name
from
  gcp_pubsub_topic
where
  kms_key_name = '';
```


### List of regions which are allowed in message storage policy for each topic

```sql
select
  name,
  jsonb_array_elements_text(
    message_storage_policy_allowed_persistence_regions
  )
from
  gcp_pubsub_topic;
```


### Find topics with policies that grant public access

```sql
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