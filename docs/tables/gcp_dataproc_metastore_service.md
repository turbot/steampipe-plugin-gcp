---
title: "Steampipe Table: gcp_dataproc_metastore_service - Query GCP Dataproc Metastore Services using SQL"
description: "Allows users to query GCP Dataproc Metastore services, providing detailed information on service configurations, networking, and scaling settings."
folder: "Dataproc"
---

# Table: gcp_dataproc_metastore_service - Query GCP Dataproc Metastore Services using SQL

Google Cloud Dataproc Metastore is a fully managed Apache Hive metastore service that helps manage metadata for various data lakes and other data processing tools. The `gcp_dataproc_metastore_service` table in Steampipe allows you to query detailed information about Dataproc Metastore services in your GCP environment, including network settings, database configurations, state, and scaling options.

## Table Usage Guide

The `gcp_dataproc_metastore_service` table enables cloud administrators, data engineers, and DevOps teams to gather insights into their Dataproc Metastore services. You can query various aspects of the service, such as its network configurations, encryption settings, scaling policies, and metadata management activities. This table is particularly useful for managing Metastore services, monitoring their health, and ensuring that the service is configured correctly for specific use cases.

## Examples

### Basic info
Retrieve basic information about Dataproc Metastore services, including their name, location, and state.

```sql+postgres
select
  name,
  location,
  state,
  create_time,
  release_channel
from
  gcp_dataproc_metastore_service;
```

```sql+sqlite
select
  name,
  location,
  state,
  create_time,
  release_channel
from
  gcp_dataproc_metastore_service;
```

### List services with specific database types
Identify Dataproc Metastore services based on the type of database they store, such as "MYSQL".

```sql+postgres
select
  name,
  database_type,
  location,
  project
from
  gcp_dataproc_metastore_service
where
  database_type = 'MYSQL';
```

```sql+sqlite
select
  name,
  database_type,
  location,
  project
from
  gcp_dataproc_metastore_service
where
  database_type = 'MYSQL';
```

### List services with scaling configurations
Retrieve services with defined scaling configurations, which can be useful for monitoring resource scaling and ensuring proper performance.

```sql+postgres
select
  name,
  scaling_config,
  project,
  location
from
  gcp_dataproc_metastore_service
where
  scaling_config is not null;
```

```sql+sqlite
select
  name,
  scaling_config,
  project,
  location
from
  gcp_dataproc_metastore_service
where
  scaling_config is not null;
```

### List services with network configurations
Fetch Dataproc Metastore services based on their network settings, such as VPC networks.

```sql+postgres
select
  name,
  network,
  network_config,
  location
from
  gcp_dataproc_metastore_service
where
  network is not null;
```

```sql+sqlite
select
  name,
  network,
  network_config,
  location
from
  gcp_dataproc_metastore_service
where
  network is not null;
```

### List services by state
Identify Metastore services based on their current state, such as "ACTIVE" or "FAILED."

```sql+postgres
select
  name,
  state,
  state_message,
  location
from
  gcp_dataproc_metastore_service
where
  state = 'ACTIVE';
```

```sql+sqlite
select
  name,
  state,
  state_message,
  location
from
  gcp_dataproc_metastore_service
where
  state = 'ACTIVE';
```

### Services with their network configurations
Retrieve Metastore services along with their associated network configurations.

```sql+postgres
select
  m.name as service_name,
  m.location,
  m.network,
  m.network_config,
  n.name as network_name,
  n.auto_create_subnetworks,
  n.peerings
from
  gcp_dataproc_metastore_service m
join
  gcp_compute_network n on n.name = split_part(m.network, '/', 5);
```

```sql+sqlite
select
  m.name as service_name,
  m.location,
  m.network,
  m.network_config,
  n.name as network_name,
  n.auto_create_subnetworks,
  n.peerings
from
  gcp_dataproc_metastore_service m
join
  gcp_compute_network n on n.name = json_extract(m.network, '$[5]');
```