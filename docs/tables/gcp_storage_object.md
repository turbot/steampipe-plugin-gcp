# Table: gcp_storage_object

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
  gcp_storage_object
where
  bucket = 'steampipe-test';
```

### Get a specific object in a bucket

```sql
select
  id,
  name,
  bucket,
  size,
  storage_class,
  time_created
from
  gcp_storage_object
where
  bucket = 'steampipe-test'
  and name = 'test/logs/2021/03/01/12/abc.txt';
```

### List storage objects with a prefix in a bucket

```sql
select
  id,
  name,
  bucket,
  size,
  storage_class,
  time_created
from
  gcp_storage_object
where
  bucket = 'steampipe-test'
  and prefix = 'test/logs/2021/03/01/12';
```

### List storage objects encrypted with customer managed keys

```sql
select
  id,
  name,
  bucket,
  kms_key_name
from
  gcp_storage_object
where
  bucket = 'steampipe-test'
  and kms_key_name != '';
```

### Get total objects and size of each bucket

```sql
select
  bucket,
  count(*) as total_objects,
  sum(size) as total_size_bytes
from
  gcp_storage_object o,
  gcp_storage_bucket b
where
  o.bucket = b.name
group by
  bucket;
```

### List of members and their associated IAM roles for each objects

```sql
select
  bucket,
  name,
  p -> 'members' as member,
  p ->> 'role' as role,
  p ->> 'version' as version
from
  gcp_storage_object,
  jsonb_array_elements(iam_policy -> 'bindings') as p
where
  bucket = 'steampipe-test';
```

### List of storage objects whose retention period is less than 7 days

```sql
select
  bucket,
  name,
  extract(epoch from (retention_expiration_time - current_timestamp)) as retention_period_secs
from
  gcp_storage_object
where
  extract(epoch from (retention_expiration_time - current_timestamp)) < 604800
  and bucket = 'steampipe-test';
```

### Get accsess controls on each object in a bucket

```sql
select
  bucket,
  name as object_name,
  a ->> 'entity' as entity,
  a ->> 'role' as role,
  a ->> 'email' as email,
  a ->> 'domain' as domain,
  a ->> 'projectTeam' as project_team
from
  gcp_storage_object,
  jsonb_array_elements(acl) as a
where
  bucket = 'steampipe-test';
```
