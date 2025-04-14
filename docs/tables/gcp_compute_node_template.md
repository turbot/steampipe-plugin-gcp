---
title: "Steampipe Table: gcp_compute_node_template - Query Google Cloud Compute Node Templates using SQL"
description: "Allows users to query Google Cloud Compute Node Templates, providing detailed information about node templates within a GCP project."
folder: "Compute"
---

# Table: gcp_compute_node_template - Query Google Cloud Compute Node Templates using SQL

A Google Cloud Compute Node Template is a resource within Google Cloud's Compute Engine that defines the properties of a node group to be created. Node Templates specify the server configuration of instances in a node group. This includes the type of CPU, amount of memory, and disk size and type.

## Table Usage Guide

The `gcp_compute_node_template` table provides insights into node templates within Google Cloud's Compute Engine. As a Cloud Engineer, explore template-specific details through this table, including CPU and memory configurations, disk size and type, and associated metadata. Utilize it to uncover information about node templates, such as those with specific configurations, the properties of each node template, and the verification of node template properties.

## Examples

### List of n2-node-80-640 type node templates
Determine the areas in which specific node templates of type 'n2-node-80-640' are used in your Google Cloud Platform. This can be helpful to understand the distribution and usage of this particular node type across different locations.

```sql+postgres
select
  name,
  id,
  location,
  node_type
from
  gcp_compute_node_template
where
  node_type = 'n2-node-80-640';
```

```sql+sqlite
select
  name,
  id,
  location,
  node_type
from
  gcp_compute_node_template
where
  node_type = 'n2-node-80-640';
```

### List of node templates where cpu overcommit is enabled
Determine the areas in which CPU overcommit is enabled within node templates, to manage resource allocation effectively and optimize cloud infrastructure performance.

```sql+postgres
select
  name,
  id,
  node_type
from
  gcp_compute_node_template
where
  cpu_overcommit_type = 'ENABLED';
```

```sql+sqlite
select
  name,
  id,
  node_type
from
  gcp_compute_node_template
where
  cpu_overcommit_type = 'ENABLED';
```

### Count of node templates per location
Determine the distribution of node templates across different locations in your Google Cloud Platform. This can help you understand the geographical spread of your compute resources for better resource management and planning.

```sql+postgres
select
  location,
  count(*)
from
  gcp_compute_node_template
group by
  location;
```

```sql+sqlite
select
  location,
  count(*)
from
  gcp_compute_node_template
group by
  location;
```

### Find unused node templates
Explore which node templates in your Google Cloud Platform are not currently being used. This can help in identifying unused resources, potentially leading to cost savings and better resource management.

```sql+postgres
select
  t.name,
  t.id
from
  gcp_compute_node_template as t
left join
  gcp_compute_node_group as g on g.node_template = t.self_link
where
  g is null;
```

```sql+sqlite
select
  t.name,
  t.id
from
  gcp_compute_node_template as t
left join
  gcp_compute_node_group as g on g.node_template = t.self_link
where
  g.node_template is null;
```