# Table: gcp_compute_disk_metric_write_ops_daily

Google cloud Monitoring Metrics provide data about the performance of your systems. The `gcp_compute_disk_metric_write_ops_daily` table provides metric statistics at 5 min intervals for the past year.

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
  gcp_compute_disk_metric_write_ops_daily
order by
  instance_id;
```

### Intervals averaging over 100 write ops

```sql
select
  instance_id,
  round(minimum::numeric,2) as min_write_ops,
  round(maximum::numeric,2) as max_write_ops,
  round(average::numeric,2) as avg_write_ops,
  sample_count
from
  gcp_compute_disk_metric_write_ops_daily
where average > 10
order by
  instance_id;
```

### Intervals averaging fewer than 1 write ops

```sql
select
  instance_id,
  round(minimum::numeric,2) as min_write_ops,
  round(maximum::numeric,2) as max_write_ops,
  round(average::numeric,2) as avg_write_ops,
  sample_count
from
  gcp_compute_disk_metric_write_ops_daily
where average < 1
order by
  instance_id;
```