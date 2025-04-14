---
title: "Steampipe Table: gcp_compute_target_pool - Query Google Cloud Compute Target Pools using SQL"
description: "Allows users to query Google Cloud Compute Target Pools, providing insights into the state, region, and session affinity of each target pool."
folder: "Compute"
---

# Table: gcp_compute_target_pool - Query Google Cloud Compute Target Pools using SQL

Google Cloud Compute Target Pools are a group of instances located in the same region that are used for forwarding rules. They are used to define where incoming traffic should be directed. The instances within a target pool can be added or removed as per the requirements.

## Table Usage Guide

The `gcp_compute_target_pool` table provides insights into the target pools within Google Cloud Compute Engine. As a network administrator, explore target pool-specific details through this table, including the state, region, and session affinity of each target pool. Utilize it to uncover information about target pools, such as their backup pools, failover ratios, and the health checks they are associated with.

## Examples

### Basic info
Explore which Google Cloud Platform (GCP) compute target pools are available in your environment. This can help in managing load balancing by identifying the specific locations of these resources.

```sql+postgres
select
  name,
  id,
  location
from
  gcp_compute_target_pool;
```

```sql+sqlite
select
  name,
  id,
  location
from
  gcp_compute_target_pool;
```

### List of target pools and attached instances that receives incoming traffic
Explore which target pools and their attached instances are set to receive incoming traffic. This can be used to determine the configuration and traffic management of your network, ensuring optimal performance and security.

```sql+postgres
select
  name,
  id,
  split_part(i, '/', 11) as instance_name
from
  gcp_compute_target_pool,
  jsonb_array_elements_text(instances) as i;
```

```sql+sqlite
Error: SQLite does not support split_part function.
```

### List of Health checks attached to each target pool
Explore the health check mechanisms associated with each target pool to effectively monitor and manage your resources in Google Cloud Platform. This can help you maintain optimal performance and quickly address any issues that arise.

```sql+postgres
select
  name,
  id,
  split_part(h, '/', 10) as health_check
from
  gcp_compute_target_pool,
  jsonb_array_elements_text(health_checks) as h;
```

```sql+sqlite
Error: SQLite does not support split or string_to_array functions.
```