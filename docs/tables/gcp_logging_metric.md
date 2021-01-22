# Table: gcp_logging_metric

Logs-based metrics are Cloud Monitoring metrics that are based on the content of log entries.

### Filter info of each metric

```sql
select
  name as metric_name,
  filter
from
  gcp_logging_metric;
```


### List of all DELTA metric kind

```sql
select
  name,
  metric_descriptor_metric_kind
from
  gcp_logging_metric
where
  metric_descriptor_metric_kind = 'DELTA';
```
