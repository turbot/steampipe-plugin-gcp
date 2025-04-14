---
title: "Steampipe Table: gcp_compute_node_group - Query Google Cloud Platform Compute Node Groups using SQL"
description: "Allows users to query Compute Node Groups in Google Cloud Platform, providing insights into the configuration, status, and metadata of node groups."
folder: "Compute"
---

# Table: gcp_compute_node_group - Query Google Cloud Platform Compute Node Groups using SQL

A Compute Node Group is a resource in Google Cloud Platform's Compute Engine service. It represents a group of Compute Engine instances that are part of the same node group, enabling you to manage, scale, and deploy instances in a unified way. Node Groups are typically used to manage homogeneous instances that need to be created and managed as a group.

## Table Usage Guide

The `gcp_compute_node_group` table provides insights into Compute Node Groups within Google Cloud Platform's Compute Engine service. As a cloud engineer or system administrator, you can explore details about each node group, including its configuration, operational status, and associated metadata. This table is particularly useful for managing and monitoring groups of instances, ensuring they are correctly configured and performing as expected.

## Examples

### Node group basic info
Explore the status and size of your Google Cloud Platform compute node groups to understand their current state and capacity. This can help manage resources effectively and optimize cloud infrastructure.

```sql+postgres
select
  name,
  status,
  size,
  self_link
from
  gcp_compute_node_group;
```

```sql+sqlite
select
  name,
  status,
  size,
  self_link
from
  gcp_compute_node_group;
```

### List of node groups where the autoscaler is not enabled
Explore which node groups within your Google Cloud Platform's Compute Engine have not enabled the autoscaler. This is useful to ensure optimal resource allocation and cost efficiency by automatically adjusting the number of nodes based on the workload.

```sql+postgres
select
  name,
  id,
  status,
  autoscaling_policy_mode
from
  gcp_compute_node_group
where
  autoscaling_policy_mode <> 'ON';
```

```sql+sqlite
select
  name,
  id,
  status,
  autoscaling_policy_mode
from
  gcp_compute_node_group
where
  autoscaling_policy_mode <> 'ON';
```

### List of node groups with default maintenance settings
Analyze the settings to understand which node groups are operating with default maintenance policies in your Google Cloud Platform compute engine. This is useful to ensure that your nodes are being maintained according to your preferred standards.

```sql+postgres
select
  name,
  id,
  status,
  autoscaling_policy_mode
from
  gcp_compute_node_group
where
  maintenance_policy = 'DEFAULT';
```

```sql+sqlite
select
  name,
  id,
  status,
  autoscaling_policy_mode
from
  gcp_compute_node_group
where
  maintenance_policy = 'DEFAULT';
```

### List node types for node groups
Analyze the settings to understand the different types of nodes within each group in your Google Cloud Platform's compute engine. This is useful for managing resources and optimizing your cloud infrastructure.

```sql+postgres
select
  g.name,
  g.id,
  g.location,
  t.node_type
from
  gcp_compute_node_group as g,
  gcp_compute_node_template as t
where
  g.node_template = t.self_link;
```

```sql+sqlite
select
  g.name,
  g.id,
  g.location,
  t.node_type
from
  gcp_compute_node_group as g
join
  gcp_compute_node_template as t
on
  g.node_template = t.self_link;
```