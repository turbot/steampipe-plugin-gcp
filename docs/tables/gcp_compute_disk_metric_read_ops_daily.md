---
title: "Steampipe Table: gcp_compute_disk_metric_read_ops_daily - Query GCP Compute Disk Metrics using SQL"
description: "Allows users to query GCP Compute Disk Metrics, specifically the daily read operations count, providing insights into disk usage patterns and potential anomalies."
folder: "Compute"
---

# Table: gcp_compute_disk_metric_read_ops_daily - Query GCP Compute Disk Metrics using SQL

Google Cloud Compute Disks are persistent, high-performance block storage for Google Cloud's virtual machines (VMs). These disks are designed to offer reliable and efficient storage for your VMs, with the added benefit of easy integration with Google Cloud's suite of data management tools. They are suitable for both boot and non-boot purposes, and come in a variety of types to suit different needs.

## Table Usage Guide

The `gcp_compute_disk_metric_read_ops_daily` table provides insights into Compute Disk Metrics within Google Cloud Platform. As a system administrator, explore disk-specific details through this table, such as daily read operations count, to monitor disk usage patterns and identify potential anomalies. Utilize it to uncover information about disks, such as those with high read operations, suggesting high disk usage and potential need for resource optimization.

GCP Monitoring metrics provide data about the performance of your systems. The `gcp_compute_disk_metric_read_ops_daily` table provides metric statistics at 24 hour intervals for the last year.

## Examples

### Basic info
Explore the daily read operations of your Google Cloud Platform compute disk. This query helps you understand the operational range and average count, allowing you to optimize disk usage and performance.

```sql+postgres
select
  name,
  minimum,
  maximum,
  average,
  sample_count
from
  gcp_compute_disk_metric_read_ops_daily
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
  gcp_compute_disk_metric_read_ops_daily
order by
  name;
```

### Intervals averaging over 100 read ops
Explore which disk operations in your GCP Compute environment are averaging over 100 read operations, allowing you to identify potential bottlenecks and optimize for better performance. This can be particularly useful in identifying high-usage disks that may need additional resources or configuration changes.

```sql+postgres
select
  name,
  round(minimum::numeric,2) as min_read_ops,
  round(maximum::numeric,2) as max_read_ops,
  round(average::numeric,2) as avg_read_ops,
  sample_count
from
  gcp_compute_disk_metric_read_ops_daily
where
  average > 10
order by
  name;
```

```sql+sqlite
select
  name,
  round(minimum,2) as min_read_ops,
  round(maximum,2) as max_read_ops,
  round(average,2) as avg_read_ops,
  sample_count
from
  gcp_compute_disk_metric_read_ops_daily
where
  average > 10
order by
  name;
```

### Intervals averaging fewer than 10 read ops
Analyze the settings to understand the performance of your GCP compute disks, specifically those with an average of fewer than 10 read operations. This can help in identifying underutilized resources and optimizing your storage configuration.

```sql+postgres
select
  name,
  round(minimum::numeric,2) as min_read_ops,
  round(maximum::numeric,2) as max_read_ops,
  round(average::numeric,2) as avg_read_ops,
  sample_count
from
  gcp_compute_disk_metric_read_ops_daily
where
  average < 1
order by
  name;
```

```sql+sqlite
select
  name,
  round(minimum,2) as min_read_ops,
  round(maximum,2) as max_read_ops,
  round(average,2) as avg_read_ops,
  sample_count
from
  gcp_compute_disk_metric_read_ops_daily
where
  average < 1
order by
  name;
```