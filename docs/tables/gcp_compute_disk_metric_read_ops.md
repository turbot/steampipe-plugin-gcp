# Table: gcp_disk_metric_read_ops_daily

Google cloud Monitoring Metrics provide data about the performance of your systems. The `gcp_disk_metric_read_ops_daily` table provides metric statistics at 24 hour intervals for the past year.

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
  gcp_disk_metric_read_ops_daily
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
  gcp_disk_metric_read_ops_daily
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
  gcp_disk_metric_read_ops_daily
where average < 1
order by
  name;
```