# Table: gcp_sql_database_instance_metric_cpu_utilization

Google cloud Monitoring Metrics provide data about the performance of your systems. The `gcp_sql_database_instance_metric_cpu_utilization` table provides metric statistics at 24 hour intervals for the past year.

## Examples

### Basic info

```sql
select
  instance_id,
  minimum,
  maximum,
  average,
  sample_count,
  timestamp
from
  gcp_sql_database_instance_metric_cpu_utilization
order by
  instance_id;
```

### Intervals averaging over 80%

```sql
select
  instance_id,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  gcp_sql_database_instance_metric_cpu_utilization
where average > 80
order by
  instance_id;
```

### Intervals averaging < 1%

```sql
select
  instance_id,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  gcp_sql_database_instance_metric_cpu_utilization
where average < 1
order by
  instance_id;
```