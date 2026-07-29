---
title: "Steampipe Table: gcp_cloud_asset_resource - Query GCP Cloud Asset Resources using SQL"
description: "Allows users to query the full resource representation of every asset in a GCP project, providing a detailed configuration inventory across all GCP services."
folder: "Cloud Asset"
---

# Table: gcp_cloud_asset_resource - Query GCP Cloud Asset Resources using SQL

GCP Cloud Asset Inventory keeps a full representation of each resource, including its service-specific configuration data. This provides a single place to inspect resource configurations across all GCP services.

## Table Usage Guide

The `gcp_cloud_asset_resource` table returns one row per asset with its full resource representation. The `data` column contains the service-specific resource body, whose shape varies by `asset_type`. For a lightweight inventory without resource data, use the [gcp_cloud_asset](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_cloud_asset) table instead.

## Examples

### Basic info

Get the resource representation details of every asset in the project.

```sql+postgres
select
  name,
  asset_type,
  location,
  parent,
  update_time
from
  gcp_cloud_asset_resource;
```

```sql+sqlite
select
  name,
  asset_type,
  location,
  parent,
  update_time
from
  gcp_cloud_asset_resource;
```

### Inspect the full configuration of a specific asset type

Review the service-specific configuration body of each resource, for example BigQuery datasets.

```sql+postgres
select
  name,
  location,
  data
from
  gcp_cloud_asset_resource
where
  asset_type = 'bigquery.googleapis.com/Dataset';
```

```sql+sqlite
select
  name,
  location,
  data
from
  gcp_cloud_asset_resource
where
  asset_type = 'bigquery.googleapis.com/Dataset';
```

### Count resources by location

Identify where resources are deployed across regions and zones.

```sql+postgres
select
  location,
  count(*)
from
  gcp_cloud_asset_resource
group by
  location
order by
  count desc;
```

```sql+sqlite
select
  location,
  count(*)
from
  gcp_cloud_asset_resource
group by
  location
order by
  count(*) desc;
```
