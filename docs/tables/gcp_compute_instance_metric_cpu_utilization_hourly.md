---
title: "Steampipe Table: gcp_compute_instance_metric_cpu_utilization_hourly - Query GCP Compute Engine Instance Metrics using SQL"
description: "Allows users to query Compute Engine Instance Metrics in GCP, specifically the hourly CPU utilization, providing insights into compute resource usage and potential performance bottlenecks."
folder: "Compute"
---

# Table: gcp_compute_instance_metric_cpu_utilization_hourly - Query GCP Compute Engine Instance Metrics using SQL

The Google Compute Engine is a service within Google Cloud Platform that provides highly customizable virtual machines with best-in-class features. It offers predefined virtual machines with specific amounts of CPU, memory, and storage to accommodate the needs of different workloads. Compute Engine also allows users to create custom machine types optimized for specific needs.

## Table Usage Guide

The `gcp_compute_instance_metric_cpu_utilization_hourly` table provides insights into the CPU utilization of Compute Engine instances on an hourly basis. As a system administrator or DevOps engineer, you can use this table to monitor and analyze the performance of your GCP virtual machines, helping you to optimize resource allocation and troubleshoot performance issues. It offers detailed information about CPU usage, allowing you to identify patterns and potential bottlenecks in your compute resources.

GCP Monitoring Metrics provide data about the performance of your systems. The `gcp_compute_instance_metric_cpu_utilization_hourly` table provides metric statistics at 60 minute intervals for the most recent 5 days.

## Examples

### Basic info
Explore the CPU utilization metrics of your Google Cloud Compute instances on an hourly basis. This can help you identify instances where resources might be under or over-utilized, enabling efficient resource allocation and cost optimization.

```sql+postgres
select
  name,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  gcp_compute_instance_metric_cpu_utilization_hourly
order by
  name,
  timestamp;
```

```sql+sqlite
select
  name,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  gcp_compute_instance_metric_cpu_utilization_hourly
order by
  name,
  timestamp;
```

### CPU Over 80% average
Identify instances where the average CPU utilization exceeds 80% on your Google Cloud Platform compute instances. This can assist in pinpointing potential performance issues or areas where resource allocation may need to be adjusted.

```sql+postgres
select
  name,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  gcp_compute_instance_metric_cpu_utilization_hourly
where
  average > 0.80
order by
  name,
  timestamp;
```

```sql+sqlite
select
  name,
  timestamp,
  round(minimum,2) as min_cpu,
  round(maximum,2) as max_cpu,
  round(average,2) as avg_cpu,
  sample_count
from
  gcp_compute_instance_metric_cpu_utilization_hourly
where
  average > 0.80
order by
  name,
  timestamp;
```