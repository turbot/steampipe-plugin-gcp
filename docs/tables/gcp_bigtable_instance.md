---
title: "Steampipe Table: gcp_bigtable_instance - Query Google Cloud BigTable Instances using SQL"
description: "Allows users to query Google Cloud BigTable Instances, specifically the details about each instance like its ID, name, project, number of nodes, storage type, etc."
folder: "Bigtable"
---

# Table: gcp_bigtable_instance - Query Google Cloud BigTable Instances using SQL

Google Cloud BigTable is a scalable, fully-managed NoSQL wide-column database that is suitable for both real-time access and analytics workloads. It can handle massive workloads, offering low latency and high throughput, and is widely used for applications in ad tech, finance, and IoT among others. The instances in BigTable serve as containers for your data where you can perform operations like creating or deleting tables.

## Table Usage Guide

The `gcp_bigtable_instance` table provides insights into BigTable instances within Google Cloud Platform. As a data engineer or a database administrator, explore instance-specific details through this table, including instance ID, name, project, number of nodes, storage type, and more. Utilize it to uncover information about instances, such as those with specific storage types, the number of nodes in each instance, and the state of the instance.

## Examples

### Basic info
Explore which Google Cloud Bigtable instances are currently active and where they are located to better manage your resources and optimize your database operations.

```sql+postgres
select
  name,
  instance_type,
  state,
  location
from
  gcp_bigtable_instance;
```

```sql+sqlite
select
  name,
  instance_type,
  state,
  location
from
  gcp_bigtable_instance;
```

### Get members and their associated IAM roles for each instance
Discover the segments that include members and their associated roles within each instance. This is useful for understanding the distribution of roles and responsibilities within your Google Cloud Platform Bigtable instances.

```sql+postgres
select
  name,
  location,
  jsonb_array_elements_text(p -> 'members') as member,
  p ->> 'role' as role
from
  gcp_bigtable_instance,
  jsonb_array_elements(iam_policy -> 'bindings') as p;
```

```sql+sqlite
select
  name,
  location,
  json_extract(p.value, '$.members') as member,
  json_extract(p.value, '$.role') as role
from
  gcp_bigtable_instance,
  json_each(iam_policy, '$.bindings') as p;
```

### List instances whose members have Bigtable admin access
Discover the segments that have Bigtable admin access in order to manage and control user access rights effectively. This is useful for maintaining security and ensuring only authorized individuals have administrative privileges.

```sql+postgres
select
  name,
  instance_type,
  jsonb_array_elements_text(i -> 'members') as members,
  i ->> 'role' as role
from
  gcp_bigtable_instance,
  jsonb_array_elements(iam_policy -> 'bindings') as i
where
  i ->> 'role' like '%bigtable.admin';
```

```sql+sqlite
select
  name,
  instance_type,
  json_extract(i.value, '$.members') as members,
  json_extract(i.value, '$.role') as role
from
  gcp_bigtable_instance,
  json_each(iam_policy, '$.bindings') as i
where
  json_extract(i.value, '$.role') like '%bigtable.admin';
```

### Count the number of instances per instance type
Explore the distribution of instances across various types in your Google Cloud Bigtable to better manage resources and optimize performance.

```sql+postgres
select
  instance_type,
  count(name)
from
  gcp_bigtable_instance
group by
  instance_type;
```

```sql+sqlite
select
  instance_type,
  count(name)
from
  gcp_bigtable_instance
group by
  instance_type;
```