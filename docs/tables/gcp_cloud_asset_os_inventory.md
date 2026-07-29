---
title: "Steampipe Table: gcp_cloud_asset_os_inventory - Query GCP Cloud Asset OS Inventories using SQL"
description: "Allows users to query the runtime OS inventory of compute instances in a GCP project, including OS details and installed packages, for patch and vulnerability management."
folder: "Cloud Asset"
---

# Table: gcp_cloud_asset_os_inventory - Query GCP Cloud Asset OS Inventories using SQL

GCP Cloud Asset Inventory records runtime OS inventory information reported by the OS Config agent on compute instances, including operating system details and installed packages. This is useful for patch and vulnerability management across a fleet.

## Table Usage Guide

The `gcp_cloud_asset_os_inventory` table returns one row per compute instance with OS inventory data. Only instances running the OS Config agent with [OS inventory management enabled](https://cloud.google.com/compute/docs/instances/os-inventory-management) are returned.

## Examples

### Basic info

Get the OS details of every instance reporting inventory.

```sql+postgres
select
  name,
  hostname,
  os_long_name,
  architecture,
  kernel_version,
  inventory_update_time
from
  gcp_cloud_asset_os_inventory;
```

```sql+sqlite
select
  name,
  hostname,
  os_long_name,
  architecture,
  kernel_version,
  inventory_update_time
from
  gcp_cloud_asset_os_inventory;
```

### Count instances by operating system

Review the OS distribution across the fleet.

```sql+postgres
select
  os_short_name,
  os_version,
  count(*)
from
  gcp_cloud_asset_os_inventory
group by
  os_short_name,
  os_version
order by
  count desc;
```

```sql+sqlite
select
  os_short_name,
  os_version,
  count(*)
from
  gcp_cloud_asset_os_inventory
group by
  os_short_name,
  os_version
order by
  count(*) desc;
```

### List instances with stale inventory

Find instances that have not reported inventory recently, which may indicate a stopped instance or a broken OS Config agent.

```sql+postgres
select
  name,
  hostname,
  inventory_update_time
from
  gcp_cloud_asset_os_inventory
where
  inventory_update_time < now() - interval '1 day';
```

```sql+sqlite
select
  name,
  hostname,
  inventory_update_time
from
  gcp_cloud_asset_os_inventory
where
  inventory_update_time < datetime('now', '-1 day');
```
