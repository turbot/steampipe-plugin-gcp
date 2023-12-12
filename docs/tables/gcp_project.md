---
title: "Steampipe Table: gcp_project - Query Google Cloud Platform Projects using SQL"
description: "Allows users to query Projects in Google Cloud Platform, specifically providing details about the project's ID, name, labels, and lifecycle state."
---

# Table: gcp_project - Query Google Cloud Platform Projects using SQL

A Google Cloud Platform Project acts as an organizational unit within GCP where resources are allocated. It is used to group resources that belong to the same logical application or business unit. Each project is linked to a billing account and can have users, roles, and permissions assigned to it.

## Table Usage Guide

The `gcp_project` table provides insights into Projects within Google Cloud Platform. As a DevOps engineer, explore project-specific details through this table, including ID, name, labels, and lifecycle state. Utilize it to uncover information about projects, such as their associated resources, user roles, permissions, and billing details.

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
  gcp_project;
```

```sql+sqlite
select
  name,
  project_id,
  project_number,
  lifecycle_state,
  create_time
from
  gcp_project;
```

### Get access approval settings for all projects
Explore the access approval settings across your various projects. This can help you understand and manage permissions and approvals more effectively.

```sql+postgres
select
  name,
  jsonb_pretty(access_approval_settings) as access_approval_settings
from
  gcp_project;
```

```sql+sqlite
select
  name,
  access_approval_settings
from
  gcp_project;
```