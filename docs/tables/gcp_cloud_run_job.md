---
title: "Steampipe Table: gcp_cloud_run_job - Query GCP Cloud Run Jobs using SQL"
description: "Allows users to query GCP Cloud Run Jobs, specifically the configurations and executions, providing insights into the application's deployment history."
folder: "Cloud Run"
---

# Table: gcp_cloud_run_job - Query GCP Cloud Run Jobs using SQL

Google Cloud Run is a managed compute platform that enables you to run stateless containers that are invocable via HTTP requests. Cloud Run is serverless: it abstracts away all infrastructure management, so you can focus on what matters most — building great applications. It automatically scales up or down from zero to N depending on traffic.

## Table Usage Guide

The `gcp_cloud_run_job` table provides insights into Cloud Run jobs within Google Cloud Platform (GCP). As a developer or DevOps engineer, explore job-specific details through this table, including configurations and executions. Utilize it to uncover information about jobs, such as the application's traffic flow, deployment history, and the current state of the job.

## Examples

### Basic info
Explore the basic details of your Google Cloud Run jobs, including their names, and client versions. This information can help you understand the configuration and status of your jobs, which is useful for managing and optimizing your cloud resources.

```sql+postgres
select
  name,
  client,
  client_version,
  create_time,
  creator,
  generation,
  launch_stage
from
  gcp_cloud_run_job;
```

```sql+sqlite
select
  name,
  client,
  client_version,
  create_time,
  creator,
  generation,
  launch_stage
from
  gcp_cloud_run_job;
```

### Count of jobs by launch stage
Determine the distribution of jobs based on their launch stages. This can help in understanding how many jobs are in different stages of their lifecycle, providing insights for resource allocation and strategic planning.

```sql+postgres
select
  launch_stage,
  count(*)
from
  gcp_cloud_run_job
group by
  launch_stage;
```

```sql+sqlite
select
  launch_stage,
  count(*)
from
  gcp_cloud_run_job
group by
  launch_stage;
```

### List cloud-run jobs that are reconciling
Analyze the settings to understand which cloud-run jobs are currently in the process of reconciling. This can be useful for tracking and managing jobs that may be undergoing changes or updates.

```sql+postgres
select
  name,
  client,
  client_version,
  create_time,
  creator,
  generation,
  launch_stage,
  reconciling
from
  gcp_cloud_run_job
where
  reconciling;
```

```sql+sqlite
select
  name,
  client,
  client_version,
  create_time,
  creator,
  generation,
  launch_stage,
  reconciling
from
  gcp_cloud_run_job
where
  reconciling = 1;
```

### List jobs created in the last 30 days
Discover the jobs that were established in the past 30 days to gain insights into recent activities and understand the context of their creation. This could be useful in tracking the growth of jobs over time or identifying any unexpected or unauthorized job creation.

```sql+postgres
select
  name,
  create_time,
  creator,
  launch_stage
from
  gcp_cloud_run_job
where
  create_time >= now() - interval '30' day;
```

```sql+sqlite
select
  name,
  create_time,
  creator,
  launch_stage
from
  gcp_cloud_run_job
where
  create_time >= datetime('now', '-30 day');
```

### Get condition details of jobs
This example allows you to gain insights into the status and condition details of various jobs in the Google Cloud Run environment. It can be used to understand the health of jobs, the reasons for their current state, and when they last transitioned, which can assist in troubleshooting and maintaining job stability.

```sql+postgres
select
  name,
  c ->> 'ExecutionReason' as execution_reason,
  c ->> 'LastTransitionTime' as last_transition_time,
  c ->> 'Message' as message,
  c ->> 'Reason' as reason,
  c ->> 'RevisionReason' as revision_reason,
  c ->> 'State' as state,
  c ->> 'Type' as type
from
  gcp_cloud_run_job,
  jsonb_array_elements(conditions) as c;
```

```sql+sqlite
select
  name,
  json_extract(c.value, '$.ExecutionReason') as execution_reason,
  json_extract(c.value, '$.LastTransitionTime') as last_transition_time,
  json_extract(c.value, '$.Message') as message,
  json_extract(c.value, '$.Reason') as reason,
  json_extract(c.value, '$.RevisionReason') as revision_reason,
  json_extract(c.value, '$.State') as state,
  json_extract(c.value, '$.Type') as type
from
  gcp_cloud_run_job,
  json_each(conditions) as c;
```

### Get associated members or principals, with a role of jobs
Attaching an Identity and Access Management (IAM) policy to a Google Cloud Run job involves setting permissions for that particular job. Google Cloud Run jobs use IAM for access control, and by configuring IAM policies, you can define who has what type of access to your Cloud Run jobs.

```sql+postgres
select
  name,
  i -> 'Condition' as condition,
  i -> 'Members' as members,
  i ->> 'Role' as role
from
  gcp_cloud_run_job,
  jsonb_array_elements(iam_policy -> 'Bindings') as i;
```

```sql+sqlite
select
  name,
  json_extract(i.value, '$.Condition') as condition,
  json_extract(i.value, '$.Members') as members,
  json_extract(i.value, '$.Role') as role
from
  gcp_cloud_run_job,
  json_each(json_extract(iam_policy, '$.Bindings')) as i;
```

### Get template details of jobs
Explore the various attributes of your cloud-based jobs, such as encryption keys, container details, and scaling parameters. This query is useful to gain an understanding of your job configurations and identify areas for potential adjustments or enhancements.

```sql+postgres
select
  name,
  template ->> 'Containers' as containers,
  template ->> 'EncryptionKey' as encryption_key,
  template ->> 'ExecutionEnvironment' as execution_environment,
  template ->> 'MaxRetries' as max_retries,
  template ->> 'ServiceAccount' as service_account,
  template ->> 'Timeout' as timeout,
  template ->> 'Volumes' as volumes,
  template ->> 'VpcAccess' as vpc_access
from
  gcp_cloud_run_job;
```

```sql+sqlite
select
  name,
  json_extract(template, '$.Containers') as containers,
  json_extract(template, '$.EncryptionKey') as encryption_key,
  json_extract(template, '$.ExecutionEnvironment') as execution_environment,
  json_extract(template, '$.MaxRetries') as max_retries,
  json_extract(template, '$.ServiceAccount') as service_account,
  json_extract(template, '$.Timeout') as timeout,
  json_extract(template, '$.Volumes') as volumes,
  json_extract(template, '$.VpcAccess') as vpc_access
from
  gcp_cloud_run_job;
```
