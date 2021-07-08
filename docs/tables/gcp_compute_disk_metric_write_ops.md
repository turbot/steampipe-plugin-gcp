# Table: gcp_compute_disk_metric_write_ops

GCP Monitoring metrics provide data about the performance of your systems. The `gcp_compute_disk_metric_write_ops` table provides metric statistics at 5 minute intervals for the most recent 5 days.

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
  gcp_compute_disk_metric_write_ops
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
  gcp_compute_disk_metric_write_ops
where average > 10
order by
  name;
```
