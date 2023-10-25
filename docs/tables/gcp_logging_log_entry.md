# Table: gcp_logging_log_entry

In Google Cloud Platform (GCP), a logging log entry represents a single log event captured by GCP's logging service. It contains information about a specific occurrence or action that took place within a GCP resource or service. Each log entry contains various metadata and data fields that provide details about the event, such as the log severity, timestamp, log message, log name, resource information, and any additional structured data associated with the event.

**Important notes:**

- For improved performance, it is advised that you use the optional qual `timestamp` to limit the result set to a specific time period.
- This table supports optional quals. Queries with optional quals are optimised to use Logging filters. Optional quals are supported for the following columns:
  - `resource_type`
  - `severity`
  - `log_name`
  - `span_id`
  - `text_payload`
  - `receive_timestamp`
  - `timestamp`
  - `trace`
  - `log_entry_operation_id`
  - `filter`

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

### List log entries that occurred between five to ten minutes ago

```sql
select
  log_name,
  insert_id,
  receive_timestamp,
  trace_sampled,
  severity,
  resource_type
from
  gcp_logging_log_entry
where
  log_name = 'projects/parker-abbb/logs/cloudaudit.googleapis.com%2Factivity'
and
  timestamp between (now() - interval '10 minutes') and (now() - interval '5 minutes')
order by
  receive_timestamp asc;
```

### Get the last log entries

```sql
select
  log_name,
  insert_id,
  log_entry_operation_last,
  receive_timestamp,
  resource_type,
  severity,
  text_payload
from
  gcp_logging_log_entry
where log_entry_operation_last;
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
  log_name = 'projects/parker-abbb/logs/cloudaudit.googleapis.com%2Factivity';
```

## Filter examples

For more information on Logging log entry filters, please refer to [Filter Pattern Syntax](https://cloud.google.com/logging/docs/view/logging-query-language).

### List log entries of Compute Engine VMs with serverity error

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
  filter = 'resource.type = "gce_instance" AND (severity = ERROR OR "error")';
```

### List events originating from a specific IP address range that occurred over the last hour

```sql
select
  log_name,
  insert_id,
  receive_timestamp,
  resource_type,
  severity,
  timestamp,
  resource_labels
from
  gcp_logging_log_entry
where
  filter = 'logName = "projects/my_project/logs/my_log" AND ip_in_net(jsonPayload.realClientIP, "10.1.2.0/24")'
  and timestamp >= now() - interval '1 hour'
order by
  receive_timestamp asc;
```

### Get prot payload details of each log entry

```sql
select
  insert_id,
  log_name,
  proto_payload -> 'authenticationInfo' as authentication_info,
  proto_payload -> 'authorizationInfo' as authorization_info,
  proto_payload -> 'serviceName' as service_name,
  proto_payload -> 'resourceName' as resource_name,
  proto_payload ->> '@type' as proto_payload_type,
  proto_payload ->> 'methodName' as method_name,
  proto_payload ->> 'callerIp' as caller_ip
from
  gcp_logging_log_entry
where
  filter = 'resource.type = "gce_instance" AND (severity = ERROR OR "error")';
```