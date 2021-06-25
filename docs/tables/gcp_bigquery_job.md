# Table: gcp_bigquery_job

Jobs are actions that BigQuery runs on your behalf to load data, export data, query data, or copy data. Once a BigQuery job is created, it cannot be changed or deleted.

## Examples

### Basic info

```sql
select
  job_id,
  self_link,
  creation_time,
  location
from
  gcp_bigquery_job;
```

### List jobs which are currently running

```sql
select
  job_id,
  self_link,
  creation_time,
  location
from
  gcp_bigquery_job
where
  state = 'RUNNING';
```
