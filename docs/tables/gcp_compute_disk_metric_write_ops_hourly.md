# Table: gcp_compute_disk_metric_write_ops_hourly

Google cloud Monitoring Metrics provide data about the performance of your systems. The `gcp_compute_disk_metric_write_ops_hourly` table provides metric statistics at 1 hour intervals for the past year.

## Examples

### Basic info

```sql
select
  name,
  minimum,
  maximum,
  average,
  sample_count
from
  gcp_compute_disk_metric_write_ops_hourly
order by
  name;
```

### Intervals averaging over 100 write ops

```sql
select
  name,
  round(minimum::numeric,2) as min_write_ops,
  round(maximum::numeric,2) as max_write_ops,
  round(average::numeric,2) as avg_write_ops,
  sample_count
from
  gcp_compute_disk_metric_write_ops_hourly
where average > 10
order by
  name;
```