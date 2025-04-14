---
title: "Steampipe Table: gcp_compute_machine_type - Query Google Cloud Platform Compute Machine Types using SQL"
description: "Allows users to query Compute Machine Types in Google Cloud Platform, providing detailed information about available machine types and their specifications."
folder: "Compute"
---

# Table: gcp_compute_machine_type - Query Google Cloud Platform Compute Machine Types using SQL

Google Cloud Platform's Compute Engine provides predefined machine types that you can use to get the right balance of resources for your applications. Each machine type comes with a specific number of vCPUs and amount of memory, and they are grouped into different families based on the ratio of CPU to memory. Machine types help you to meet your specific performance and cost requirements.

## Table Usage Guide

The `gcp_compute_machine_type` table provides insights into the available machine types within Google Cloud Platform's Compute Engine. As a cloud architect or DevOps engineer, you can explore machine type-specific details through this table, including vCPU count, memory, and associated metadata. Utilize it to understand the specifications of each machine type, aiding in the selection of the most suitable machine type for your applications based on performance requirements and cost efficiency.

## Examples

### Basic info
Assess the elements within your Google Cloud Platform to understand the capacity and capabilities of each machine type. This can help optimize resource allocation, by identifying the maximum number of persistent disks and their total size each machine can support.

```sql+postgres
select
  name,
  id,
  description,
  guest_cpus,
  maximum_persistent_disks,
  maximum_persistent_disks_size_gb
from
  gcp_compute_machine_type;
```

```sql+sqlite
select
  name,
  id,
  description,
  guest_cpus,
  maximum_persistent_disks,
  maximum_persistent_disks_size_gb
from
  gcp_compute_machine_type;
```

### List machine types with more than 48 cores
Determine the areas in which machine types have high processing power, specifically those with more than 48 cores. This can help in identifying high-performance options for resource-intensive applications.

```sql+postgres
select
  name,
  id,
  description,
  guest_cpus
from
  gcp_compute_machine_type
where
  guest_cpus >= 48;
```

```sql+sqlite
select
  name,
  id,
  description,
  guest_cpus
from
  gcp_compute_machine_type
where
  guest_cpus >= 48;
```

### List machine types with shared CPUs
Determine the types of machines that utilize shared CPUs to optimize resource allocation and enhance performance efficiency.

```sql+postgres
select
  name,
  id,
  is_shared_cpu
from
  gcp_compute_machine_type
where
  is_shared_cpu;
```

```sql+sqlite
select
  name,
  id,
  is_shared_cpu
from
  gcp_compute_machine_type
where
  is_shared_cpu = 1;
```

### Get accelerator configurations assigned to each machine type
Analyze machine configurations to understand the number and type of accelerators assigned to each. This is beneficial in optimizing resource allocation and performance in a GCP compute environment.

```sql+postgres
select
  name,
  id,
  a -> 'guestAcceleratorCount' as guest_accelerator_count,
  a ->> 'guestAcceleratorType' as guest_accelerator_type
from
  gcp_compute_machine_type,
  jsonb_array_elements(accelerators) as a;
```

```sql+sqlite
select
  m.name,
  m.id,
  json_extract(a.value, '$.guestAcceleratorCount') as guest_accelerator_count,
  json_extract(a.value, '$.guestAcceleratorType') as guest_accelerator_type
from
  gcp_compute_machine_type as m,
  json_each(accelerators) as a;
```

### Display the categorization of machine types by zone
Explore which machine types are most prevalent in each zone. This can help in understanding the distribution of resources and planning for future infrastructure needs.

```sql+postgres
select
  name,
  zone,
  count(name) as numbers_of_machine_type
from
  gcp_compute_machine_type
group by
  name,
  zone;
```

```sql+sqlite
select
  name,
  zone,
  count(name) as numbers_of_machine_type
from
  gcp_compute_machine_type
group by
  name,
  zone;
```