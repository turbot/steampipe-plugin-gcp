---
title: "Steampipe Table: gcp_compute_disk - Query Google Cloud Compute Engine Disks using SQL"
description: "Allows users to query Google Cloud Compute Engine Disks, specifically providing detailed information about each disk, including its size, type, and associated instances."
folder: "Compute"
---

# Table: gcp_compute_disk - Query Google Cloud Compute Engine Disks using SQL

Google Cloud Compute Engine Disks are persistent, high-performance block storage for Google Cloud's Virtual Machines (VMs). They are used to store data and serve as the primary storage for data used by VMs. These disks are automatically encrypted, durable, and offer up to 64 TB of space.

## Table Usage Guide

The `gcp_compute_disk` table provides insights into disks within Google Cloud Compute Engine. As a system administrator, you can explore disk-specific details through this table, including their sizes, types, and associated instances. Utilize it to monitor and manage your storage resources effectively, ensuring optimal performance and cost-efficiency.

## Examples

### Basic info
Explore which Google Cloud Platform (GCP) compute disks are being used, their locations, and their respective sizes. This information can be beneficial for managing storage resources and optimizing costs.

```sql+postgres
select
  name,
  id,
  size_gb as disk_size_in_gb,
  type_name,
  zone_name,
  region_name,
  location_type
from
  gcp_compute_disk;
```

```sql+sqlite
select
  name,
  id,
  size_gb as disk_size_in_gb,
  type_name,
  zone_name,
  region_name,
  location_type
from
  gcp_compute_disk;
```

### List disks encrypted with Google-managed key
Explore which disks are encrypted using a Google-managed key to ensure compliance with your organization's data security policies. This can help in identifying potential security vulnerabilities and maintaining data privacy standards.

```sql+postgres
select
  name,
  id,
  zone_name,
  disk_encryption_key_type
from
  gcp_compute_disk
where
  disk_encryption_key_type = 'Google managed';
```

```sql+sqlite
select
  name,
  id,
  zone_name,
  disk_encryption_key_type
from
  gcp_compute_disk
where
  disk_encryption_key_type = 'Google managed';
```

### List disks that are not in use
Discover the segments that include unused disks in your Google Cloud Platform compute disk storage. This can be beneficial in identifying potential areas for cost optimization and resource management.

```sql+postgres
select
  name,
  id,
  users
from
  gcp_compute_disk
where
  users is null;
```

```sql+sqlite
select
  name,
  id,
  users
from
  gcp_compute_disk
where
  users is null;
```

### List regional disks
Explore which disks are regionally located in your Google Cloud Platform's compute engine. This is useful for understanding the distribution of your resources and ensuring data is stored in the appropriate geographical areas.

```sql+postgres
select
  name,
  region_name
from
  gcp_compute_disk
where
  location_type = 'REGIONAL';
```

```sql+sqlite
select
  name,
  region_name
from
  gcp_compute_disk
where
  location_type = 'REGIONAL';
```

### Count the number of disks per availability zone
Analyze the distribution of your storage resources by determining the total number of disks available in each zone. This information can be utilized to efficiently manage and balance your storage resources across different zones.

```sql+postgres
select
  zone_name,
  count(*)
from
  gcp_compute_disk
group by
  zone_name
order by
  count desc;
```

```sql+sqlite
select
  zone_name,
  count(*)
from
  gcp_compute_disk
group by
  zone_name
order by
  count(*) desc;
```

### List disks ordered by size
Analyze your Google Cloud Platform's compute disk storage to understand which disks are consuming the most space. This can help manage storage efficiently by identifying disks that may need to be resized or cleaned up.

```sql+postgres
select
  name,
  size_gb
from
  gcp_compute_disk
order by
  size_gb desc;
```

```sql+sqlite
select
  name,
  size_gb
from
  gcp_compute_disk
order by
  size_gb desc;
```