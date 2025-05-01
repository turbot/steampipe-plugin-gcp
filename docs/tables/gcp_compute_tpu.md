---
title: "Steampipe Table: gcp_compute_tpu - Query Google Cloud TPUs using SQL"
description: "Allows users to query Google Cloud TPUs, providing insights into TPU node configurations, states, and associated metadata."
folder: "Compute"
---

# Table: gcp_compute_tpu - Query Google Cloud TPUs using SQL

Google Cloud TPUs (Tensor Processing Units) are specialized hardware accelerators designed to speed up specific machine learning workloads. They are particularly optimized for TensorFlow, Google's open-source machine learning framework. TPUs can significantly accelerate machine learning training and inference for your applications.

## Table Usage Guide

The `gcp_compute_tpu` table provides insights into TPU nodes within Google Cloud Platform. As a machine learning engineer or data scientist, explore TPU-specific details through this table, including node configurations, operational states, and associated metadata. You can use this table to gather information about your TPU resources, such as:

- Identifying TPU nodes in specific states
- Monitoring TPU node health and performance
- Verifying TPU configurations and network settings
- Tracking TPU resource utilization and scheduling

## Examples

### Basic info
```sql
select
  name,
  id,
  state,
  accelerator_type,
  zone
from
  gcp_compute_tpu;
```

### List TPUs that are not running
```sql
select
  name,
  state,
  health_description,
  zone
from
  gcp_compute_tpu
where
  state != 'READY';
```

### List TPUs by accelerator type
```sql
select
  accelerator_type,
  count(*) as count,
  array_agg(name) as tpu_names
from
  gcp_compute_tpu
group by
  accelerator_type;
```

### Get network configuration details for each TPU
```sql
select
  name,
  network,
  cidr_block,
  network_endpoint
from
  gcp_compute_tpu;
```

### List TPUs with their creation time and runtime version
```sql
select
  name,
  create_time,
  runtime_version,
  zone
from
  gcp_compute_tpu
order by
  create_time;
```

### Get TPUs with their associated service accounts
```sql
select
  name,
  service_account,
  zone
from
  gcp_compute_tpu
where
  service_account is not null;
```

### List TPUs with their tags
```sql
select
  name,
  tags,
  zone
from
  gcp_compute_tpu
where
  tags is not null;
```

### Get TPUs with scheduling configurations
```sql
select
  name,
  scheduling_config,
  zone
from
  gcp_compute_tpu
where
  scheduling_config is not null;
```

### List TPUs with any reported symptoms
```sql
select
  name,
  symptoms,
  health_description,
  zone
from
  gcp_compute_tpu
where
  symptoms is not null;
```

### Get TPU distribution across zones
```sql
select
  zone,
  count(*) as tpu_count,
  array_agg(name) as tpu_names
from
  gcp_compute_tpu
group by
  zone
order by
  tpu_count desc;
``` 