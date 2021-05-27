# Table: gcp_kms_key

A Cloud KMS key is a named object containing one or more key versions, along with metadata for the key. A key exists on exactly one key ring tied to a specific location.

## Examples

### Basic info

```sql
select
  name,
  create_time,
  rotation_period
from
  gcp_kms_key;
```

### List keys older than 30 days

```sql
select
  name,
  create_time,
  rotation_period
from
  gcp_kms_key
where
  create_time <= (current_date - interval '30' day)
order by
  create_time;
```

### List keys with rotation period more than 7776000s (90 days)

```sql
select
  name,
  create_time,
  rotation_period
from
  gcp_kms_key
where
  split_part(rotation_period, 's', 1) :: int > 7776000;
```

### List publicly accessible keys

```sql
select
  distinct name,
  key_ring_name,
  location
from
  gcp_kms_key,
  jsonb_array_elements(iam_policy -> 'bindings') as b
where
  b -> 'members' ?| array['allAuthenticatedUsers', 'allUsers'];
```
