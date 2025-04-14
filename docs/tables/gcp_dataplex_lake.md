---
title: "Steampipe Table: gcp_dataplex_lake - Query GCP Dataplex Lakes using SQL"
description: "Allows users to query GCP Dataplex Lakes, providing detailed information about each lake's configuration, status, and metadata."
folder: "Dataplex"
---

# Table: gcp_dataplex_lake - Query GCP Dataplex Lakes using SQL

GCP Dataplex Lakes are managed data lakes that provide unified analytics and governance for data at scale. Dataplex simplifies data management by automating discovery, organization, and management of data across various storage systems.

## Table Usage Guide

The `gcp_dataplex_lake` table allows data engineers and cloud administrators to query and manage Dataplex Lakes within their GCP environment. You can retrieve information about the lake's configuration, status, associated metastore, and more. This table is useful for monitoring and managing the state and metadata of Dataplex Lakes.

## Examples

### Basic info
Retrieve a list of all Dataplex Lakes in your GCP account to get an overview of your managed data lakes.

```sql+postgres
select
  display_name,
  name,
  state,
  create_time,
  service_account
from
  gcp_dataplex_lake;
```

```sql+sqlite
select
  display_name,
  name,
  state,
  create_time,
  service_account
from
  gcp_dataplex_lake;
```

### Dataplex Lakes by location
Explore which regions have the most Dataplex Lakes to understand your data infrastructure distribution better.

```sql+postgres
select
  location,
  count(*)
from
  gcp_dataplex_lake
group by
  location;
```

```sql+sqlite
select
  location,
  count(*)
from
  gcp_dataplex_lake
group by
  location;
```

### Get details of lakes with a specific state
Retrieve Dataplex Lakes in a specific state (e.g., `ACTIVE`) to monitor their status.

```sql+postgres
select
  name,
  state,
  create_time,
  update_time
from
  gcp_dataplex_lake
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
  gcp_dataplex_lake
where
  state = 'ACTIVE';
```

### Get Dataplex Lakes with the associated metastore settings
List all Dataplex Lakes that have an associated Dataproc Metastore, including their metastore settings and status.

```sql+postgres
select
  name,
  metastore ->> 'service' as metastore_service,
  metastore_status ->> 'state' as metastore_state,
  location
from
  gcp_dataplex_lake
where
  metastore is not null;
```

```sql+sqlite
select
  name,
  json_extract(metastore, '$.service') as metastore_service,
  json_extract(metastore_status, '$.state') as metastore_state,
  location
from
  gcp_dataplex_lake
where
  metastore is not null;
```
