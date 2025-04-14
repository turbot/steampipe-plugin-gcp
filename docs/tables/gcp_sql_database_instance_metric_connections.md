---
title: "Steampipe Table: gcp_sql_database_instance_metric_connections - Query Google Cloud SQL Database Instance Connections using SQL"
description: "Allows users to query Google Cloud SQL Database Instance Connections, specifically the current client connections to a Cloud SQL database instance, providing insights into connection patterns and potential issues."
folder: "SQL"
---

# Table: gcp_sql_database_instance_metric_connections - Query Google Cloud SQL Database Instance Connections using SQL

Google Cloud SQL is a fully-managed database service that makes it easy to set up, maintain, manage, and administer relational databases on Google Cloud Platform. The service provides a highly available and scalable cloud database environment, supporting both MySQL and PostgreSQL databases. It allows users to focus on application development, freeing them from the typical database administration tasks.

## Table Usage Guide

The `gcp_sql_database_instance_metric_connections` table provides insights into the current client connections to a Google Cloud SQL Database Instance. As a database administrator, you can explore connection-specific details through this table, including the connection count and related metrics. Utilize it to monitor connection patterns, diagnose potential connection issues, and ensure optimal database performance.

GCP Monitoring metrics provide data about the performance of your systems. The `gcp_sql_database_instance_metric_connections` table provides metric statistics at 5 minute intervals for the most recent 5 days.

## Examples

### Basic info
Explore the performance of your Google Cloud SQL databases by analyzing the minimum, maximum, and average number of connections over time. This can help in identifying potential bottlenecks and planning for capacity upgrades.

```sql+postgres
select
  instance_id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  gcp_sql_database_instance_metric_connections
order by
  instance_id;
```

```sql+sqlite
select
  instance_id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  gcp_sql_database_instance_metric_connections
order by
  instance_id;
```

### Intervals averaging over 100 connections
Analyze the settings to understand instances where the average number of connections to a Google Cloud SQL database exceeds 100. This can be useful for identifying potential bottlenecks or periods of high demand in your database infrastructure.

```sql+postgres
select
  instance_id,
  timestamp,
  round(minimum::numeric,2) as min_connection,
  round(maximum::numeric,2) as max_connection,
  round(average::numeric,2) as avg_connection,
  sample_count
from
  gcp_sql_database_instance_metric_connections
where
  average > 100
order by
  instance_id;
```

```sql+sqlite
select
  instance_id,
  timestamp,
  round(minimum,2) as min_connection,
  round(maximum,2) as max_connection,
  round(average,2) as avg_connection,
  sample_count
from
  gcp_sql_database_instance_metric_connections
where
  average > 100
order by
  instance_id;
```

### Intervals averaging fewer than 10 connections
Explore which instances in your database are operating with an average of fewer than 10 connections. This can help in optimizing resource allocation by identifying under-utilized instances.

```sql+postgres
select
  instance_id,
  timestamp,
  round(minimum::numeric,2) as min_connection,
  round(maximum::numeric,2) as max_connection,
  round(average::numeric,2) as avg_connection,
  sample_count
from
  gcp_sql_database_instance_metric_connections
where
  average < 10
order by
  instance_id;
```

```sql+sqlite
select
  instance_id,
  timestamp,
  round(minimum,2) as min_connection,
  round(maximum,2) as max_connection,
  round(average,2) as avg_connection,
  sample_count
from
  gcp_sql_database_instance_metric_connections
where average < 10
order by
  instance_id;
```