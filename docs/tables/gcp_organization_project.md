---
title: "Steampipe Table: gcp_organization_project - Query Google Cloud Platform Projects using SQL"
description: "Allows users to query Projects in Google Cloud Platform, specifically providing details about the project's ID, name, labels, and lifecycle state."
folder: "Organization"
---

# Table: gcp_organization_project - Query Google Cloud Platform Projects using SQL

**Note: this table is a variant of the `gcp_project` table which does not filter on the GCP project attached to connection, and thus, will return all projects that the credentials used by the connection have access to. Using this table in aggregator connections can produce unexpected duplicate results.**

A Google Cloud Platform Project acts as an organizational unit within GCP where resources are allocated. It is used to group resources that belong to the same logical application or business unit. Each project is linked to a billing account and can have users, roles, and permissions assigned to it.

## Table Usage Guide

The `gcp_organization_project` table provides insights into Projects within Google Cloud Platform. As a DevOps engineer, explore project-specific details through this table, including ID, name, labels, and lifecycle state. Utilize it to uncover information about projects, such as their associated resources, user roles, permissions, and billing details.

## Examples

### Basic info
Explore which Google Cloud Platform projects are active, by looking at their lifecycle state and creation time. This can help you manage resources effectively and keep track of ongoing projects.

```sql+postgres
select
  name,
  project_id,
  project_number,
  lifecycle_state,
  create_time
from
  gcp_organization_project;
```

```sql+sqlite
select
  name,
  project_id,
  project_number,
  lifecycle_state,
  create_time
from
  gcp_organization_project;
```

### Get access approval settings for all projects
Explore the access approval settings across your various projects. This can help you understand and manage permissions and approvals more effectively.

```sql+postgres
select
  name,
  jsonb_pretty(access_approval_settings) as access_approval_settings
from
  gcp_organization_project;
```

```sql+sqlite
select
  name,
  access_approval_settings
from
  gcp_organization_project;
```

### Get parent and organization ID for all projects
Get the organization ID across your various projects.

```sql+postgres
select
  project_id,
  parent ->> 'id' as parent_id,
  parent ->> 'type' as parent_type,
  case when jsonb_array_length(ancestors) > 1 then ancestors -> -1 -> 'resourceId' ->> 'id' else null end as organization_id
from
  gcp_project;
```

```sql+sqlite
select
  project_id,
  parent ->> 'id' as parent_id,
  parent ->> 'type' as parent_type,
  case when json_array_length(ancestors) > 1 then ancestors -> -1 -> 'resourceId' ->> 'id' else null end as organization_id
from
  gcp_project;
```
