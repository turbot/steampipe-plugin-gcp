# Table: gcp_sql_database_instance_connections_daily

GCP Monitoring metrics provide data about the performance of your systems. The `gcp_sql_database_instance_connections_daily` table provides metric statistics at 24 hour intervals for the past year.

## Examples

### Basic info

```sql
select
  instance_id,
  minimum,
  maximum,
  average,
  sample_count
from
  gcp_sql_database_connections_daily
order by
  instance_id;
```

### Intervals averaging over 100 connections

```sql
select
  instance_id,
  round(minimum::numeric,2) as min_connection,
  round(maximum::numeric,2) as max_connection,
  round(average::numeric,2) as avg_connection,
  sample_count
from
  gcp_sql_database_connections_daily
where average > 100
order by
  instance_id;
```

### Intervals averaging fewer than 10 connections

```sql
select
  instance_id,
  round(minimum::numeric,2) as min_connection,
  round(maximum::numeric,2) as max_connection,
  round(average::numeric,2) as avg_connection,
  sample_count
from
  gcp_sql_database_connections_daily
where average < 10
order by
  instance_id;
```
