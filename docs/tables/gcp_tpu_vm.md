---
title: "Steampipe Table: gcp_tpu_vm - Query Google Cloud TPU VMs using SQL"
description: "Allows users to query Google Cloud TPU VMs, specifically providing detailed information about TPU node configurations, states, network settings, and associated metadata."
folder: "TPU"
---

# Table: gcp_tpu_vm - Query Google Cloud TPU VMs using SQL

Google Cloud TPU VMs (Tensor Processing Units) are specialized hardware accelerators designed to speed up specific machine learning workloads. They are particularly optimized for TensorFlow, Google's open-source machine learning framework. TPUs provide high-performance, custom-developed ASICs (Application-Specific Integrated Circuits) that can significantly accelerate machine learning training and inference for your applications.

## Table Usage Guide

The `gcp_tpu_vm` table provides insights into TPU nodes within Google Cloud Platform. As a machine learning engineer or data scientist, explore TPU-specific details through this table to monitor and manage your machine learning infrastructure effectively. Use it to track TPU node states, verify configurations, and ensure optimal resource utilization across your machine learning workloads.

## Examples

### Basic info
Explore the fundamental details of your TPU nodes to understand their current state and configuration. This can help in managing and monitoring your machine learning infrastructure effectively.

```sql+postgres
select
  name,
  id,
  state,
  accelerator_type,
  zone
from
  gcp_tpu_vm;
```

```sql+sqlite
select
  name,
  id,
  state,
  accelerator_type,
  zone
from
  gcp_tpu_vm;
```

### List non-operational TPU nodes
Identify TPU nodes that are not in a ready state to help troubleshoot potential issues and maintain optimal resource availability for your machine learning workloads.

```sql+postgres
select
  name,
  state,
  health_description,
  zone
from
  gcp_tpu_vm
where
  state != 'READY';
```

```sql+sqlite
select
  name,
  state,
  health_description,
  zone
from
  gcp_tpu_vm
where
  state != 'READY';
```

### TPU distribution by accelerator type
Analyze the distribution of TPU nodes across different accelerator types to understand your resource allocation and assist with capacity planning for machine learning workloads.

```sql+postgres
select
  accelerator_type,
  count(*) as count,
  array_agg(name) as tpu_names
from
  gcp_tpu_vm
group by
  accelerator_type;
```

```sql+sqlite
select
  accelerator_type,
  count(*) as count,
  group_concat(name) as tpu_names
from
  gcp_tpu_vm
group by
  accelerator_type;
```

### Network configuration analysis
Examine the network settings of your TPU nodes to verify connectivity configurations and ensure proper network security settings are in place.

```sql+postgres
select
  name,
  network,
  cidr_block,
  network_endpoints
from
  gcp_tpu_vm;
```

```sql+sqlite
select
  name,
  network,
  cidr_block,
  network_endpoints
from
  gcp_tpu_vm;
```

### TPU nodes by creation time
Track when your TPU nodes were created to understand resource provisioning patterns and manage lifecycle effectively.

```sql+postgres
select
  name,
  create_time,
  zone
from
  gcp_tpu_vm
order by
  create_time;
```

```sql+sqlite
select
  name,
  create_time,
  zone
from
  gcp_tpu_vm
order by
  create_time;
```

### TPU nodes with service accounts
Identify TPU nodes that have associated service accounts to audit security configurations and ensure appropriate access controls are in place.

```sql+postgres
select
  name,
  service_account,
  zone
from
  gcp_tpu_vm
where
  service_account is not null;
```

```sql+sqlite
select
  name,
  service_account,
  zone
from
  gcp_tpu_vm
where
  service_account is not null;
```

### Tagged TPU nodes
List TPU nodes with associated tags to understand your resource organization and management strategies.

```sql+postgres
select
  name,
  tags,
  zone
from
  gcp_tpu_vm
where
  tags is not null;
```

```sql+sqlite
select
  name,
  tags,
  zone
from
  gcp_tpu_vm
where
  tags is not null;
```

### TPU nodes with scheduling configurations
Analyze TPU nodes with specific scheduling configurations to optimize resource utilization and manage costs effectively.

```sql+postgres
select
  name,
  scheduling_config,
  zone
from
  gcp_tpu_vm
where
  scheduling_config is not null;
```

```sql+sqlite
select
  name,
  scheduling_config,
  zone
from
  gcp_tpu_vm
where
  scheduling_config is not null;
```

### TPU nodes with health issues
Monitor TPU nodes that are reporting symptoms or health issues to maintain optimal performance and reliability of your machine learning infrastructure.

```sql+postgres
select
  name,
  symptoms,
  health_description,
  zone
from
  gcp_tpu_vm
where
  symptoms is not null;
```

```sql+sqlite
select
  name,
  symptoms,
  health_description,
  zone
from
  gcp_tpu_vm
where
  symptoms is not null;
```

### Geographic distribution of TPU nodes
Analyze the distribution of TPU nodes across different zones to understand your regional resource allocation and ensure appropriate availability for your machine learning workloads.

```sql+postgres
select
  zone,
  count(*) as tpu_count,
  array_agg(name) as tpu_names
from
  gcp_tpu_vm
group by
  zone
order by
  tpu_count desc;
```

```sql+sqlite
select
  zone,
  count(*) as tpu_count,
  group_concat(name) as tpu_names
from
  gcp_tpu_vm
group by
  zone
order by
  tpu_count desc;
``` 