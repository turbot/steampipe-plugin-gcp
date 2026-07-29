---
title: "Steampipe Table: gcp_cloud_asset_iam_policy - Query GCP Cloud Asset IAM Policies using SQL"
description: "Allows users to query the resource-level IAM policies of every asset in a GCP project, providing a project-wide view of role bindings for security auditing."
folder: "Cloud Asset"
---

# Table: gcp_cloud_asset_iam_policy - Query GCP Cloud Asset IAM Policies using SQL

GCP Cloud Asset Inventory records the IAM policy set directly on each resource. This provides a single place to audit role bindings across all resources in a project, regardless of which service they belong to.

## Table Usage Guide

The `gcp_cloud_asset_iam_policy` table returns one row per asset that has a resource-level IAM policy set on it. Assets without a directly attached IAM policy are not returned; policies inherited from ancestors are not included.

## Examples

### Basic info

Get the IAM policy details of every asset that has one.

```sql+postgres
select
  name,
  asset_type,
  bindings,
  update_time
from
  gcp_cloud_asset_iam_policy;
```

```sql+sqlite
select
  name,
  asset_type,
  bindings,
  update_time
from
  gcp_cloud_asset_iam_policy;
```

### List role bindings per member

Fan out the bindings to see which members hold which roles on which resources. This is useful for auditing who has access to what across the entire project.

```sql+postgres
select
  name,
  asset_type,
  b ->> 'role' as role,
  jsonb_array_elements_text(b -> 'members') as member
from
  gcp_cloud_asset_iam_policy,
  jsonb_array_elements(bindings) as b;
```

```sql+sqlite
select
  a.name,
  a.asset_type,
  json_extract(b.value, '$.role') as role,
  m.value as member
from
  gcp_cloud_asset_iam_policy a,
  json_each(a.bindings) b,
  json_each(json_extract(b.value, '$.members')) m;
```

### Find resources granting access to allUsers or allAuthenticatedUsers

Detect resources that are publicly accessible through their resource-level IAM policy.

```sql+postgres
select
  name,
  asset_type,
  b ->> 'role' as role
from
  gcp_cloud_asset_iam_policy,
  jsonb_array_elements(bindings) as b
where
  b -> 'members' ?| array['allUsers', 'allAuthenticatedUsers'];
```

```sql+sqlite
select
  a.name,
  a.asset_type,
  json_extract(b.value, '$.role') as role
from
  gcp_cloud_asset_iam_policy a,
  json_each(a.bindings) b,
  json_each(json_extract(b.value, '$.members')) m
where
  m.value in ('allUsers', 'allAuthenticatedUsers');
```
