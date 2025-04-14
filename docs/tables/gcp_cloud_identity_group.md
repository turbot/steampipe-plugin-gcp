---
title: "Steampipe Table: gcp_cloud_identity_group - Query Google Cloud Platform Cloud Identity Groups using SQL"
description: "Allows users to query Cloud Identity Groups in Google Cloud Platform, specifically the details of each group including its name, description, dynamic settings, and labels."
folder: "Cloud Identity"
---

# Table: gcp_cloud_identity_group - Query Google Cloud Platform Cloud Identity Groups using SQL

Google Cloud Identity Groups is a service within Google Cloud Platform that helps you manage access to your resources. It allows you to create groups for your workspace, manage group memberships, and provide access to resources based on group membership. Google Cloud Identity Groups makes it easier to manage access to resources at scale.

## Table Usage Guide

The `gcp_cloud_identity_group` table provides insights into Cloud Identity Groups within Google Cloud Platform. As a security engineer, explore group-specific details through this table, including group names, descriptions, dynamic settings, and labels. Utilize it to uncover information about groups, such as the group's access permissions, the members within the group, and the resources accessible to the group.

**Important Notes**
- You must specify the parent resource in the `where` clause (`where parent='C046psxkn'`) to list the identity groups.

## Examples

### Basic info
Explore which Google Cloud Identity groups are associated with a specific parent group. This can be useful for understanding group hierarchies and the distribution of resources within your Google Cloud project.

```sql+postgres
select
  name,
  display_name,
  description,
  create_time,
  location,
  project
from
  gcp_cloud_identity_group
where
  parent = 'C046psxkn';
```

```sql+sqlite
select
  name,
  display_name,
  description,
  create_time,
  location,
  project
from
  gcp_cloud_identity_group
where
  parent = 'C046psxkn';
```

### Get details for a specific group
Explore the specifics of a particular group in Google Cloud Platform's Cloud Identity service. This can be useful in understanding the group's creation time, location, and associated project, aiding in effective group management and security oversight.

```sql+postgres
select
  name,
  display_name,
  description,
  create_time,
  location,
  project
from
  gcp_cloud_identity_group
where
  name = 'group_name';
```

```sql+sqlite
select
  name,
  display_name,
  description,
  create_time,
  location,
  project
from
  gcp_cloud_identity_group
where
  name = 'group_name';
```

### Get dynamic group settings
Analyze the settings to understand the status and configuration of dynamic groups within a specific project in Google Cloud Identity. This can be useful for managing and monitoring group membership based on user-defined rules.

```sql+postgres
select
  name,
  display_name,
  dynamic_group_metadata ->> 'Status' as dynamic_group_status,
  queries ->> 'Query' as dynamic_group_query,
  queries ->> 'ResourceType' as dynamic_group_query_resource_type,
  project
from
  gcp_cloud_identity_group,
  jsonb_array_elements(dynamic_group_metadata -> 'Queries') as queries
where
  parent = 'C046psxkn';
```

```sql+sqlite
select
  g.name,
  g.display_name,
  json_extract(g.dynamic_group_metadata, '$.Status') as dynamic_group_status,
  json_extract(queries.value, '$.Query') as dynamic_group_query,
  json_extract(queries.value, '$.ResourceType') as dynamic_group_query_resource_type,
  g.project
from
  gcp_cloud_identity_group as g,
  json_each(json_extract(g.dynamic_group_metadata, '$.Queries')) as queries
where
  g.parent = 'C046psxkn';
```

### List groups created in the last 7 days
Explore which groups have been formed within the last week in the GCP Cloud Identity service. This can be useful for keeping track of recent group additions and ensuring proper access controls are in place.

```sql+postgres
select
  name,
  display_name,
  description,
  create_time,
  location,
  project
from
  gcp_cloud_identity_group
where
  parent = 'C046psxkn'
  and create_time > now() - interval '7' day;
```

```sql+sqlite
select
  name,
  display_name,
  description,
  create_time,
  location,
  project
from
  gcp_cloud_identity_group
where
  parent = 'C046psxkn'
  and create_time > datetime('now', '-7 day');
```