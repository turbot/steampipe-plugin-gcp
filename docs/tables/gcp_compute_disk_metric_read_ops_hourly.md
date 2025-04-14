---
title: "Steampipe Table: gcp_compute_disk_metric_read_ops_hourly - Query Google Cloud Platform Compute Engine Disks using SQL"
description: "Allows users to query Compute Engine Disks in Google Cloud Platform, specifically the hourly read operations metric, providing insights into disk usage patterns and potential performance issues."
folder: "Compute"
---

# Table: gcp_compute_disk_metric_read_ops_hourly - Query Google Cloud Platform Compute Engine Disks using SQL

Google Cloud Compute Engine Disks are persistent, high-performance block storage for Google Cloud Platform virtual machines. They are designed to offer reliable and efficient storage for your workloads, with features such as automatic encryption, snapshot capabilities, and seamless integration with Google Cloud Platform services. Compute Engine Disks provide the flexibility to balance cost and performance for your storage needs.

## Table Usage Guide

The `gcp_compute_disk_metric_read_ops_hourly` table provides insights into Compute Engine Disks within Google Cloud Platform. As a DevOps engineer, explore disk-specific details through this table, including hourly read operations metrics. Utilize it to uncover information about disk usage patterns, such as high read operations, which could indicate potential performance issues.

GCP Monitoring metrics provide data about the performance of your systems. The `gcp_compute_disk_metric_read_ops_hourly` table provides metric statistics at 1 hour intervals for the most recent 60 days.

## Examples

### Basic info
Explore the range and average of read operations on your Google Cloud Platform compute disks over an hourly period. This can help you understand disk usage patterns and identify potential areas for performance optimization.

```sql+postgres
select
  name,
  minimum,
  maximum,
  average,
  sample_count
from
  gcp_compute_disk_metric_read_ops_hourly
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
  gcp_compute_disk_metric_read_ops_hourly
order by
  name;
```

### Intervals averaging over 100 read operations
Explore which disk operations have an average of over 100 read operations. This can be helpful in identifying areas where resource usage may be high, potentially indicating a need for optimization or increased capacity.

```sql+postgres
select
  name,
  round(minimum::numeric,2) as min_read_ops,
  round(maximum::numeric,2) as max_read_ops,
  round(average::numeric,2) as avg_read_ops,
  sample_count
from
  gcp_compute_disk_metric_read_ops_hourly
where
  average > 100
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
  gcp_compute_disk_metric_read_ops_hourly
where
  average > 100
order by
  name;
```

### Intervals averaging fewer than 10 read operations
Determine the areas in which disk operations in Google Cloud Compute are underutilized by identifying intervals where read operations average less than ten. This can help in optimizing resource allocation and managing costs effectively.

```sql+postgres
select
  name,
  round(minimum::numeric,2) as min_read_ops,
  round(maximum::numeric,2) as max_read_ops,
  round(average::numeric,2) as avg_read_ops,
  sample_count
from
  gcp_compute_disk_metric_read_ops_hourly
where
  average < 10
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
  gcp_compute_disk_metric_read_ops_hourly
where
  average < 10
order by
  name;
```