---
title: "Steampipe Table: gcp_dataplex_zone - Query GCP Dataplex Zones using SQL"
description: "Allows users to query GCP Dataplex Zones, providing detailed information about each zone's configuration, status, and associated lake."
folder: "Dataplex"
---

# Table: gcp_dataplex_zone - Query GCP Dataplex Zones using SQL

GCP Dataplex Zones are logical groupings within a Dataplex Lake, designed to organize and manage data based on various criteria, such as data type, security requirements, or lifecycle management. Zones allow for granular control over data management and access within a Lake.

## Table Usage Guide

The `gcp_dataplex_zone` table allows data engineers and cloud administrators to query and manage Dataplex Zones within their GCP environment. You can retrieve information about the zone's configuration, status, associated lake, and more. This table is useful for monitoring and managing the state and metadata of Dataplex Zones.

## Examples

### Basic info
Retrieve a list of all Dataplex Zones in your GCP account to get an overview of your managed data zones.

```sql+postgres
select
  name,
  display_name,
  lake_name,
  create_time,
  type
from
  gcp_dataplex_zone;
```

```sql+sqlite
select
  name,
  display_name,
  lake_name,
  create_time,
  type
from
  gcp_dataplex_zone;
```

### Dataplex zones by type
Explore the different types of Dataplex Zones to understand how your data is organized and managed within lakes.

```sql+postgres
select
  type,
  count(*)
from
  gcp_dataplex_zone
group by
  type;
```

```sql+sqlite
select
  type,
  count(*)
from
  gcp_dataplex_zone
group by
  type;
```

### Get details of zones in a specific state
Retrieve Dataplex Zones in a specific state (e.g., `ACTIVE`) to monitor their status.

```sql+postgres
select
  name,
  state,
  create_time,
  update_time
from
  gcp_dataplex_zone
where
  state = 'ACTIVE';
```

```sql+sqlite
select
  name,
  state,
  create_time,
  update_time
from
  gcp_dataplex_zone
where
  state = 'ACTIVE';
```

### Dataplex zones with their associated lakes
This is useful for understanding the relationship between zones and lakes in your Dataplex environment.

```sql+postgres
select
  z.name as zone_name,
  z.type as zone_type,
  z.state as zone_state,
  z.create_time as zone_create_time,
  l.name as lake_name,
  l.location as lake_location,
  l.state as lake_state
from
  gcp_dataplex_zone as z
join
  gcp_dataplex_lake as l
on
  z.lake_name = l.name;
```

```sql+sqlite
select
  z.name as zone_name,
  z.type as zone_type,
  z.state as zone_state,
  z.create_time as zone_create_time,
  l.name as lake_name,
  l.location as lake_location,
  l.state as lake_state
from
  gcp_dataplex_zone as z
join
  gcp_dataplex_lake as l
on
  z.lake_name = l.name;
```
