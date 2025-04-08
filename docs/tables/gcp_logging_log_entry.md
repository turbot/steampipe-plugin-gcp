---
title: "Steampipe Table: gcp_logging_log_entry - Query Google Cloud Logging Log Entries using SQL"
description: "Allows users to query Log Entries in Google Cloud Logging, providing insights into system events, application behavior, and user actions."
folder: "Cloud Logging"
---

# Table: gcp_logging_log_entry - Query Google Cloud Logging Log Entries using SQL

Google Cloud Logging is a service that stores logs from your applications, systems, and services on Google Cloud Platform (GCP). It allows you to analyze and export logs to various destinations for long-term storage or further analysis. Google Cloud Logging helps you understand how your applications and services are performing and how they are being used.

## Table Usage Guide

The `gcp_logging_log_entry` table provides insights into Log Entries within Google Cloud Logging. As a System Administrator, explore log entry-specific details through this table, including severity, timestamp, and associated metadata. Utilize it to uncover information about system events, application behavior, and user actions, which can be useful for debugging, auditing, and performance optimization.

**Important Notes:**
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
  - `operation_id`
  - `filter`: For additional details regarding the filter string, please refer to the documentation at https://cloud.google.com/logging/docs/view/logging-query-language.

## Examples

### Basic info
Explore the operations recorded in the Google Cloud Platform's logging service to gain insights into the sequence and timing of events. This can help you understand the operational flow and identify potential issues based on event timings.

```sql+postgres
select
  log_name,
  insert_id,
  log_entry_operation_first,
  log_entry_operation_id,
  receive_timestamp
from
  gcp_logging_log_entry;
```

```sql+sqlite
select
  log_name,
  insert_id,
  operation_id,
  receive_timestamp
from
  gcp_logging_log_entry;
```

### Get log entries by resource type
Analyze the settings to understand the various log entries associated with a specific type of audited resource. This can be particularly useful for pinpointing operational issues or potential security concerns tied to that resource type.

```sql+postgres
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

```sql+sqlite
select
  log_name,
  insert_id,
  resource_type,
  span_id,
  text_payload
from
  gcp_logging_log_entry
where
  resource_type = 'audited_resource';
```

### List log entries with NOTICE severity
Discover the segments that have log entries with a NOTICE severity. This can help monitor system activity and identify potential issues that may not be critical but are still noteworthy.

```sql+postgres
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

```sql+sqlite
select
  log_name,
  insert_id,
  severity,
  span_id,
  timestamp
from
  gcp_logging_log_entry
where
  severity = 'NOTICE';
```

### List log entries in last 30 days
Explore the recent activity within the last month by identifying the log entries. This helps in monitoring system performance and tracking changes, thereby aiding in system maintenance and troubleshooting.

```sql+postgres
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

```sql+sqlite
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
  timestamp >= datetime('now', '-30 day');
```

### List log entries that occurred between five to ten minutes ago
Explore the recent activities in your project by identifying log entries that occurred within a specific time frame, in this case, between five to ten minutes ago. This can help in monitoring real-time activities or detecting any irregularities within a short span of time.

```sql+postgres
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

```sql+sqlite
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
  timestamp between (datetime('now', '-10 minutes')) and (datetime('now', '-5 minutes'))
order by
  receive_timestamp asc;
```

### Get the last log entries
Explore the most recent activities in your system by checking the last entries in the logs. This can help you monitor system operations, identify potential issues, and maintain a secure and efficient environment.

```sql+postgres
select
  log_name,
  insert_id,
  operation ->> 'Last' as log_entry_operation_last,
  receive_timestamp,
  resource_type,
  severity,
  text_payload
from
  gcp_logging_log_entry
where
  (operation ->> 'Last')::boolean;
```

```sql+sqlite
select
  log_name,
  insert_id,
  json_extract(operation, '$.Last') as log_entry_operation_last,
  receive_timestamp,
  resource_type,
  severity,
  text_payload
from
  gcp_logging_log_entry
where
  json_extract(operation, '$.Last') = 'true';
```

### Filter log entries by log name
Explore the specific log entries by defining a particular log name. This can help in narrowing down the search and making the process of monitoring and troubleshooting more efficient.

```sql+postgres
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

```sql+sqlite
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

### Get split details of each log entry
Extracting detailed information about specific log entries in a structured and relational manner. It allows for a deeper analysis of the logs by providing contextual information like the sequence of the log entry

```sql+postgres
select
  log_name,
  insert_id,
  split ->> 'Index' as split_index,
  split ->> 'TotalSplits' as total_splits,
  split ->> 'Uid' as split_uid
from
  gcp_logging_log_entry;
```

```sql+sqlite
select
  log_name,
  insert_id,
  json_extract(split, '$.Index') as split_index,
  json_extract(split, '$.TotalSplits') as total_splits,
  json_extract(split, '$.Uid') as split_uid
from
  gcp_logging_log_entry;
```

### Get operation details of each log entry
Retrieve the specifics of operation-related details from log entry records. This query can be instrumental in acquiring information regarding the initial operation, concluding operation, and the source of each operation.

```sql+postgres
select
  log_name,
  insert_id,
  operation_id,
  operation ->> 'Producer' as operation_producer,
  operation ->> 'First' as operation_first,
  operation ->> 'Last' as operation_last
from
  gcp_logging_log_entry;
```

```sql+sqlite
select
  log_name,
  insert_id,
  operation_id,
  json_extract(operation, '$.Producer') as operation_producer,
  json_extract(operation, '$.First') as operation_first,
  json_extract(operation, '$.Last') as operation_last
from
  gcp_logging_log_entry;
```

## Filter examples

For more information on Logging log entry filters, please refer to [Filter Pattern Syntax](https://cloud.google.com/logging/docs/view/logging-query-language).

### List log entries of Compute Engine VMs with serverity error
Discover the segments that have logged errors on your Google Compute Engine virtual machines. This query is beneficial in identifying and troubleshooting system faults, ensuring smooth operation of your VMs.

```sql+postgres
select
  log_name,
  insert_id,
  receive_timestamp,
  resource_type,
  severity
from
  gcp_logging_log_entry
where
  filter = 'resource.type = "gce_instance" and (severity = ERROR OR "error")';
```

```sql+sqlite
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
  filter = 'resource.type = "gce_instance" and (severity = "ERROR" OR severity = "error")';
```

### List events originating from a specific IP address range that occurred over the last hour
This query allows you to monitor and analyze events from a specific IP address range that have occurred in the last hour. It's a useful tool for real-time network security monitoring, helping to promptly identify unusual activity or potential security threats.

```sql+postgres
select
  log_name,
  insert_id,
  receive_timestamp,
  resource_type,
  severity,
  timestamp
from
  gcp_logging_log_entry
where
  filter = 'logName = "projects/my_project/logs/my_log" AND ip_in_net(jsonPayload.realClientIP, "10.1.2.0/24")'
  and timestamp >= now() - interval '1 hour'
order by
  receive_timestamp asc;
```

### Get proto payload details of each log entry
The query is useful for extracting specific information from log entries in a GCP logging system, particularly for entries related to Google Compute Engine (GCE) instances with errors. Extracting specific information from log entries in a GCP logging system, particularly for entries related to Google Compute Engine (GCE) instances with errors.

```sql+postgres
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

```sql+sqlite
select
  insert_id,
  log_name,
  json_extract(proto_payload, '$.authenticationInfo') AS authentication_info,
  json_extract(proto_payload, '$.authorizationInfo') AS authorization_info,
  json_extract(proto_payload, '$.serviceName') AS service_name,
  json_extract(proto_payload, '$.resourceName') AS resource_name,
  json_extract(proto_payload, '$.@type') AS proto_payload_type,
  json_extract(proto_payload, '$.methodName') AS method_name,
  json_extract(proto_payload, '$.callerIp') AS caller_ip
from
  gcp_logging_log_entry
where
  filter = 'resource.type = "gce_instance" AND (severity = ERROR OR severity = "error")';
```