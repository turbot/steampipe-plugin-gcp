---
title: "Steampipe Table: gcp_cloud_asset_access_policy - Query GCP Cloud Asset Access Policies using SQL"
description: "Allows users to query Access Context Manager policies, access levels, and service perimeters recorded in GCP Cloud Asset Inventory."
folder: "Cloud Asset"
---

# Table: gcp_cloud_asset_access_policy - Query GCP Cloud Asset Access Policies using SQL

GCP Access Context Manager lets organizations define fine-grained, attribute-based access control for projects and resources, including VPC Service Controls perimeters. GCP Cloud Asset Inventory records these policies as assets.

## Table Usage Guide

The `gcp_cloud_asset_access_policy` table returns one row per Access Context Manager asset. Each row carries exactly one of `access_policy`, `access_level`, or `service_perimeter`, depending on the `asset_type`; the other two columns are null.

## Examples

### Basic info

Get every Access Context Manager asset in scope.

```sql+postgres
select
  name,
  asset_type,
  access_policy,
  access_level,
  service_perimeter
from
  gcp_cloud_asset_access_policy;
```

```sql+sqlite
select
  name,
  asset_type,
  access_policy,
  access_level,
  service_perimeter
from
  gcp_cloud_asset_access_policy;
```

### List service perimeters and their protected resources

Review VPC Service Controls perimeters and which resources they protect.

```sql+postgres
select
  name,
  service_perimeter ->> 'title' as title,
  service_perimeter ->> 'perimeterType' as perimeter_type,
  service_perimeter -> 'status' -> 'resources' as protected_resources
from
  gcp_cloud_asset_access_policy
where
  service_perimeter is not null;
```

```sql+sqlite
select
  name,
  json_extract(service_perimeter, '$.title') as title,
  json_extract(service_perimeter, '$.perimeterType') as perimeter_type,
  json_extract(service_perimeter, '$.status.resources') as protected_resources
from
  gcp_cloud_asset_access_policy
where
  service_perimeter is not null;
```

### List access levels

Review the conditions under which requests are permitted.

```sql+postgres
select
  name,
  access_level ->> 'title' as title,
  access_level -> 'basic' as basic_level,
  access_level -> 'custom' as custom_level
from
  gcp_cloud_asset_access_policy
where
  access_level is not null;
```

```sql+sqlite
select
  name,
  json_extract(access_level, '$.title') as title,
  json_extract(access_level, '$.basic') as basic_level,
  json_extract(access_level, '$.custom') as custom_level
from
  gcp_cloud_asset_access_policy
where
  access_level is not null;
```
