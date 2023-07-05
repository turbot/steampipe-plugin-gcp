# Table: gcp_logging_log_entry

In Google Cloud Platform (GCP), a logging log entry represents a single log event captured by GCP's logging service. It contains information about a specific occurrence or action that took place within a GCP resource or service. Each log entry contains various metadata and data fields that provide details about the event, such as the log severity, timestamp, log message, log name, resource information, and any additional structured data associated with the event.

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