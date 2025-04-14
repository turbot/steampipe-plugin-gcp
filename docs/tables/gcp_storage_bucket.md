---
title: "Steampipe Table: gcp_storage_bucket - Query Google Cloud Storage Buckets using SQL"
description: "Allows users to query Google Cloud Storage Buckets, providing detailed information about each bucket's configuration, access controls, and associated metadata."
folder: "Cloud Storage"
---

# Table: gcp_storage_bucket - Query Google Cloud Storage Buckets using SQL

Google Cloud Storage is a service within Google Cloud Platform that provides scalable, durable, and highly available object storage. It offers multiple storage classes, versioning, fine-grained access controls, and other features for managing data. Google Cloud Storage is designed to help organizations of all sizes securely store and retrieve any amount of data at any time.

## Table Usage Guide

The `gcp_storage_bucket` table provides insights into Storage Buckets within Google Cloud Storage. As a Cloud Engineer, explore bucket-specific details through this table, including configurations, access controls, and associated metadata. Utilize it to uncover information about buckets, such as their storage class, location, versioning status, and access control policies.

## Examples

### List of buckets where versioning is not enabled
Discover the segments that have not enabled versioning within their storage buckets. This is particularly useful to identify potential risk areas where data loss could occur due to overwriting or deleting of files.

```sql+postgres
select
  name,
  location,
  versioning_enabled
from
  gcp_storage_bucket
where
  not versioning_enabled;
```

```sql+sqlite
select
  name,
  location,
  versioning_enabled
from
  gcp_storage_bucket
where
  versioning_enabled is not 1;
```

### List of members and their associated iam roles for the bucket
Discover the segments that illustrate the relationship between members and their respective roles for a specific storage bucket in GCP. This could be useful in assessing access permissions and managing security within your cloud storage environment.

```sql+postgres
select
  name,
  location,
  p -> 'members' as member,
  p ->> 'role' as role
from
  gcp_storage_bucket,
  jsonb_array_elements(iam_policy -> 'bindings') as p;
```

```sql+sqlite
select
  name,
  location,
  json_extract(p.value, '$.members') as member,
  json_extract(p.value, '$.role') as role
from
  gcp_storage_bucket,
  json_each(iam_policy, '$.bindings') as p;
```

### Lifecycle rule of each storage bucket
Explore the lifecycle rules of your storage buckets to understand how they're configured. This can help in managing resources more effectively by determining when certain actions, such as transitioning to a different storage class or deleting objects, are set to occur.

```sql+postgres
select
  name,
  p -> 'action' ->> 'storageClass'  as storage_class,
  p -> 'action' ->> 'type'  as action_type,
  p -> 'condition' ->> 'age' as age_in_days
from
  gcp_storage_bucket,
  jsonb_array_elements(lifecycle_rules) as p;
```

```sql+sqlite
select
  name,
  json_extract(p.value, '$.action.storageClass') as storage_class,
  json_extract(p.value, '$.action.type') as action_type,
  json_extract(p.value, '$.condition.age') as age_in_days
from
  gcp_storage_bucket,
  json_each(lifecycle_rules) as p;
```

### List of storage buckets whose retention period is less than 7 days
Explore which storage buckets have a retention period of less than a week. This can be useful in identifying potential data loss risks due to short retention periods.

```sql+postgres
select
  name,
  retention_policy ->> 'retentionPeriod' as retention_period
from
  gcp_storage_bucket
where
  retention_policy ->> 'retentionPeriod' < 604800 :: text;
```

```sql+sqlite
select
  name,
  json_extract(retention_policy, '$.retentionPeriod') as retention_period
from
  gcp_storage_bucket
where
  cast(json_extract(retention_policy, '$.retentionPeriod') as integer) < 604800;
```