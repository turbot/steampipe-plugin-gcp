---
title: "Steampipe Table: gcp_compute_project_metadata - Query Google Cloud Compute Engine Project Metadata using SQL"
description: "Allows users to query Project Metadata in Google Cloud Compute Engine, specifically the project-wide metadata that includes common instance metadata and enable-oslogin metadata."
folder: "Compute"
---

# Table: gcp_compute_project_metadata - Query Google Cloud Compute Engine Project Metadata using SQL

Google Cloud Compute Engine Project Metadata is a set of data about a Google Cloud Compute Engine project. It includes common instance metadata that applies to all instances in the project, and enable-oslogin metadata that controls the OS Login feature for all instances in the project. This metadata can be used to configure or manage the behavior of the instances in the project.

## Table Usage Guide

The `gcp_compute_project_metadata` table provides insights into the metadata of projects within Google Cloud Compute Engine. As a Cloud Engineer, you can explore project-specific details through this table, including common instance metadata and enable-oslogin metadata. Utilize it to manage and configure the behavior of all instances in your projects, and to control the OS Login feature for all instances.

## Examples

### Basic info
Analyze the settings to understand the default service accounts and their creation timestamps within your Google Cloud Platform project. This can help you manage your resources and monitor any changes made over time.

```sql+postgres
select
  name,
  id,
  default_service_account,
  creation_timestamp
from
  gcp_compute_project_metadata;
```

```sql+sqlite
select
  name,
  id,
  default_service_account,
  creation_timestamp
from
  gcp_compute_project_metadata;
```

### Check if OS Login is enabled for Linux instances in the project
Determine the areas in which OS Login is not activated for Linux instances within a project. This insight can help enhance security by ensuring that all instances are properly configured for OS Login.

```sql+postgres
select
  name,
  id
from
  gcp_compute_project_metadata,
  jsonb_array_elements(common_instance_metadata -> 'items') as q
where
  common_instance_metadata -> 'items' @> '[{"key": "enable-oslogin"}]'
  and q ->> 'key' ilike 'enable-oslogin'
  and q ->> 'value' not ilike 'TRUE';
```

```sql+sqlite
select
  m.name,
  m.id
from
  gcp_compute_project_metadata as m,
  json_each(common_instance_metadata, '$.items') as q
where
  json_extract(common_instance_metadata, '$.items') like '%"key": "enable-oslogin"%'
  and json_extract(q.value, '$.key') like 'enable-oslogin'
  and json_extract(q.value, '$.value') not like 'TRUE';
```