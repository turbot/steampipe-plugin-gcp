---
title: "Steampipe Table: gcp_compute_disk_metric_write_ops_hourly - Query Google Cloud Compute Engine Disks using SQL"
description: "Allows users to query Google Cloud Compute Engine Disks, specifically the hourly write operations metrics, providing insights into disk usage patterns and potential performance issues."
folder: "Compute"
---

# Table: gcp_compute_disk_metric_write_ops_hourly - Query Google Cloud Compute Engine Disks using SQL

Google Cloud Compute Engine Disks are persistent, high-performance block storage for Google Cloud's Virtual Machines (VMs). They offer a range of options to accommodate varying storage capacity, performance, and cost needs. These disks can be attached to instances within the same region.

## Table Usage Guide

The `gcp_compute_disk_metric_write_ops_hourly` table provides insights into the hourly write operations of Google Cloud Compute Engine Disks. As a system administrator or DevOps engineer, explore disk-specific details through this table, including the number of write operations, associated metadata, and timestamps. Utilize it to monitor disk usage patterns, optimize disk performance, and troubleshoot potential issues.

GCP Monitoring metrics provide data about the performance of your systems. The `gcp_compute_disk_metric_write_ops_hourly` table provides metric statistics at 1 hour intervals for the most recent 60 days.

## Examples

### Basic info
Explore the performance of your Google Cloud Platform (GCP) compute disks by analyzing their hourly write operations. This can help you understand disk usage patterns, identify potential bottlenecks, and plan for capacity accordingly.

```sql+postgres
select
  name,
  minimum,
  maximum,
  average,
  sample_count
from
  gcp_compute_disk_metric_write_ops_hourly
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
  gcp_compute_disk_metric_write_ops_hourly
order by
  name;
```

### Intervals averaging over 100 write ops
Explore which disk operations have an average of over 10 write operations, allowing you to identify potential high-usage instances and optimize for better performance.

```sql+postgres
select
  name,
  round(minimum::numeric,2) as min_write_ops,
  round(maximum::numeric,2) as max_write_ops,
  round(average::numeric,2) as avg_write_ops,
  sample_count
from
  gcp_compute_disk_metric_write_ops_hourly
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
  gcp_compute_disk_metric_write_ops_hourly
where
  average > 10
order by
  name;
```