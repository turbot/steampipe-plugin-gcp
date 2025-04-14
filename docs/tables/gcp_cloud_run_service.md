---
title: "Steampipe Table: gcp_cloud_run_service - Query GCP Cloud Run Services using SQL"
description: "Allows users to query GCP Cloud Run Services, specifically the configurations, revisions, and routes, providing insights into the application's traffic flow and deployment history."
folder: "Cloud Run"
---

# Table: gcp_cloud_run_service - Query GCP Cloud Run Services using SQL

Google Cloud Run is a managed compute platform that enables you to run stateless containers that are invocable via HTTP requests. Cloud Run is serverless: it abstracts away all infrastructure management, so you can focus on what matters most — building great applications. It automatically scales up or down from zero to N depending on traffic.

## Table Usage Guide

The `gcp_cloud_run_service` table provides insights into Cloud Run services within Google Cloud Platform (GCP). As a developer or DevOps engineer, explore service-specific details through this table, including configurations, revisions, and routes. Utilize it to uncover information about services, such as the application's traffic flow, deployment history, and the current state of the service.

## Examples

### Basic info
Explore the basic details of your Google Cloud Run services, including their names, descriptions, and client versions. This information can help you understand the configuration and status of your services, which is useful for managing and optimizing your cloud resources.

```sql+postgres
select
  name,
  description,
  client,
  client_version,
  create_time,
  creator,
  generation,
  launch_stage
from
  gcp_cloud_run_service;
```

```sql+sqlite
select
  name,
  description,
  client,
  client_version,
  create_time,
  creator,
  generation,
  launch_stage
from
  gcp_cloud_run_service;
```

### Count of services by launch stage
Determine the distribution of services based on their launch stages. This can help in understanding how many services are in different stages of their lifecycle, providing insights for resource allocation and strategic planning.

```sql+postgres
select
  launch_stage,
  count(*)
from
  gcp_cloud_run_service
group by
  launch_stage;
```

```sql+sqlite
select
  launch_stage,
  count(*)
from
  gcp_cloud_run_service
group by
  launch_stage;
```

### List cloud-run services that are reconciling
Analyze the settings to understand which cloud-run services are currently in the process of reconciling. This can be useful for tracking and managing services that may be undergoing changes or updates.

```sql+postgres
select
  name,
  description,
  client,
  client_version,
  create_time,
  creator,
  generation,
  launch_stage,
  reconciling
from
  gcp_cloud_run_service
where
  reconciling;
```

```sql+sqlite
select
  name,
  description,
  client,
  client_version,
  create_time,
  creator,
  generation,
  launch_stage,
  reconciling
from
  gcp_cloud_run_service
where
  reconciling = 1;
```

### List services created in the last 30 days
Discover the services that were established in the past 30 days to gain insights into recent activities and understand the context of their creation. This could be useful in tracking the growth of services over time or identifying any unexpected or unauthorized service creation.

```sql+postgres
select
  name,
  description,
  create_time,
  creator,
  launch_stage
from
  gcp_cloud_run_service
where
  create_time >= now() - interval '30' day;
```

```sql+sqlite
select
  name,
  description,
  create_time,
  creator,
  launch_stage
from
  gcp_cloud_run_service
where
  create_time >= datetime('now', '-30 day');
```

### List services of ingress type INGRESS_TRAFFIC_ALL
Analyze the settings to understand which cloud run services are configured to allow all types of ingress traffic. This can be useful for assessing potential security risks associated with unrestricted ingress access.

```sql+postgres
select
  name,
  description,
  client,
  client_version,
  create_time,
  ingress
from
  gcp_cloud_run_service
where
  ingress = 'INGRESS_TRAFFIC_ALL';
```

```sql+sqlite
select
  name,
  description,
  client,
  client_version,
  create_time,
  ingress
from
  gcp_cloud_run_service
where
  ingress = 'INGRESS_TRAFFIC_ALL';
```

### Get condition details of services
This example allows you to gain insights into the status and condition details of various services in the Google Cloud Run environment. It can be used to understand the health of services, the reasons for their current state, and when they last transitioned, which can assist in troubleshooting and maintaining service stability.

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
  gcp_cloud_run_service,
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
  gcp_cloud_run_service,
  json_each(conditions) as c;
```

### Get associated members or principals, with a role of services
Attaching an Identity and Access Management (IAM) policy to a Google Cloud Run service involves setting permissions for that particular service. Google Cloud Run services use IAM for access control, and by configuring IAM policies, you can define who has what type of access to your Cloud Run services.

```sql+postgres
select
  name,
  i -> 'Condition' as condition,
  i -> 'Members' as members,
  i ->> 'Role' as role
from
  gcp_cloud_run_service,
  jsonb_array_elements(iam_policy -> 'Bindings') as i;
```

```sql+sqlite
select
  name,
  json_extract(i.value, '$.Condition') as condition,
  json_extract(i.value, '$.Members') as members,
  json_extract(i.value, '$.Role') as role
from
  gcp_cloud_run_service,
  json_each(json_extract(iam_policy, '$.Bindings')) as i;
```

### Get template details of services
Explore the various attributes of your cloud-based services, such as encryption keys, container details, and scaling parameters. This query is useful to gain an understanding of your service configurations and identify areas for potential adjustments or enhancements.

```sql+postgres
select
  name,
  template ->> 'Annotations' as template_annotations,
  template ->> 'Containers' as containers,
  template ->> 'EncryptionKey' as encryption_key,
  template ->> 'ExecutionEnvironment' as execution_environment,
  template ->> 'Revision' as revision,
  template ->> 'Scaling' as scaling,
  template ->> 'ServiceAccount' as service_account,
  template ->> 'SessionAffinity' as session_affinity,
  template ->> 'Timeout' as timeout,
  template ->> 'Volumes' as volumes,
  template ->> 'VpcAccess' as vpc_access
from
  gcp_cloud_run_service;
```

```sql+sqlite
select
  name,
  json_extract(template, '$.Annotations') as template_annotations,
  json_extract(template, '$.Containers') as containers,
  json_extract(template, '$.EncryptionKey') as encryption_key,
  json_extract(template, '$.ExecutionEnvironment') as execution_environment,
  json_extract(template, '$.Revision') as revision,
  json_extract(template, '$.Scaling') as scaling,
  json_extract(template, '$.ServiceAccount') as service_account,
  json_extract(template, '$.SessionAffinity') as session_affinity,
  json_extract(template, '$.Timeout') as timeout,
  json_extract(template, '$.Volumes') as volumes,
  json_extract(template, '$.VpcAccess') as vpc_access
from
  gcp_cloud_run_service;
```

### Get target traffic details of services
Gain insights into the distribution of traffic across different revisions and tags of your services. This is useful for understanding how your traffic is being balanced and identifying potential areas for optimization.

```sql+postgres
select
  name,
  t ->> 'Percent' as percent,
  t ->> 'Revision' as revision,
  t ->> 'Tag' as tag,
  t ->> 'Type' as type
from
  gcp_cloud_run_service,
  jsonb_array_elements(traffic) as t;
```

```sql+sqlite
select
  name,
  json_extract(t.value, '$.Percent') as percent,
  json_extract(t.value, '$.Revision') as revision,
  json_extract(t.value, '$.Tag') as tag,
  json_extract(t.value, '$.Type') as type
from
  gcp_cloud_run_service,
  json_each(traffic) as t;
```