# Table: gcp_disk_metric_read_ops_hourly

Google cloud Monitoring Metrics provide data about the performance of your systems. The `gcp_disk_metric_read_ops_hourly` table provides metric statistics at 24 hour intervals for the past year.

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
  gcp_disk_metric_read_ops_hourly
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
  gcp_disk_metric_read_ops_hourly
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
  gcp_disk_metric_read_ops_hourly
where average < 10
order by
  instance_id;
```
