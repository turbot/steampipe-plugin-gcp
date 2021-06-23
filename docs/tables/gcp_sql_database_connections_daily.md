# Table: gcp_sql_database_connection_daily

Google cloud Monitoring Metrics provide data about the performance of your systems. The `gcp_sql_database_connection_daily` table provides metric statistics at 24 hour intervals for the last year.

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
  gcp_sql_database_connection_daily
order by
  instance_id;
```

### Connection Over 100 average

```sql
select
  instance_id,
  round(minimum::numeric,2) as min_connection,
  round(maximum::numeric,2) as max_connection,
  round(average::numeric,2) as avg_connection,
  sample_count
from
  gcp_sql_database_connection_daily
where average > 100
order by
  instance_id;
```

### Connection daily average < 10

```sql
select
  instance_id,
  round(minimum::numeric,2) as min_connection,
  round(maximum::numeric,2) as max_connection,
  round(average::numeric,2) as avg_connection,
  sample_count
from
  gcp_sql_database_connection_daily
where average < 10
order by
  instance_id;
```
