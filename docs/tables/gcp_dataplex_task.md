---
title: "Steampipe Table: gcp_dataplex_task - Query GCP Dataplex Tasks using SQL"
description: "Allows users to query GCP Dataplex Tasks, providing detailed information about each task's configuration, execution, and status."
folder: "Dataplex"
---

# Table: gcp_dataplex_task - Query GCP Dataplex Tasks using SQL

GCP Dataplex Tasks are jobs that run on a scheduled basis or are triggered by specific events within a Dataplex Lake. These tasks can be used to manage and process data, including running custom Spark jobs or scheduled notebooks.

## Table Usage Guide

The `gcp_dataplex_task` table allows data engineers and cloud administrators to query and manage Dataplex Tasks within their GCP environment. You can retrieve information about a task’s configuration, execution status, and associated lake. This table is useful for monitoring and managing the state and execution of Dataplex Tasks.

## Examples

### List all Dataplex Tasks
Retrieve a list of all Dataplex Tasks in your GCP account to get an overview of your scheduled and triggered jobs.

```sql+postgres
select
  display_name,
  name,
  lake_name,
  state,
  create_time
from
  gcp_dataplex_task;
```

```sql+sqlite
select
  display_name,
  name,
  lake_name,
  state,
  create_time
from
  gcp_dataplex_task;
```

### Get details of task trigger specifications
This query extracts details about how tasks are triggered, including the schedule, type, start time, and maximum retries.

```sql+postgres
select
  name,
  trigger_spec ->> 'type' as trigger_type,
  trigger_spec ->> 'schedule' as trigger_schedule,
  trigger_spec ->> 'startTime' as trigger_start_time,
  trigger_spec ->> 'maxRetries' as trigger_max_retries
from
  gcp_dataplex_task;
```

```sql+sqlite
select
  name,
  json_extract(trigger_spec, '$.type') as trigger_type,
  json_extract(trigger_spec, '$.schedule') as trigger_schedule,
  json_extract(trigger_spec, '$.startTime') as trigger_start_time,
  json_extract(trigger_spec, '$.maxRetries') as trigger_max_retries
from
  gcp_dataplex_task;
```

### Get execution specifications for tasks
This query retrieves the execution specifications for each task, including the service account used, project, and maximum job execution lifetime.

```sql+postgres
select
  name,
  execution_spec ->> 'serviceAccount' as service_account,
  execution_spec ->> 'project' as project,
  execution_spec ->> 'maxJobExecutionLifetime' as max_execution_lifetime
from
  gcp_dataplex_task;
```

```sql+sqlite
select
  name,
  json_extract(execution_spec, '$.serviceAccount') as service_account,
  json_extract(execution_spec, '$.project') as project,
  json_extract(execution_spec, '$.maxJobExecutionLifetime') as max_execution_lifetime
from
  gcp_dataplex_task;
```

### Get the latest execution status of tasks
This query retrieves the latest execution status for each task, including the state, trigger, and any messages associated with the last job execution.

```sql+postgres
select
  name,
  execution_status -> 'latestJob' ->> 'state' as latest_job_state,
  execution_status -> 'latestJob' ->> 'trigger' as latest_job_trigger,
  execution_status -> 'latestJob' ->> 'message' as latest_job_message,
  execution_status -> 'latestJob' -> 'executionSpec' ->> 'kmsKey' as latest_job_kms_key
from
  gcp_dataplex_task;
```

```sql+sqlite
select
  name,
  json_extract(execution_status, '$.latestJob.state') as latest_job_state,
  json_extract(execution_status, '$.latestJob.trigger') as latest_job_trigger,
  json_extract(execution_status, '$.latestJob.message') as latest_job_message,
  json_extract(execution_status, '$.latestJob.executionSpec.kmsKey') as latest_job_kms_key
from
  gcp_dataplex_task;
```

### Dataplex tasks with their associated lakes
This is useful for understanding how tasks are distributed across different lakes in your Dataplex environment.

```sql+postgres
select
  t.name as task_name,
  t.state as task_state,
  t.create_time as task_create_time,
  l.name as lake_name,
  l.location as lake_location,
  l.state as lake_state
from
  gcp_dataplex_task as t
join
  gcp_dataplex_lake as l
on
  t.lake_name = l.name;
```

```sql+sqlite
select
  t.name as task_name,
  t.state as task_state,
  t.create_time as task_create_time,
  l.name as lake_name,
  l.location as lake_location,
  l.state as lake_state
from
  gcp_dataplex_task as t
join
  gcp_dataplex_lake as l
on
  t.lake_name = l.name;
```
