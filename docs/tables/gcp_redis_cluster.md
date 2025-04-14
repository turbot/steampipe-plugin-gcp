---
title: "Steampipe Table: gcp_redis_cluster - Query Google Cloud Platform Memorystore Redis Clusters using SQL"
description: "Allows users to query Memorystore Redis Clusters on Google Cloud Platform, providing detailed information about each cluster."
folder: "Memorystore"
---

# Table: gcp_redis_cluster - Query Google Cloud Platform Memorystore Redis Clusters using SQL

Google Cloud Platform's Memorystore Redis Cluster service is a fully managed service that powers applications with low-latency data access. It provides secure and highly available Redis clusters while Google handles all the underlying infrastructure for you. Redis clusters are ideal for use cases such as caching, session storage, gaming leaderboards, real-time analytics, and queueing.

## Table Usage Guide

The `gcp_redis_cluster` table provides insights into Memorystore Redis Clusters within Google Cloud Platform. As a DevOps engineer, you can explore cluster-specific details through this table, including the cluster's ID, name, region, and current status. Utilize it to monitor and manage your Redis clusters, ensuring they are configured correctly and running efficiently.

## Examples

### Basic info
Explore which Google Cloud Platform (GCP) Memorystore Redis clusters have been created, along with their creation times, locations, memory sizes, and IP addresses. This is useful for gaining insights into your GCP Memorystore Redis clusters' configurations and understanding how your resources are being utilized.

```sql+postgres
select
  name,
  create_time,
  location,
  size_gb,
  precise_size_gb,
  psc_connections[0] ->> 'address' as address
from
  gcp_redis_cluster;
```

```sql+sqlite
select
  name,
  create_time,
  location,
  size_gb,
  precise_size_gb,
  psc_connections[0] -> 'address' as address
from
  gcp_redis_cluster;
```

### List clusters that have IAM authorization enabled
Discover the segments that have enabled IAM authorization to enhance security measures and maintain data privacy within your GCP Memorystore Redis clusters. This can be particularly useful in identifying potential vulnerabilities and ensuring compliance with best practices.

```sql+postgres
select
  name,
  create_time,
  location,
  psc_connections[0] ->> 'address' as address
from
  gcp_redis_cluster
where
  authorization_mode = 1;
```

```sql+sqlite
select
  name,
  create_time,
  location,
  psc_connections[0] -> 'address' as address
from
  gcp_redis_cluster
where
  authorization_mode = 1;
```

### List clusters created in the last 7 days
Discover the segments that have been newly added within the past week. This is beneficial in monitoring the growth and changes in your database over a short period.

```sql+postgres
select
  name,
  create_time,
  location,
  psc_connections[0] ->> 'address' as address
from
  gcp_redis_cluster
where
  create_time >= current_timestamp - interval '7 days';
```

```sql+sqlite
select
  name,
  create_time,
  location,
  psc_connections[0] -> 'address' as address
from
  gcp_redis_cluster
where
  create_time >= datetime('now', '-7 days');
```

### Get node details of each cluster
Gain insights into the specific details of each cluster node in your GCP Memorystore Redis database, such as the creation time and location. This can be particularly useful for troubleshooting or for optimizing your database's performance and security.

```sql+postgres
select
  name,
  create_time,
  location,
  node_type,
  size_gb,
  replica_count,
  shard_count
from
  gcp_redis_cluster
where
  name = 'cluster-test'
  and location = 'europe-west9';
```

```sql+sqlite
select
  name,
  create_time,
  location,
  node_type,
  size_gb,
  replica_count,
  shard_count
from
  gcp_redis_cluster
where
  name = 'cluster-test'
  and location = 'europe-west9';
```

### List clusters that have in-transit encryption disabled
Discover the segments where in-transit encryption is disabled in clusters. This is particularly useful in identifying potential security risks and ensuring data protection standards are maintained.

```sql+postgres
select
  name,
  create_time,
  location,
  psc_connections[0] ->> 'address' as address
from
  gcp_redis_cluster
where
  transit_encryption_mode != 2;
```

```sql+sqlite
select
  name,
  create_time,
  location,
  psc_connections[0] -> 'address' as address
from
  gcp_redis_cluster
where
  transit_encryption_mode != 2;
```

<!--
FIXME: this is missing from the Go SDK
https://github.com/googleapis/google-cloud-go/issues/11061

### List the maintenance details of clusters
Explore the maintenance characteristics of your clusters to identify when and how often maintenance is performed, as well as the versions available for maintenance. This can help you manage and plan your maintenance activities more effectively.

```sql+postgres
select
  name,
  create_time,
  location,
  maintenance_policy,
  maintenance_schedule
from
  gcp_redis_cluster;
```

```sql+sqlite
select
  name,
  create_time,
  location,
  maintenance_policy,
  maintenance_schedule
from
  gcp_redis_cluster;
```
-->
