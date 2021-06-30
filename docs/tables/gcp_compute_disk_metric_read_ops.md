# Table: gcp_compute_disk_metric_read_ops

GCP Monitoring metrics provide data about the performance of your systems. The `gcp_compute_disk_metric_read_ops` table provides metric statistics at 5 minute intervals for the most recent 5 days.

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
  gcp_compute_disk_metric_read_ops
order by
  name;
```

### Intervals averaging over 100 read ops

```sql
select
  name,
  round(minimum::numeric,2) as min_read_ops,
  round(maximum::numeric,2) as max_read_ops,
  round(average::numeric,2) as avg_read_ops,
  sample_count
from
  gcp_compute_disk_metric_read_ops
where average > 10
order by
  name;
```

### Intervals averaging fewer than 10 read ops

```sql
select
  name,
  round(minimum::numeric,2) as min_read_ops,
  round(maximum::numeric,2) as max_read_ops,
  round(average::numeric,2) as avg_read_ops,
  sample_count
from
  gcp_compute_disk_metric_read_ops
where average < 10
order by
  name;
```
