# Table: gcp_storage_bucket

Storage buckets are the basic containers that hold data. Everything that you store in cloud Storage must be contained in a bucket.

## Examples

### List of buckets where versioning is not enabled

```sql
select
  name,
  location,
  versioning_enabled
from
  gcp_storage_bucket
where
  not versioning_enabled;
```


### List of members and their associated iam roles for the bucket

```sql
select
  name,
  location,
  p -> 'members' as member,
  p ->> 'role' as role
from
  gcp_storage_bucket,
  jsonb_array_elements(iam_policy -> 'bindings') as p;
```


### Lifecycle rule of each storage bucket

```sql
select
  name,
  p -> 'action' ->> 'storageClass'  as storage_class,
  p -> 'action' ->> 'type'  as action_type,
  p -> 'condition' ->> 'age' as age_in_days
from
  gcp_storage_bucket,
  jsonb_array_elements(lifecycle_rules) as p;
```


### List of storage buckets whose retention period is less than 7 days

```sql
select
  name,
  retention_policy ->> 'retentionPeriod' as retention_period
from
  gcp_storage_bucket
where
  retention_policy ->> 'retentionPeriod' < 604800 :: text;
```