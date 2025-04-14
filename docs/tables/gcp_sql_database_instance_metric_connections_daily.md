---
title: "Steampipe Table: gcp_sql_database_instance_metric_connections_daily - Query GCP SQL Database Instance Metrics using SQL"
description: "Allows users to query GCP SQL Database Instance Metrics, specifically daily connection metrics for Google Cloud SQL instances, providing insights into database connection patterns and potential anomalies."
folder: "SQL"
---

# Table: gcp_sql_database_instance_metric_connections_daily - Query GCP SQL Database Instance Metrics using SQL

Google Cloud SQL is a fully-managed database service that makes it easy to set up, maintain, manage, and administer your relational databases on Google Cloud Platform. You can use Cloud SQL with MySQL, PostgreSQL, or SQL Server. The service provides detailed metrics on database instances, including connection metrics, to help monitor the health and performance of your databases.

## Table Usage Guide

The `gcp_sql_database_instance_metric_connections_daily` table provides insights into daily connection metrics for Google Cloud SQL instances. As a database administrator, you can explore instance-specific details through this table, including the number of established connections, aborted connections, and failed connection attempts. Use it to monitor database connection patterns, identify potential anomalies, and optimize your database performance.

GCP Monitoring metrics provide data about the performance of your systems. The `gcp_sql_database_instance_metric_connections_daily` table provides metric statistics at 24 hour intervals for the past year.

## Examples

### Basic info
Analyze the daily connection metrics of Google Cloud SQL database instances to gain insights into their usage patterns. This could be useful for capacity planning and performance optimization.

```sql+postgres
select
  instance_id,
  minimum,
  maximum,
  average,
  sample_count
from
  gcp_sql_database_instance_metric_connections_daily
order by
  instance_id;
```

```sql+sqlite
select
  instance_id,
  minimum,
  maximum,
  average,
  sample_count
from
  gcp_sql_database_instance_metric_connections_daily
order by
  instance_id;
```

### Intervals averaging over 100 connections
Explore which Google Cloud SQL database instances have an average daily connection count exceeding 100. This is useful for identifying potentially over-utilized instances that may require capacity scaling or performance optimization.

```sql+postgres
select
  instance_id,
  round(minimum::numeric,2) as min_connection,
  round(maximum::numeric,2) as max_connection,
  round(average::numeric,2) as avg_connection,
  sample_count
from
  gcp_sql_database_instance_metric_connections_daily
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
  gcp_sql_database_instance_metric_connections_daily
where
  average > 100
order by
  instance_id;
```

### Intervals averaging fewer than 10 connections
Analyze the settings to understand the performance of your Google Cloud SQL Database instances, specifically those that average fewer than 10 connections daily. This can be beneficial in identifying underutilized instances and optimizing resource allocation.

```sql+postgres
select
  instance_id,
  round(minimum::numeric,2) as min_connection,
  round(maximum::numeric,2) as max_connection,
  round(average::numeric,2) as avg_connection,
  sample_count
from
  gcp_sql_database_instance_metric_connections_daily
where
  average < 10
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
  gcp_sql_database_instance_metric_connections_daily
where
  average < 10
order by
  instance_id;
```