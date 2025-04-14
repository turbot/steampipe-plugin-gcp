---
title: "Steampipe Table: gcp_redis_instance - Query Google Cloud Platform Memorystore Redis Instances using SQL"
description: "Allows users to query Memorystore Redis Instances on Google Cloud Platform, providing detailed information about each instance."
folder: "Memorystore"
---

# Table: gcp_redis_instance - Query Google Cloud Platform Memorystore Redis Instances using SQL

Google Cloud Platform's Memorystore Redis service is a fully managed service that powers applications with low-latency data access. It provides secure and highly available Redis instances while Google handles all the underlying infrastructure for you. Redis instances are ideal for use cases such as caching, session storage, gaming leaderboards, real-time analytics, and queueing.

## Table Usage Guide

The `gcp_redis_instance` table provides insights into Memorystore Redis Instances within Google Cloud Platform. As a DevOps engineer, you can explore instance-specific details through this table, including the instance's ID, name, region, and current status. Utilize it to monitor and manage your Redis instances, ensuring they are configured correctly and running efficiently.

## Examples

### Basic info
Explore which Google Cloud Platform (GCP) Memorystore Redis instances have been created, along with their display names, creation times, locations, memory sizes, and reserved IP ranges. This is useful for gaining insights into your GCP Memorystore Redis instances' configurations and understanding how your resources are being utilized.

```sql+postgres
select
  name,
  display_name,
  create_time,
  location_id,
  memory_size_gb,
  reserved_ip_range
from
  gcp_redis_instance;
```

```sql+sqlite
select
  name,
  display_name,
  create_time,
  location_id,
  memory_size_gb,
  reserved_ip_range
from
  gcp_redis_instance;
```

### List instances that have authentication enabled
Discover the segments that have enabled authentication to enhance security measures and maintain data privacy within your GCP Memorystore Redis instances. This can be particularly useful in identifying potential vulnerabilities and ensuring compliance with best practices.

```sql+postgres
select
  name,
  display_name,
  create_time,
  location_id,
  memory_size_gb,
  reserved_ip_range
from
  gcp_redis_instance
where
  auth_enabled;
```

```sql+sqlite
select
  name,
  display_name,
  create_time,
  location_id,
  memory_size_gb,
  reserved_ip_range
from
  gcp_redis_instance
where
  auth_enabled = 1;
```

### List instances created in the last 7 days
Discover the segments that have been newly added within the past week. This is beneficial in monitoring the growth and changes in your database over a short period.

```sql+postgres
select
  name,
  display_name,
  create_time,
  location_id,
  memory_size_gb,
  reserved_ip_range
from
  gcp_redis_instance
where
  create_time >= current_timestamp - interval '7 days';
```

```sql+sqlite
select
  name,
  display_name,
  create_time,
  location_id,
  memory_size_gb,
  reserved_ip_range
from
  gcp_redis_instance
where
  create_time >= datetime('now', '-7 days');
```

### List the node details of each instance
Gain insights into the specific details of each instance node in your GCP Memorystore Redis database, such as the creation time and location. This can be particularly useful for troubleshooting or for optimizing your database's performance and security.

```sql+postgres
select
  name,
  display_name,
  create_time,
  location_id,
  jsonb_pretty(nodes) as instance_nodes
from
  gcp_redis_instance
where
  name = 'instance-test'
  and location_id = 'europe-west3-c';
```

```sql+sqlite
select
  name,
  display_name,
  create_time,
  location_id,
  nodes as instance_nodes
from
  gcp_redis_instance
where
  name = 'instance-test'
  and location_id = 'europe-west3-c';
```

### List instances encrypted with customer-managed keys
Discover the segments that utilize customer-managed encryption for their instances. This is useful to assess the security measures and determine areas that might need additional protection.

```sql+postgres
select
  name,
  display_name,
  create_time,
  location_id,
  memory_size_gb,
  reserved_ip_range
from
  gcp_redis_instance
where
  customer_managed_key is not null;
```

```sql+sqlite
select
  name,
  display_name,
  create_time,
  location_id,
  memory_size_gb,
  reserved_ip_range
from
  gcp_redis_instance
where
  customer_managed_key is not null;
```

### List instances that have in-transit encryption disabled
Discover the segments where in-transit encryption is disabled in instances. This is particularly useful in identifying potential security risks and ensuring data protection standards are maintained.

```sql+postgres
select
  name,
  display_name,
  create_time,
  location_id,
  memory_size_gb,
  reserved_ip_range
from
  gcp_redis_instance
where
  transit_encryption_mode != 1;
```

```sql+sqlite
select
  name,
  display_name,
  create_time,
  location_id,
  memory_size_gb,
  reserved_ip_range
from
  gcp_redis_instance
where
  transit_encryption_mode != 1;
```

### List the maintenance details of instances
Explore the maintenance characteristics of your instances to identify when and how often maintenance is performed, as well as the versions available for maintenance. This can help you manage and plan your maintenance activities more effectively.

```sql+postgres
select
  name,
  display_name,
  create_time,
  location_id,
  maintenance_policy,
  maintenance_schedule,
  maintenance_version,
  available_maintenance_versions
from
  gcp_redis_instance;
```

```sql+sqlite
select
  name,
  display_name,
  create_time,
  location_id,
  maintenance_policy,
  maintenance_schedule,
  maintenance_version,
  available_maintenance_versions
from
  gcp_redis_instance;
```

### List instances with direct peering access
Explore which instances have direct peering access in order to better manage your network and ensure secure connections. This can be especially useful for maintaining optimal performance and security in your GCP Memorystore Redis instances.

```sql+postgres
select
  name,
  display_name,
  create_time,
  location_id,
  memory_size_gb,
  reserved_ip_range
from
  gcp_redis_instance
where
  connect_mode = 1;
```

```sql+sqlite
select
  name,
  display_name,
  create_time,
  location_id,
  memory_size_gb,
  reserved_ip_range
from
  gcp_redis_instance
where
  connect_mode = 1;
```
