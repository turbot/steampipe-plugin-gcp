---
title: "Steampipe Table: gcp_compute_instance_group_manager - Query Google Cloud Compute Engine Instance Group Managers using SQL"
description: "Allows users to query Google Cloud Compute Engine Instance Group Managers, providing insights into the configuration, status, and properties of these group managers."
folder: "Compute"
---

# Table: gcp_compute_instance_group_manager - Query Google Cloud Compute Engine Instance Group Managers using SQL

Google Cloud Compute Engine Instance Group Managers manage [Managed Instance Groups (MIG)](https://cloud.google.com/compute/docs/instance-groups#managed_instance_groups), and are ideal for highly available applications that require a lot of computing power and need to scale rapidly to meet demand. They offer a range of features including autoscaling, autohealing, regional (multiple zone) deployment, and automatic updating.

## Table Usage Guide

The `gcp_compute_instance_group_manager` table provides insights into instance group managers within Google Cloud Compute Engine. As a system administrator, you can explore group-specific details through this table, including configuration, associated instances, and autoscaling policies. Utilize it to monitor the status of your instance groups, manage load balancing, and plan for capacity adjustments.

## Examples

### Basic Info
Discover the segments of your Google Cloud Platform (GCP) that contain instance group managers, gaining insights into aspects like size and location. This can help in project management and resource allocation within the GCP infrastructure.

```sql+postgres
select
  name,
  description,
  self_link,
  instance_group,
  location,
  akas,
  project
from
  gcp_compute_instance_group_manager;
```

```sql+sqlite
select
  name,
  description,
  self_link,
  instance_group,
  location,
  akas,
  project
from
  gcp_compute_instance_group_manager;
```

### Get instance group details of each instance group manager
Get the size of the instance groups managed by instance group managers.

```sql+postgres
select
  m.name,
  g.name as group_name,
  g.size as group_size
from
  gcp_compute_instance_group_manager as m,
  gcp_compute_instance_group as g
where
  m.instance_group ->> 'name' = g.name;
```

```sql+sqlite
select
  m.name,
  g.name as group_name,
  g.size as group_size
from
  gcp_compute_instance_group_manager as m,
  gcp_compute_instance_group as g
where
  m.instance_group -> 'name' = g.name;
```
