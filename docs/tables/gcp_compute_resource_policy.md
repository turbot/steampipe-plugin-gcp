---
title: "Steampipe Table: gcp_compute_resource_policy - Query GCP Compute Engine Resource Policies using SQL"
description: "Allows users to query Resource Policies in Google Cloud's Compute Engine, providing insights into the rules that schedule disk snapshot and VM operations."
folder: "Compute"
---

# Table: gcp_compute_resource_policy - Query GCP Compute Engine Resource Policies using SQL

Resource Policies in Google Cloud's Compute Engine are used to schedule operations for your instances. For example, they can be used to schedule periodic snapshot creation for persistent disk, VM start and stop schedules, and more. These policies help to automate routine tasks, which can increase operational efficiency and reduce the potential for error.

## Table Usage Guide

The `gcp_compute_resource_policy` table offers insights into Resource Policies within Google Cloud's Compute Engine. As a cloud engineer, you can leverage this table to explore policy-specific details, including the scheduled operations, the frequency of these operations, and the instances to which they apply. Use this table to understand your resource scheduling policies, verify their configurations, and ensure they are operating as intended.

## Examples

### Basic info
Explore which GCP compute resource policies are currently active by assessing their status, providing a quick way to monitor and manage your resources effectively.

```sql+postgres
select
  name,
  status,
  self_link
from
  gcp_compute_resource_policy;
```

```sql+sqlite
select
  name,
  status,
  self_link
from
  gcp_compute_resource_policy;
```

### List policies used to schedule an instance
Explore which policies are used to schedule instances in your GCP Compute Engine. This can help you understand and manage your resource allocation more effectively.

```sql+postgres
select
  p.name as policy_name,
  i.name,
  p.instance_schedule_policy
from
  gcp_compute_resource_policy as p
  join gcp_compute_instance as i on i.resource_policies ?| array[p.self_link]
where
  p.instance_schedule_policy is not null;
```

```sql+sqlite
select
  p.name as policy_name,
  i.name,
  p.instance_schedule_policy
from
  gcp_compute_resource_policy as p
  join gcp_compute_instance as i on json_extract(i.resource_policies, p.self_link) is not null
where
  p.instance_schedule_policy is not null;
```

### List invalid policies
Explore which policies in your Google Cloud Platform compute resources are invalid. This can be beneficial for maintaining optimal resource management and troubleshooting potential issues.

```sql+postgres
select
  name,
  self_link,
  status
from
  gcp_compute_resource_policy
where
  status = 'INVALID';
```

```sql+sqlite
select
  name,
  self_link,
  status
from
  gcp_compute_resource_policy
where
  status = 'INVALID';
```