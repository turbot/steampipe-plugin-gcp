---
title: "Steampipe Table: gcp_compute_instance_metric_cpu_utilization_daily - Query GCP Compute Engine Instances using SQL"
description: "Allows users to query daily CPU utilization metrics of GCP Compute Engine Instances, providing insights into their performance and resource usage patterns."
folder: "Compute"
---

# Table: gcp_compute_instance_metric_cpu_utilization_daily - Query GCP Compute Engine Instances using SQL

Google Cloud Compute Engine is a service that provides secure and customizable compute instances that can be used to build and host your applications. These instances are highly scalable and flexible, offering a variety of machine types to suit your needs. Compute Engine instances can be managed through the Google Cloud Console, the RESTful API, or the command-line interface.

## Table Usage Guide

The `gcp_compute_instance_metric_cpu_utilization_daily` table provides insights into the daily CPU utilization metrics of Google Cloud Compute Engine instances. As a system administrator or a DevOps engineer, you can explore instance-specific details through this table, including the CPU usage patterns, to manage and optimize resource allocation effectively. Use this table to monitor the performance of your instances, identify those with high CPU usage, and make informed decisions about scaling your resources.

GCP Monitoring Metrics provide data about the performance of your systems. The `gcp_compute_instance_metric_cpu_utilization_daily` table provides metric statistics at 24 hour intervals for the most recent 5 days.

## Examples

### Basic info
Analyze the daily CPU utilization metrics of your Google Cloud Compute instances to gain insights into usage patterns and performance. This can be particularly useful for capacity planning, identifying resource-intensive instances, and optimizing costs.

```sql+postgres
select
  name,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  gcp_compute_instance_metric_cpu_utilization_daily
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
  gcp_compute_instance_metric_cpu_utilization_daily
order by
  name,
  timestamp;
```

### CPU Over 80% average
Analyze the instances where CPU utilization exceeds 80% on average. This query can help in identifying potential performance issues and ensuring optimal resource allocation.

```sql+postgres
select
  name,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  gcp_compute_instance_metric_cpu_utilization_daily
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
  gcp_compute_instance_metric_cpu_utilization_daily
where
  average > 0.80
order by
  name,
  timestamp;
```