# Table: gcp_cloud_run_service

Google Cloud Run is a fully managed compute platform offered by Google Cloud Platform (GCP) that is designed for running containerized applications. It allows developers to deploy containerized applications quickly and easily without having to manage the underlying infrastructure. Cloud Run abstracts away many of the complexities of managing servers and scaling applications, making it an excellent choice for building and deploying microservices, APIs, web applications, and more.

## Examples

### Basic info

```sql
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

```sql
select
  launch_stage,
  count(*)
from
  gcp_cloud_run_service
group by
  launch_stage;
```

### List cloud-run services that are reconciling

```sql
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

### List services created in the last 30 days

```sql
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

### List services of ingress type INGRESS_TRAFFIC_ALL

```sql
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

```sql
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

### Get template details of services

```sql
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

### Get target traffic details of services

```sql
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
