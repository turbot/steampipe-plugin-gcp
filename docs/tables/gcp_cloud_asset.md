---
title: "Steampipe Table: gcp_cloud_asset - Query GCP Cloud Asset using SQL"
description: "Allows users to query GCP Cloud Asset, specifically visibility and control over their cloud resources, ensuring that they can manage these assets effectively in terms of security, compliance, and operational efficiency."
folder: "Cloud Asset"
---

# Table: gcp_cloud_asset - Query GCP Cloud Asset using SQL

GCP Cloud Asset is a powerful tool for organizations to maintain visibility and control over their cloud resources, ensuring that they can manage these assets effectively in terms of security, compliance, and operational efficiency.

## Table Usage Guide

The `gcp_cloud_asset` table provides an management system for resources and policies within GCP. It allows users to keep track of their cloud assets across various GCP services.

## Examples

### Basic info

It provides a quick snapshot of all assets in the GCP environment. This is helpful for administrators and cloud architects to get an overview of the resources, their types, and recent updates.

```sql+postgres
select
  name,
  asset_type,
  update_time,
  ancestors
from
  gcp_cloud_asset;
```

```sql+sqlite
select
  name,
  asset_type,
  update_time,
  ancestors
from
  gcp_cloud_asset;
```

### Get access policy of the resources

This query is particularly useful for administrators and security professionals who need to oversee and manage access policies within a GCP environment. It provides a detailed view of how access is controlled and managed across various cloud assets.

```sql+postgres
select
  name,
  access_policy ->> 'Etag' as access_policy_etag,
  access_policy ->> 'Name' as access_policy_name,
  access_policy ->> 'Parent' as access_policy_parent,
  access_policy -> 'Scopes' as access_policy_scopes
from
  gcp_cloud_asset;
```

```sql+sqlite
select
  name,
  json_extract(access_policy, '$.Etag') as access_policy_etag,
  json_extract(access_policy, '$.Name') as access_policy_name,
  json_extract(access_policy, '$.Parent') as access_policy_parent,
  json_extract(access_policy, '$.Scopes') as access_policy_scopes
from
  gcp_cloud_asset;
```
