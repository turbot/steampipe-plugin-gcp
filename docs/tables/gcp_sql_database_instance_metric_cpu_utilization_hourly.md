---
title: "Steampipe Table: gcp_sql_database_instance_metric_cpu_utilization_hourly - Query Google Cloud SQL Database Instances using SQL"
description: "Allows users to query Google Cloud SQL Database Instances, specifically hourly CPU utilization metrics, providing insights into resource usage and potential performance bottlenecks."
folder: "SQL"
---

# Table: gcp_sql_database_instance_metric_cpu_utilization_hourly - Query Google Cloud SQL Database Instances using SQL

Google Cloud SQL is a fully-managed database service that makes it easy to set up, maintain, manage, and administer your relational databases on Google Cloud Platform. The service offers seamless scalability, secure connections, and the flexibility to support various SQL workloads. Google Cloud SQL Database Instances are the building blocks of this service where the databases are hosted.

## Table Usage Guide

The `gcp_sql_database_instance_metric_cpu_utilization_hourly` table provides insights into the CPU utilization of Google Cloud SQL Database Instances on an hourly basis. As a database administrator or DevOps engineer, you can use this table to monitor and analyze the CPU usage of your instances, helping you to identify resource-intensive operations, potential performance issues, and opportunities for optimization. This information is crucial for maintaining the efficiency and reliability of your database operations.

GCP Monitoring metrics provide data about the performance of your systems. The `gcp_sql_database_instance_metric_cpu_utilization_hourly` table provides metric statistics at 1 hour intervals for the most recent 60 days.

## Examples

### Basic info
Explore the CPU utilization of your Google Cloud SQL database instances over the past hour. This can help you understand their performance and identify any instances that may be under or over-utilized.

```sql+postgres
select
  instance_id,
  minimum,
  maximum,
  average,
  sample_count,
  timestamp
from
  gcp_sql_database_instance_metric_cpu_utilization_hourly
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
  gcp_sql_database_instance_metric_cpu_utilization_hourly
order by
  instance_id;
```

### Intervals averaging over 80%
Identify instances where the CPU utilization of your Google Cloud SQL database averages over 80% to help manage resources and optimize performance.

```sql+postgres
select
  instance_id,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  gcp_sql_database_instance_metric_cpu_utilization_hourly
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
  gcp_sql_database_instance_metric_cpu_utilization_hourly
where
  average > 80
order by
  instance_id;
```

### Intervals averaging < 1%
Analyze the CPU utilization of your Google Cloud SQL database instances to identify those with an average CPU utilization of less than 1%. This can be particularly useful in optimizing resource allocation and reducing costs by pinpointing underutilized instances.

```sql+postgres
select
  instance_id,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  gcp_sql_database_instance_metric_cpu_utilization_hourly
where
  average < 1
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
  gcp_sql_database_instance_metric_cpu_utilization_hourly
where
  average < 1
order by
  instance_id;
```