# Table: gcp_logging_metric

Logs-based metrics are Cloud Monitoring metrics that are based on the content of log entries.

### Filter info of each metric

```sql
select
  name as metric_name,
  description,
  filter
from
  gcp_logging_metric;
```


### Bucket configuration details of the logging metrics

```sql
select
  name,
  exponential_buckets_options_growth_factor,
  exponential_buckets_options_num_finite_buckets,
  exponential_buckets_options_scale,
  linear_buckets_options_num_finite_buckets,
  linear_buckets_options_offset,
  linear_buckets_options_width,
  explicit_buckets_options_bounds
from
  gcp_logging_metric;
```
