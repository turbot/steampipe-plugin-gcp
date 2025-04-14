---
title: "Steampipe Table: gcp_compute_disk_metric_read_ops - Query Google Cloud Compute Engine Disk Read Operations using SQL"
description: "Allows users to query Disk Read Operations in Google Cloud's Compute Engine, specifically the number of read operations completed, providing insights into disk usage and potential performance issues."
folder: "Compute"
---

# Table: gcp_compute_disk_metric_read_ops - Query Google Cloud Compute Engine Disk Read Operations using SQL

Google Cloud's Compute Engine is an Infrastructure as a Service that allows you to run your large-scale computing workloads on virtual machines hosted on Google's infrastructure. Disk Read Operations represents the count of read operations completed by the Compute Engine. This count can be used to analyze disk usage and performance.

## Table Usage Guide

The `gcp_compute_disk_metric_read_ops` table provides insights into Disk Read Operations within Google Cloud's Compute Engine. As a System Administrator or a DevOps engineer, explore disk-specific details through this table, including the number of read operations. Utilize it to uncover information about disk usage, such as frequent read operations, which can help in identifying potential performance issues and optimizing disk usage.

GCP Monitoring metrics provide data about the performance of your systems. The `gcp_compute_disk_metric_read_ops` table provides metric statistics at 5 minute intervals for the most recent 5 days.

## Examples

### Basic info
Explore which Google Cloud Platform Compute Disk has the highest and lowest read operations by assessing the minimum, maximum, and average values. This can help optimize disk usage and improve system performance.

```sql+postgres
select
  name,
  minimum,
  maximum,
  average,
  sample_count
from
  gcp_compute_disk_metric_read_ops
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
  gcp_compute_disk_metric_read_ops
order by
  name;
```

### Intervals averaging over 100 read ops
Explore which disk operations have an average read operation count over 10. This can help in identifying potential areas of high disk usage and performance bottlenecks.

```sql+postgres
select
  name,
  round(minimum::numeric,2) as min_read_ops,
  round(maximum::numeric,2) as max_read_ops,
  round(average::numeric,2) as avg_read_ops,
  sample_count
from
  gcp_compute_disk_metric_read_ops
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
  gcp_compute_disk_metric_read_ops
where
  average > 10
order by
  name;
```

### Intervals averaging fewer than 10 read ops
Analyze disk performance by identifying those with an average of less than 10 read operations. This can assist in pinpointing underutilized resources and optimizing system performance.

```sql+postgres
select
  name,
  round(minimum::numeric,2) as min_read_ops,
  round(maximum::numeric,2) as max_read_ops,
  round(average::numeric,2) as avg_read_ops,
  sample_count
from
  gcp_compute_disk_metric_read_ops
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
  gcp_compute_disk_metric_read_ops
where
  average < 10
order by
  name;
```