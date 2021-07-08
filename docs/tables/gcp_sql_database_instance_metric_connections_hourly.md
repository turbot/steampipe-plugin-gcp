# Table: gcp_sql_database_instance_metric_connections_hourly

GCP Monitoring metrics provide data about the performance of your systems. The `gcp_sql_database_instance_metric_connections_hourly` table provides metric statistics at 1 hour intervals for the most recent 60 days.

## Examples

### Basic info

```sql
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

```sql
select
  instance_id,
  timestamp,
  round(minimum::numeric,2) as min_connection,
  round(maximum::numeric,2) as max_connection,
  round(average::numeric,2) as avg_connection,
  sample_count
from
  gcp_sql_database_instance_metric_connections_hourly
where average > 100
order by
  instance_id;
```

### Intervals averaging fewer than 10 connections

```sql
select
  instance_id,
  timestamp,
  round(minimum::numeric,2) as min_connection,
  round(maximum::numeric,2) as max_connection,
  round(average::numeric,2) as avg_connection,
  sample_count
from
  gcp_sql_database_instance_metric_connections_hourly
where average < 10
order by
  instance_id;
```
