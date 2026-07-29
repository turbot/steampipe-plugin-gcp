---
title: "Steampipe Table: gcp_cloud_asset - Query GCP Cloud Assets using SQL"
description: "Allows users to query GCP Cloud Assets, providing a basic inventory of all resources and policies within a GCP project."
folder: "Cloud Asset"
---

# Table: gcp_cloud_asset - Query GCP Cloud Assets using SQL

GCP Cloud Asset Inventory is a powerful tool for organizations to maintain visibility and control over their cloud resources, ensuring that they can manage these assets effectively in terms of security, compliance, and operational efficiency.

## Table Usage Guide

The `gcp_cloud_asset` table provides a basic inventory of all assets within a GCP project: the asset name, type, ancestry path, and last update time. It allows users to keep track of their cloud assets across various GCP services.

To query content attached to assets, use the content-specific tables instead:

| Table | Content |
| ----- | ------- |
| [gcp_cloud_asset_resource](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_cloud_asset_resource) | Full resource representations |
| [gcp_cloud_asset_iam_policy](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_cloud_asset_iam_policy) | Resource-level IAM policies |
| [gcp_cloud_asset_org_policy](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_cloud_asset_org_policy) | Organization policies set on assets |
| [gcp_cloud_asset_os_inventory](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_cloud_asset_os_inventory) | Runtime OS inventory of compute instances |
| [gcp_cloud_asset_access_policy](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_cloud_asset_access_policy) | Access Context Manager policies, levels, and service perimeters |
| [gcp_cloud_asset_relationship](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_cloud_asset_relationship) | Asset relationships (SCC Premium/Enterprise only) |

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

### Count assets by type

Identify which services make up the bulk of your cloud footprint.

```sql+postgres
select
  asset_type,
  count(*)
from
  gcp_cloud_asset
group by
  asset_type
order by
  count desc;
```

```sql+sqlite
select
  asset_type,
  count(*)
from
  gcp_cloud_asset
group by
  asset_type
order by
  count(*) desc;
```

### List assets updated in the last 7 days

Track recent changes in the environment.

```sql+postgres
select
  name,
  asset_type,
  update_time
from
  gcp_cloud_asset
where
  update_time > now() - interval '7 days';
```

```sql+sqlite
select
  name,
  asset_type,
  update_time
from
  gcp_cloud_asset
where
  update_time > datetime('now', '-7 days');
```
