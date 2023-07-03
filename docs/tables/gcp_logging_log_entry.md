
# Table: gcp_logging_log_entry

AWS Logging Log Entry refers to a single log event or entry in a log stream within an AWS service's logging system. It contains information about a specific occurrence or activity that is being logged. Log entries typically include details such as timestamp, log message, log level, request ID, and other relevant metadata.

## Examples

### Basic info

```sql
select
  log_name,
  insert_id,
  log_entry_operation_first,
  log_entry_operation_id,
  receive_timestamp
from
  gcp_logging_log_entry;
```

### Get log entries by resource type

```sql
select
  log_name,
  insert_id,
  log_entry_operation_first,
  log_entry_operation_last,
  resource_type,
  span_id,
  text_payload
from
  gcp_logging_log_entry
where
  resource_type = 'audited_resource';
```

### List log entries with NOTICE severity

```sql
select
  log_name,
  insert_id,
  resource_type,
  severity,
  span_id,
  timestamp
from
  gcp_logging_log_entry
where
  severity = 'NOTICE';
```

### List log entries in last 30 days

```sql
select
  log_name,
  insert_id,
  receive_timestamp,
  trace_sampled,
  span_id,
  timestamp
from
  gcp_logging_log_entry
where
  timestamp >= now() - interval '30' day;
```

### Filter log entries by log name

```sql
select
  log_name,
  insert_id,
  log_entry_operation_first,
  log_entry_operation_last,
  receive_timestamp,
  resource_type,
  severity
from
  gcp_logging_log_entry
where
  log_name = 'projects/parker-aaa/logs/cloudaudit.googleapis.com%2Factivity';
```