---
title: "Steampipe Table: gcp_dataplex_asset - Query GCP Dataplex Assets using SQL"
description: "Allows users to query GCP Dataplex assets, providing detailed information on asset configurations, discovery status, and associated resources."
folder: "Dataplex"
---

# Table: gcp_dataplex_asset - Query GCP Dataplex Assets using SQL

Google Cloud Dataplex is an intelligent data fabric that helps you manage, monitor, and govern your data across various cloud and on-premises environments. The `gcp_dataplex_asset` table in Steampipe allows you to query information about Dataplex assets, including their resource specifications, discovery status, security status, and associated lakes and zones.

## Table Usage Guide

The `gcp_dataplex_asset` table is useful for cloud administrators, data engineers, and security professionals to gather detailed insights into their Dataplex assets. You can query various aspects of the assets, such as their resource specifications, discovery and security statuses, and their associations with specific lakes and zones. This table is particularly helpful for monitoring asset states, managing data access, and ensuring that your data is properly organized and governed.

**Important Notes**
- You must specify the `zone_name` in the `where` clause (`where zone_name='projects/{projectName}/locations/us-central1/lakes/{lakeId}/zones/{zoneId}'`) to list the assets.

## Examples

### Basic asset information
Retrieve basic information about Dataplex assets, including their name, display name, and creation time.

```sql+postgres
select
  name,
  display_name,
  create_time,
  state,
  location,
  project
from
  gcp_dataplex_asset
where
  zone_name = 'projects/parker-aaa/locations/us-central1/lakes/dasdadsa/zones/tese9392';
```

```sql+sqlite
select
  name,
  display_name,
  create_time,
  state,
  location,
  project
from
  gcp_dataplex_asset
where
  zone_name = 'projects/parker-aaa/locations/us-central1/lakes/dasdadsa/zones/tese9392';
```

### List assets by state and resource type
Identify assets that are in a specific state (e.g., READY) and of a specific resource type (e.g., BIGQUERY_DATASET).

```sql+postgres
select
  name,
  display_name,
  state,
  resource_spec ->> 'type' as resource_type
from
  gcp_dataplex_asset
where
  state = 'ACTIVE'
  and resource_spec ->> 'type' = 'BIGQUERY_DATASET'
  and zone_name = 'projects/parker-aaa/locations/us-central1/lakes/dasdadsa/zones/tese9392';
```

```sql+sqlite
select
  name,
  display_name,
  state,
  json_extract(resource_spec, '$.type') as resource_type
from
  gcp_dataplex_asset
where
  state = 'ACTIVE'
  and json_extract(resource_spec, '$.type') = 'BIGQUERY_DATASET'
  and zone_name = 'projects/parker-aaa/locations/us-central1/lakes/dasdadsa/zones/tese9392';
```

### List assets with discovery details
Retrieve assets that have discovery features enabled, showing details about the last discovery run and its duration.

```sql+postgres
select
  name,
  display_name,
  discovery_status ->> 'state' as discovery_state,
  discovery_status ->> 'lastRunTime' as last_run_time,
  discovery_status ->> 'lastRunDuration' as last_run_duration
from
  gcp_dataplex_asset
where
  discovery_status is not null
  and zone_name = 'projects/parker-aaa/locations/us-central1/lakes/dasdadsa/zones/tese9392';
```

```sql+sqlite
select
  name,
  display_name,
  json_extract(discovery_status, '$.state') as discovery_state,
  json_extract(discovery_status, '$.lastRunTime') as last_run_time,
  json_extract(discovery_status, '$.lastRunDuration') as last_run_duration
from
  gcp_dataplex_asset
where
  discovery_status is not null
  and zone_name = 'projects/parker-aaa/locations/us-central1/lakes/dasdadsa/zones/tese9392';
```

### List assets and join with zones and lakes
Retrieve a list of assets along with their associated zones and lakes by joining with the `gcp_dataplex_zone` and `gcp_dataplex_lake` tables.

```sql+postgres
select
  a.name as asset_name,
  a.display_name as asset_display_name,
  z.name as zone_name,
  l.name as lake_name
from
  gcp_dataplex_asset a
join
  gcp_dataplex_zone z on a.zone_name = z.name
join
  gcp_dataplex_lake l on z.lake_name = l.name;
```

```sql+sqlite
select
  a.name as asset_name,
  a.display_name as asset_display_name,
  z.name as zone_name,
  l.name as lake_name
from
  gcp_dataplex_asset a
join
  gcp_dataplex_zone z on a.zone_name = z.name
join
  gcp_dataplex_lake l on z.lake_name = l.name;
```