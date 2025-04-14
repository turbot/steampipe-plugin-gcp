---
title: "Steampipe Table: gcp_compute_disk_metric_write_ops - Query GCP Compute Disks using SQL"
description: "Allows users to query GCP Compute Disks, specifically focusing on write operations metrics, providing valuable insights into disk usage and performance."
folder: "Compute"
---

# Table: gcp_compute_disk_metric_write_ops - Query GCP Compute Disks using SQL

Google Cloud Platform's (GCP) Compute Disks are persistent, high-performance block storage for Google Compute Engine virtual machines. The disks are automatically encrypted, replicated across various zones, and can be easily increased in size. They provide the foundation for applications, databases, and file systems.

## Table Usage Guide

The `gcp_compute_disk_metric_write_ops` table provides insights into the write operations metrics of GCP Compute Disks. As a DevOps engineer, you can explore disk-specific details through this table, including the number of write operations and the rate of these operations. Utilize it to monitor disk performance, identify potential issues, and ensure optimal resource utilization.

GCP Monitoring metrics provide data about the performance of your systems. The `gcp_compute_disk_metric_write_ops` table provides metric statistics at 5 minute intervals for the most recent 5 days.

## Examples

### Basic info
Explore the performance of your Google Cloud Compute disks by analyzing metrics such as minimum, maximum, and average write operations. This can help in identifying any disks that may be underperforming or overutilized, enabling you to optimize your resources.

```sql+postgres
select
  name,
  minimum,
  maximum,
  average,
  sample_count
from
  gcp_compute_disk_metric_write_ops
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
  gcp_compute_disk_metric_write_ops
order by
  name;
```

### Intervals averaging over 100 write ops
Analyze the settings to understand which intervals have an average of over 100 write operations. This can help you pinpoint specific locations where high write operations occur, aiding in resource optimization and system performance enhancement.

```sql+postgres
select
  name,
  round(minimum::numeric,2) as min_write_ops,
  round(maximum::numeric,2) as max_write_ops,
  round(average::numeric,2) as avg_write_ops,
  sample_count
from
  gcp_compute_disk_metric_write_ops
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
  gcp_compute_disk_metric_write_ops
where
  average > 10
order by
  name;
```