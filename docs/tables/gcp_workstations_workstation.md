---
title: "Steampipe Table: gcp_workstations_workstation - Query Google Cloud Workstations using SQL"
description: "Query Google Cloud Workstations, including configuration, state, and metadata details for cloud-based development environments."
folder: "Workstations"
---

# Table: gcp_workstations_workstation - Query Google Cloud Workstations using SQL

Google Cloud Workstations is a fully managed solution that provides secure, cloud-based development environments. Cloud Workstations allows you to create and manage workstations for your development teams, providing a consistent environment with all the necessary tools and configurations. It enables developers to access their workstations from anywhere, scales automatically, and integrates with other Google Cloud Platform services.

## Table Usage Guide

The `gcp_workstations_workstation` table provides insights into Workstations within Google Cloud. As a developer or platform engineer, explore workstation-specific details through this table, including configurations, states, and metadata. Use it to uncover information about workstations, such as their current state, associated configurations, and the resources allocated to each workstation.

## Examples

### Basic info
Explore the basic details of your Google Cloud Workstations, including their names, states, and creation times. This information helps you understand the configuration and status of your workstations, which is useful for managing cloud development environments.

```sql+postgres
select
  name,
  display_name,
  state,
  host,
  create_time,
  uid,
  location,
  project
from
  gcp_workstations_workstation;
```

```sql+sqlite
select
  name,
  display_name,
  state,
  host,
  create_time,
  uid,
  location,
  project
from
  gcp_workstations_workstation;
```

### List running workstations
Identify workstations that are currently in a running state to monitor active development environments and understand resource utilization.

```sql+postgres
select
  name,
  display_name,
  state,
  host,
  start_time,
  location
from
  gcp_workstations_workstation
where
  state = 'STATE_RUNNING';
```

```sql+sqlite
select
  name,
  display_name,
  state,
  host,
  start_time,
  location
from
  gcp_workstations_workstation
where
  state = 'STATE_RUNNING';
```

### List workstations that are being updated
Analyze which workstations are currently reconciling (being updated). This is useful for tracking environments undergoing changes.

```sql+postgres
select
  name,
  display_name,
  state,
  reconciling,
  update_time
from
  gcp_workstations_workstation
where
  reconciling;
```

```sql+sqlite
select
  name,
  display_name,
  state,
  reconciling,
  update_time
from
  gcp_workstations_workstation
where
  reconciling = 1;
```

### List workstations created in the last 30 days
Discover workstations provisioned in the past 30 days to understand recent environment activity.

```sql+postgres
select
  name,
  display_name,
  state,
  create_time,
  location
from
  gcp_workstations_workstation
where
  create_time >= now() - interval '30' day;
```

```sql+sqlite
select
  name,
  display_name,
  state,
  create_time,
  location
from
  gcp_workstations_workstation
where
  create_time >= datetime('now', '-30 day');
```

### Count of workstations by state
Determine the distribution of workstations by their current state to aid in resource planning and management.

```sql+postgres
select
  state,
  count(*)
from
  gcp_workstations_workstation
group by
  state;
```

```sql+sqlite
select
  state,
  count(*)
from
  gcp_workstations_workstation
group by
  state;
```

### Get environment variables for workstations
Examine the environment variables configured for each workstation to understand runtime settings.

```sql+postgres
select
  name,
  display_name,
  jsonb_pretty(env) as environment_variables
from
  gcp_workstations_workstation;
```

```sql+sqlite
select
  name,
  display_name,
  env as environment_variables
from
  gcp_workstations_workstation;
```

### Get IAM policy for workstations
Retrieve IAM policy bindings associated with workstations to understand access control configurations.

```sql+postgres
select
  name,
  i -> 'condition' as condition,
  i -> 'members' as members,
  i ->> 'role' as role
from
  gcp_workstations_workstation,
  jsonb_array_elements(iam_policy -> 'bindings') as i;
```

```sql+sqlite
select
  name,
  json_extract(i.value, '$.condition') as condition,
  json_extract(i.value, '$.members') as members,
  json_extract(i.value, '$.role') as role
from
  gcp_workstations_workstation,
  json_each(json_extract(iam_policy, '$.bindings')) as i;
```

### List deleted workstations
Identify workstations that have been soft-deleted but not yet permanently removed to track recently removed resources.

```sql+postgres
select
  name,
  display_name,
  state,
  delete_time,
  location
from
  gcp_workstations_workstation
where
  delete_time is not null;
```

```sql+sqlite
select
  name,
  display_name,
  state,
  delete_time,
  location
from
  gcp_workstations_workstation
where
  delete_time is not null;
```