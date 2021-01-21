# Table: gcp_logging_exclusion

Exclusion filter lets you exclude logs on the logs sinks that route logs to your Logs Buckets

### Basic info

```sql
select
  name,
  disabled,
  filter,
  description
from
  gcp_logging_exclusion;
```


### List of exclusions which are disabled

```sql
select
  name,
  disabled
from
  gcp_logging_exclusion
where
  disabled;
```
