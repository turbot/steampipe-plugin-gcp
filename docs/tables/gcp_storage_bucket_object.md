# Table: gcp_storage_bucket_object

The Objects resource represents an object within Cloud Storage. Objects are pieces of data that you have uploaded to Cloud Storage.

## Examples

### Basic info
  
```sql
select
  id,
  name,
  bucket,
  size,
  storage_class,
  time_created
from
  gcp_storage_bucket_object;
```

### List storage objects encrypted with customer managed keys

```sql
select
  id,
  name,
  bucket,
  kms_key_name
from
  gcp_storage_bucket_object
where
  kms_key_name != '';
```

### Get total objects and size of each bucket

```sql
select
  bucket,
  count(*) as total_objects,
  sum(size) as total_size_bytes
from
  gcp_storage_bucket_object
group by
  bucket;
```

### List of members and their associated iam roles for each objects

```sql
select
  bucket,
  name,
  p -> 'members' as member,
  p ->> 'role' as role,
  p ->> 'version' as version
from
  gcp_storage_bucket_object,
  jsonb_array_elements(iam_policy -> 'bindings') as p;
```

### List of storage objects whose retention period is less than 7 days

```sql
select
  bucket,
  name,
  extract(epoch from (retention_expiration_time - current_timestamp)) as retention_period_secs
from
  gcp_storage_bucket_object
where
  extract(epoch from (retention_expiration_time - current_timestamp)) < 604800;
```