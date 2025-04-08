---
title: "Steampipe Table: gcp_alloydb_cluster - Query Google Cloud Platform AlloyDB Clusters using SQL"
description: "Allows users to query Google Cloud Platform AlloyDB Clusters, providing insights into cluster configurations, status, and associated metadata."
folder: "AlloyDB"
---

# Table: gcp_alloydb_cluster - Query Google Cloud Platform AlloyDB Clusters using SQL

Google Cloud AlloyDB is a fully managed, PostgreSQL-compatible database service optimized for performance and scalability. Built on Google's infrastructure, AlloyDB provides high availability, security, and integration with other services, making it ideal for enterprise database management solutions.

## Table Usage Guide

The `gcp_alloydb_cluster` table enables you to query information about AlloyDB clusters within Google Cloud Platform. It is useful for database administrators and developers to monitor the operational aspects of their AlloyDB environments, including configuration details, current state, and resource utilization.

## Examples

### Basic info
Retrieve basic information about your Google Cloud Platform's AlloyDB clusters. Useful for a quick overview and operational monitoring of cluster attributes.

```sql+postgres
select
  name,
  state,
  display_name,
  database_version,
  location
from
  gcp_alloydb_cluster;
```

```sql+sqlite
select
  name,
  state,
  display_name,
  database_version,
  location
from
  gcp_alloydb_cluster;
```

### List clusters by state
Identify AlloyDB clusters in a specific state to manage and troubleshoot operational needs effectively.

```sql+postgres
select
  name,
  state
from
  gcp_alloydb_cluster
where
  state = 'MAINTENANCE';
```

```sql+sqlite
select
  name,
  state
from
  gcp_alloydb_cluster
where
  state = 'MAINTENANCE';
```

### Detailed configuration of a cluster
Access detailed configuration settings of an AlloyDB cluster to understand its setup and make informed decisions about scaling and modifications.

```sql+postgres
select
  name,
  encryption_config,
  network_config
from
  gcp_alloydb_cluster
where
  display_name = 'your-cluster-name';
```

```sql+sqlite
select
  name,
  json_extract(encryption_config, '$') as encryption_config,
  json_detail(network_config, '$') as network_config
from
  gcp_alloydb_cluster
where
  display_name = 'your-cluster-name';
```

### List Clusters with Specific Database Version
This query retrieves all AlloyDB clusters that are using a specific database version, which is helpful for auditing purposes or planning upgrades.

```sql+postgres
select
  name,
  database_version,
  state,
  location
from
  gcp_alloydb_cluster
where
  database_version = 'POSTGRES_14';
```

```sql+sqlite
select
  name,
  database_version,
  state,
  location
from
  gcp_alloydb_cluster
where
  database_version = 'POSTGRES_14';
```

### Find Clusters with Encryption Enabled
This query is useful to ensure compliance by checking which clusters have encryption configured.

```sql+postgres
select
  name,
  encryption_config
from
  gcp_alloydb_cluster
where
  encryption_config is not null;
```

```sql+sqlite
select
  name,
  json_extract(encryption_config, '$')
from
  gcp_alloydb_cluster
where
  encryption_config is not null;
```