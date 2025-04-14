---
title: "Steampipe Table: gcp_storage_object - Query Google Cloud Storage Objects using SQL"
description: "Allows users to query Google Cloud Storage Objects, specifically the metadata and details of each object stored in Google Cloud Storage."
folder: "Cloud Storage"
---

# Table: gcp_storage_object - Query Google Cloud Storage Objects using SQL

Google Cloud Storage is a scalable, fully-managed, highly reliable, and cost-efficient object/blob store. It is designed to handle data from any size, type, and ingestion speed with a simple and consistent API. In addition to archiving, Google Cloud Storage offers high durability for backup, restore, and disaster recovery use cases.

## Table Usage Guide

The `gcp_storage_object` table provides insights into objects stored within Google Cloud Storage. As a data analyst or storage administrator, explore object-specific details through this table, including metadata, storage class, and associated bucket information. Utilize it to uncover information about objects, such as their size, content type, creation time, and the bucket they are stored in.

## Examples

### Basic info
Explore which objects within a specific Google Cloud Storage bucket are taking up the most space. This allows for efficient space management and can help identify potential areas for data optimization.  

```sql+postgres
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

```sql+sqlite
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
Discover the segments that contain a specific object within a certain bucket in Google Cloud Storage. This can be useful to assess the elements within that object such as its size, storage class, and the time it was created.

```sql+postgres
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

```sql+sqlite
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
Explore storage objects within a specific bucket that share a common prefix, which can be useful for organizing and identifying related files or data sets. This is especially beneficial when dealing with large amounts of data, as it allows you to quickly locate and analyze related objects.

```sql+postgres
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

```sql+sqlite
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
Explore which storage objects within a specific bucket are encrypted using customer-managed keys. This helps in assessing the level of control you have over your data encryption and security.

```sql+postgres
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

```sql+sqlite
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
Explore the total number of objects and their combined size within each storage bucket. This can help you understand your storage usage and manage your resources more effectively.

```sql+postgres
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

```sql+sqlite
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
Explore which members are associated with specific roles for each object in a GCP storage bucket. This can be particularly useful for evaluating access permissions and ensuring appropriate security measures are in place.

```sql+postgres
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

```sql+sqlite
select
  bucket,
  name,
  json_extract(p.value, '$.members') as member,
  json_extract(p.value, '$.role') as role,
  json_extract(p.value, '$.version') as version
from
  gcp_storage_object,
  json_each(iam_policy, '$.bindings') as p
where
  bucket = 'steampipe-test';
```

### List of storage objects whose retention period is less than 7 days
Explore which storage objects have a retention period of less than a week. This is useful for identifying potential data loss risks or ensuring compliance with data retention policies.

```sql+postgres
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

```sql+sqlite
select
  bucket,
  name,
  strftime('%s', retention_expiration_time) - strftime('%s', 'now') as retention_period_secs
from
  gcp_storage_object
where
  strftime('%s', retention_expiration_time) - strftime('%s', 'now') < 604800
  and bucket = 'steampipe-test';
```

### Get accsess controls on each object in a bucket
Explore the access controls assigned to each object within a specific storage bucket to understand the roles and permissions associated with different entities. This can assist in managing security and access within your storage environment.

```sql+postgres
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

```sql+sqlite
select
  bucket,
  name as object_name,
  json_extract(a.value, '$.entity') as entity,
  json_extract(a.value, '$.role') as role,
  json_extract(a.value, '$.email') as email,
  json_extract(a.value, '$.domain') as domain,
  json_extract(a.value, '$.projectTeam') as project_team
from
  gcp_storage_object,
  json_each(acl) as a
where
  bucket = 'steampipe-test';
```