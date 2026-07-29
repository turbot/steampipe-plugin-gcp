---
title: "Steampipe Table: gcp_cloud_asset_org_policy - Query GCP Cloud Asset Organization Policies using SQL"
description: "Allows users to query the organization policies set directly on assets in a GCP project, providing visibility into where policy constraints are applied."
folder: "Cloud Asset"
---

# Table: gcp_cloud_asset_org_policy - Query GCP Cloud Asset Organization Policies using SQL

GCP Organization policies let administrators set constraints on how resources can be used. GCP Cloud Asset Inventory records the organization policies set directly on each asset, providing visibility into where constraints are applied across the resource hierarchy.

## Table Usage Guide

The `gcp_cloud_asset_org_policy` table returns one row per (asset, policy) pair, so an asset with multiple constraints configured appears once per constraint. Assets without a directly attached organization policy are not returned; policies inherited from ancestors are not included.

## Examples

### Basic info

Get every organization policy set directly on an asset in the project.

```sql+postgres
select
  name,
  asset_type,
  constraint_name,
  policy_update_time
from
  gcp_cloud_asset_org_policy;
```

```sql+sqlite
select
  name,
  asset_type,
  constraint_name,
  policy_update_time
from
  gcp_cloud_asset_org_policy;
```

### List assets with a specific constraint configured

Find where a particular constraint is applied, for example domain-restricted sharing.

```sql+postgres
select
  name,
  asset_type,
  list_policy,
  boolean_policy
from
  gcp_cloud_asset_org_policy
where
  constraint_name = 'constraints/iam.allowedPolicyMemberDomains';
```

```sql+sqlite
select
  name,
  asset_type,
  list_policy,
  boolean_policy
from
  gcp_cloud_asset_org_policy
where
  constraint_name = 'constraints/iam.allowedPolicyMemberDomains';
```

### List enforced boolean constraints

Review which boolean constraints are actively enforced on assets.

```sql+postgres
select
  name,
  constraint_name
from
  gcp_cloud_asset_org_policy
where
  (boolean_policy ->> 'enforced')::bool;
```

```sql+sqlite
select
  name,
  constraint_name
from
  gcp_cloud_asset_org_policy
where
  json_extract(boolean_policy, '$.enforced');
```
