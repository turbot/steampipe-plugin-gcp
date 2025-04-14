---
title: "Steampipe Table: gcp_compute_instance_metric_cpu_utilization - Query Google Cloud Compute Engine Instance Metrics using SQL"
description: "Allows users to query Compute Engine Instance CPU Utilization Metrics in Google Cloud Platform (GCP), specifically the CPU utilization of each instance, providing insights into resource usage and potential performance issues."
folder: "Compute"
---

# Table: gcp_compute_instance_metric_cpu_utilization - Query Google Cloud Compute Engine Instance Metrics using SQL

Google Cloud Compute Engine is a service within Google Cloud Platform that offers scalable and flexible virtual machine computing capabilities. It allows you to run large-scale workloads on virtual machines hosted on Google's infrastructure. Compute Engine instances can be tailored to specific workloads for optimal performance and cost efficiency.

## Table Usage Guide

The `gcp_compute_instance_metric_cpu_utilization` table provides insights into the CPU utilization of Compute Engine Instances within Google Cloud Platform (GCP). As a system administrator or DevOps engineer, explore instance-specific details through this table, including CPU usage patterns, potential performance bottlenecks, and resource optimization opportunities. Utilize it to uncover information about instances, such as those with high CPU utilization, to make informed decisions about resource allocation and workload management.

Google Monitoring Metrics provide data about the performance of your systems. The `gcp_compute_instance_metric_cpu_utilization` table provides metric statistics at 5 minute intervals for the most recent 5 days.

## Examples

### Basic info
Explore the utilization of your Google Cloud Platform's compute instances over time. This query helps you understand the CPU usage patterns, including peaks and average usage, enabling you to optimize resource allocation and manage costs effectively.

```sql+postgres
select
  name,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  gcp_compute_instance_metric_cpu_utilization
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
  gcp_compute_instance_metric_cpu_utilization
order by
  name,
  timestamp;
```

### CPU Over 80% average
Determine the areas in which the average CPU utilization exceeds 80%. This can be used to identify potential performance issues and ensure efficient resource allocation.

```sql+postgres
select
  name,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  gcp_compute_instance_metric_cpu_utilization
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
  gcp_compute_instance_metric_cpu_utilization
where
  average > 0.80
order by
  name,
  timestamp;
```