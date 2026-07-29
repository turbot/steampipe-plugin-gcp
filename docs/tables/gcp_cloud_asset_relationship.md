---
title: "Steampipe Table: gcp_cloud_asset_relationship - Query GCP Cloud Asset Relationships using SQL"
description: "Allows users to query relationships between assets in a GCP project, such as which instances belong to which instance groups."
folder: "Cloud Asset"
---

# Table: gcp_cloud_asset_relationship - Query GCP Cloud Asset Relationships using SQL

GCP Cloud Asset Inventory records relationships between assets, such as which compute instances belong to which instance groups. See the [supported relationship types](https://cloud.google.com/asset-inventory/docs/supported-asset-types#supported_relationship_types) for the full list.

## Table Usage Guide

The `gcp_cloud_asset_relationship` table returns one row per (asset, related asset) pair.

**Important:** Asset relationships are only available to Security Command Center Premium and Enterprise tier customers. For other customers, the Cloud Asset API returns an authorization error and queries against this table fail.

## Examples

### Basic info

Get every recorded asset relationship in the project.

```sql+postgres
select
  name,
  asset_type,
  relationship_type,
  related_asset_name,
  related_asset_type
from
  gcp_cloud_asset_relationship;
```

```sql+sqlite
select
  name,
  asset_type,
  relationship_type,
  related_asset_name,
  related_asset_type
from
  gcp_cloud_asset_relationship;
```

### List instances and the instance groups they belong to

Map group membership of compute instances.

```sql+postgres
select
  name as instance,
  related_asset_name as instance_group
from
  gcp_cloud_asset_relationship
where
  relationship_type = 'INSTANCE_TO_INSTANCEGROUP';
```

```sql+sqlite
select
  name as instance,
  related_asset_name as instance_group
from
  gcp_cloud_asset_relationship
where
  relationship_type = 'INSTANCE_TO_INSTANCEGROUP';
```

### Count relationships by type

Review which relationship types exist in the project.

```sql+postgres
select
  relationship_type,
  count(*)
from
  gcp_cloud_asset_relationship
group by
  relationship_type
order by
  count desc;
```

```sql+sqlite
select
  relationship_type,
  count(*)
from
  gcp_cloud_asset_relationship
group by
  relationship_type
order by
  count(*) desc;
```
