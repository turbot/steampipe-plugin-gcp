# Table: gcp_disk_metric_read_ops

Google cloud Monitoring Metrics provide data about the performance of your systems. The `gcp_disk_metric_read_ops` table provides metric statistics at 5 min intervals for the past year.

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
  gcp_disk_metric_read_ops
order by
  instance_id;
```

### Intervals averaging over 100 read ops

```sql
select
  instance_id,
  round(minimum::numeric,2) as min_read_ops,
  round(maximum::numeric,2) as max_read_ops,
  round(average::numeric,2) as avg_read_ops,
  sample_count
from
  gcp_disk_metric_read_ops
where average > 10
order by
  instance_id;
```

### Intervals averaging fewer than 10 read ops

```sql
select
  instance_id,
  round(minimum::numeric,2) as min_read_ops,
  round(maximum::numeric,2) as max_read_ops,
  round(average::numeric,2) as avg_read_ops,
  sample_count
from
  gcp_disk_metric_read_ops
where average < 10
order by
  instance_id;
```
