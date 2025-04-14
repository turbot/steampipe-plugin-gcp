---
title: "Steampipe Table: gcp_sql_database_instance_metric_cpu_utilization_daily - Query GCP SQL Database Instance Metrics using SQL"
description: "Allows users to query SQL Database Instance Metrics in GCP, specifically the daily CPU utilization, providing insights into resource usage patterns and potential performance bottlenecks."
folder: "SQL"
---

# Table: gcp_sql_database_instance_metric_cpu_utilization_daily - Query GCP SQL Database Instance Metrics using SQL

A Google Cloud SQL Database Instance is a fully-managed relational database service that makes it easy to set up, manage, and administer relational databases on Google Cloud. It provides a cost-effective and scalable way to operate MySQL, PostgreSQL, and SQL Server instances in the cloud. This service is designed to handle the demanding, heavy-duty workloads of high-performance applications.

## Table Usage Guide

The `gcp_sql_database_instance_metric_cpu_utilization_daily` table provides insights into the CPU utilization of SQL Database Instances within Google Cloud Platform (GCP). As a Database Administrator or DevOps Engineer, explore instance-specific details through this table, including daily CPU usage patterns. Utilize it to uncover information about instances, such as those with high resource usage, helping you to optimize performance and manage costs effectively.

GCP Monitoring metrics provide data about the performance of your systems. The `gcp_sql_database_instance_metric_cpu_utilization_daily` table provides metric statistics at 24 hour intervals for the past year.

## Examples

### Basic info
Explore the daily CPU utilization metrics of your Google Cloud SQL database instances to understand their performance. This can help you identify any instances that may be under or over-utilized, allowing you to optimize resource allocation and cost.

```sql+postgres
select
  instance_id,
  minimum,
  maximum,
  average,
  sample_count,
  timestamp
from
  gcp_sql_database_instance_metric_cpu_utilization_daily
order by
  instance_id;
```

```sql+sqlite
select
  instance_id,
  minimum,
  maximum,
  average,
  sample_count,
  timestamp
from
  gcp_sql_database_instance_metric_cpu_utilization_daily
order by
  instance_id;
```

### Intervals averaging over 100%
Explore which instances have an average connection that exceeds 100%, allowing you to identify potential areas of overutilization for further investigation.

```sql+postgres
select
  instance_id,
  round(minimum::numeric,2) as min_connection,
  round(maximum::numeric,2) as max_connection,
  round(average::numeric,2) as avg_connection,
  sample_count
from
  gcp_sql_database_instance_metric_cpu_utilization_daily
where
  average > 100
order by
  instance_id;
```

```sql+sqlite
select
  instance_id,
  round(minimum,2) as min_connection,
  round(maximum,2) as max_connection,
  round(average,2) as avg_connection,
  sample_count
from
  gcp_sql_database_instance_metric_cpu_utilization_daily
where
  average > 100
order by
  instance_id;
```

### Intervals averaging < 1%
Determine the areas in which Google Cloud SQL database instances have an average CPU utilization of less than 1% per day. This can help optimize resource allocation by identifying under-utilized instances.

```sql+postgres
select
  instance_id,
  round(minimum::numeric,2) as min_connection,
  round(maximum::numeric,2) as max_connection,
  round(average::numeric,2) as avg_connection,
  sample_count
from
  gcp_sql_database_instance_metric_cpu_utilization_daily
where
  average < 1
order by
  instance_id;
```

```sql+sqlite
select
  instance_id,
  round(minimum,2) as min_connection,
  round(maximum,2) as max_connection,
  round(average,2) as avg_connection,
  sample_count
from
  gcp_sql_database_instance_metric_cpu_utilization_daily
where
  average < 1
order by
  instance_id;
```