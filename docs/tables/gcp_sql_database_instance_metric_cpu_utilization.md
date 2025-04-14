---
title: "Steampipe Table: gcp_sql_database_instance_metric_cpu_utilization - Query Google Cloud Platform SQL Database Instances using SQL"
description: "Allows users to query CPU Utilization Metrics for Google Cloud Platform SQL Database Instances. This provides insights into the CPU usage patterns and potential performance bottlenecks."
folder: "SQL"
---

# Table: gcp_sql_database_instance_metric_cpu_utilization - Query Google Cloud Platform SQL Database Instances using SQL

Google Cloud SQL is a fully-managed relational database service that makes it easy to set up, manage, and administer relational databases on Google Cloud Platform. It provides a way to run standard SQL queries on your data with ease. This service supports MySQL, PostgreSQL, and SQL Server.

## Table Usage Guide

The `gcp_sql_database_instance_metric_cpu_utilization` table provides insights into the CPU utilization of Google Cloud SQL Database Instances. As a database administrator or a cloud engineer, you can use this table to monitor CPU usage patterns and identify potential performance issues. This can be particularly useful in capacity planning and performance optimization tasks.

GCP Monitoring metrics provide data about the performance of your systems. The `gcp_sql_database_instance_metric_cpu_utilization` table provides metric statistics at 5 minute intervals for the most recent 5 days.

## Examples

### Basic info
Determine the performance of your Google Cloud SQL database instances by analyzing CPU utilization over time. This query helps in identifying instances with irregular CPU usage patterns, which could indicate potential issues or areas for optimization.

```sql+postgres
select
  instance_id,
  minimum,
  maximum,
  average,
  sample_count,
  timestamp
from
  gcp_sql_database_instance_metric_cpu_utilization
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
  gcp_sql_database_instance_metric_cpu_utilization
order by
  instance_id;
```

### Intervals averaging over 80%
Explore which instances in your Google Cloud SQL database have an average CPU utilization over 80%. This can help determine areas of potential overload or inefficiency, allowing you to better manage your resources.

```sql+postgres
select
  instance_id,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  gcp_sql_database_instance_metric_cpu_utilization
where
  average > 80
order by
  instance_id;
```

```sql+sqlite
select
  instance_id,
  round(minimum,2) as min_cpu,
  round(maximum,2) as max_cpu,
  round(average,2) as avg_cpu,
  sample_count
from
  gcp_sql_database_instance_metric_cpu_utilization
where
  average > 80
order by
  instance_id;
```

### Intervals averaging < 1%
Explore instances of Google Cloud SQL databases where the average CPU utilization is less than 1%. This can help identify under-utilized resources, potentially leading to cost savings.

```sql+postgres
select
  instance_id,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  gcp_sql_database_instance_metric_cpu_utilization
where average < 1
order by
  instance_id;
```

```sql+sqlite
select
  instance_id,
  round(minimum,2) as min_cpu,
  round(maximum,2) as max_cpu,
  round(average,2) as avg_cpu,
  sample_count
from
  gcp_sql_database_instance_metric_cpu_utilization
where average < 1
order by
  instance_id;
```