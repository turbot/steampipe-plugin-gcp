# Table: gcp_logging_bucket

Logging buckets store the logs that are routed from other projects, folders, or organizations.

## Examples

### Basic info

```sql
select
  name,
  lifecycle_state,
  description,
  retention_days
from
  gcp_logging_bucket;
```

### List locked buckets

```sql
select
  name,
  locked
from
  gcp_logging_bucket
where
  locked;
```
