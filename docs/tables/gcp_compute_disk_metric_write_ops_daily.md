---
title: "Steampipe Table: gcp_compute_disk_metric_write_ops_daily - Query GCP Compute Engine Disk Metrics using SQL"
description: "Allows users to query Compute Engine Disk Metrics in GCP, specifically the daily write operations count, providing insights into disk usage and potential efficiency improvements."
folder: "Compute"
---

# Table: gcp_compute_disk_metric_write_ops_daily - Query GCP Compute Engine Disk Metrics using SQL

Google Cloud's Compute Engine Disk is a block storage system for Google Compute Engine virtual machines. The compute engine disks provide persistent disk storage for instances in any zone or region. These disks are integrated with Google's infrastructure, ensuring data durability and security.

## Table Usage Guide

The `gcp_compute_disk_metric_write_ops_daily` table provides insights into Compute Engine Disk Metrics within Google Cloud Platform (GCP). As a system administrator, you can explore disk-specific details through this table, including daily write operations count. Utilize it to uncover information about disk usage, such as those with high write operations, thereby allowing you to identify potential areas for efficiency improvements.

GCP Monitoring metrics provide data about the performance of your systems. The `gcp_compute_disk_metric_write_ops_daily` table provides metric statistics at 24 hour intervals for the last year.

## Examples

### Basic info
Explore which Google Cloud Platform compute disk has the highest and lowest average daily write operations. This can help to identify instances where disks may be under or over-utilized, allowing for better resource allocation and performance optimization.

```sql+postgres
select
  name,
  minimum,
  maximum,
  average,
  sample_count
from
  gcp_compute_disk_metric_write_ops_daily
order by
  name;
```

```sql+sqlite
select
  name,
  minimum,
  maximum,
  average,
  sample_count
from
  gcp_compute_disk_metric_write_ops_daily
order by
  name;
```

### Intervals averaging over 100 write ops
Analyze the settings to understand which disks are experiencing a high average of write operations. This can be useful in identifying potential bottlenecks or performance issues in the system.

```sql+postgres
select
  name,
  round(minimum::numeric,2) as min_write_ops,
  round(maximum::numeric,2) as max_write_ops,
  round(average::numeric,2) as avg_write_ops,
  sample_count
from
  gcp_compute_disk_metric_write_ops_daily
where
  average > 10
order by
  name;
```

```sql+sqlite
select
  name,
  round(minimum,2) as min_write_ops,
  round(maximum,2) as max_write_ops,
  round(average,2) as avg_write_ops,
  sample_count
from
  gcp_compute_disk_metric_write_ops_daily
where
  average > 10
order by
  name;
```

### Intervals averaging fewer than 1 write ops
Analyze the usage patterns of your Google Cloud Platform compute disk by identifying those instances where the average daily write operations fall below 1. This can help in optimizing resource allocation by pinpointing under-utilized disks.

```sql+postgres
select
  name,
  round(minimum::numeric,2) as min_write_ops,
  round(maximum::numeric,2) as max_write_ops,
  round(average::numeric,2) as avg_write_ops,
  sample_count
from
  gcp_compute_disk_metric_write_ops_daily
where
  average < 1
order by
  name;
```

```sql+sqlite
select
  name,
  round(minimum,2) as min_write_ops,
  round(maximum,2) as max_write_ops,
  round(average,2) as avg_write_ops,
  sample_count
from
  gcp_compute_disk_metric_write_ops_daily
where
  average < 1
order by
  name;
```