---
title: "Steampipe Table: gcp_sql_database_instance_metric_connections_hourly - Query GCP SQL Database Instances using SQL"
description: "Allows users to query GCP SQL Database Instances, specifically the hourly metric connections, providing insights into database connection patterns and potential issues."
folder: "SQL"
---

# Table: gcp_sql_database_instance_metric_connections_hourly - Query GCP SQL Database Instances using SQL

A Google Cloud SQL Database Instance is a fully-managed relational database service in Google Cloud. It offers easy-to-use, scalable database instances that support applications of any size. It provides high performance, scalability, and convenience with features such as automated backups, replication, and failover.

## Table Usage Guide

The `gcp_sql_database_instance_metric_connections_hourly` table provides insights into the hourly metric connections of Google Cloud SQL Database Instances. As a database administrator or a DevOps engineer, you can use this table to monitor and analyze connection patterns, which can help in identifying potential issues and optimizing database performance. This table is particularly useful for tracking connection trends, identifying peak usage times, and planning capacity.

GCP Monitoring metrics provide data about the performance of your systems. The `gcp_sql_database_instance_metric_connections_hourly` table provides metric statistics at 1 hour intervals for the most recent 60 days.

## Examples

### Basic info
Explore the varying connection metrics of Google Cloud SQL databases to understand their performance over time. This can assist in identifying potential issues or areas for optimization based on the minimum, maximum, and average connections.

```sql+postgres
select
  instance_id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  gcp_sql_database_instance_metric_connections_hourly
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
  gcp_sql_database_instance_metric_connections_hourly
order by
  instance_id;
```

### Intervals averaging over 100 connections
Explore instances where database connections average over 100 per hour, offering insights into potential high-traffic periods or potential performance issues. This could be crucial for capacity planning, load balancing and optimizing database performance.

```sql+postgres
select
  instance_id,
  timestamp,
  round(minimum::numeric,2) as min_connection,
  round(maximum::numeric,2) as max_connection,
  round(average::numeric,2) as avg_connection,
  sample_count
from
  gcp_sql_database_instance_metric_connections_hourly
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
  gcp_sql_database_instance_metric_connections_hourly
where
  average > 100
order by
  instance_id;
```

### Intervals averaging fewer than 10 connections
Explore which instances are maintaining an average connection count of less than 10. This is useful for identifying underutilized resources and optimizing your database for cost efficiency.

```sql+postgres
select
  instance_id,
  timestamp,
  round(minimum::numeric,2) as min_connection,
  round(maximum::numeric,2) as max_connection,
  round(average::numeric,2) as avg_connection,
  sample_count
from
  gcp_sql_database_instance_metric_connections_hourly
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
  gcp_sql_database_instance_metric_connections_hourly
where
  average < 10
order by
  instance_id;
```