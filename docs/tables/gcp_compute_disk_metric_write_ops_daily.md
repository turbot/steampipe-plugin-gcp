# Table: gcp_compute_disk_metric_write_ops_daily

GCP Monitoring metrics provide data about the performance of your systems. The `gcp_compute_disk_metric_write_ops_daily` table provides metric statistics at 24 hour intervals for the last year.

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
  gcp_compute_disk_metric_write_ops_daily
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
  gcp_compute_disk_metric_write_ops_daily
where average > 10
order by
  name;
```

### Intervals averaging fewer than 1 write ops

```sql
select
  name,
  round(minimum::numeric,2) as min_write_ops,
  round(maximum::numeric,2) as max_write_ops,
  round(average::numeric,2) as avg_write_ops,
  sample_count
from
  gcp_compute_disk_metric_write_ops_daily
where average < 1
order by
  name;
```
