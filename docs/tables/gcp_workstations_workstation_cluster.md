---
title: "Steampipe Table: gcp_workstations_workstation_cluster - Query Google Cloud Workstations Clusters using SQL"
description: "Query Google Cloud Workstations Clusters, including configuration, state, and network details for cloud-based development environments."
folder: "Workstations"
---

# Table: gcp_workstations_workstation_cluster - Query Google Cloud Workstations Clusters using SQL

Google Cloud Workstations Clusters are logical groupings of workstation configurations and workstations. A cluster defines the network and subnetwork where workstations will be created, and provides a way to organize and manage multiple workstation configurations and their associated workstations.

## Table Usage Guide

The `gcp_workstations_workstation_cluster` table provides insights into Workstations Clusters within Google Cloud. As a developer or platform engineer, explore cluster-specific details through this table, including network configurations, states, and metadata. Use it to uncover information about clusters, such as their current state, associated network settings, and the resources allocated to each cluster.

## Examples

### Basic info
Explore the basic details of your Google Cloud Workstations Clusters, including their names and creation times. This information helps you understand the configuration and status of your clusters, which is useful for managing cloud development environments.

```sql+postgres
select
  name,
  display_name,
  create_time,
  update_time,
  uid,
  location,
  project
from
  gcp_workstations_workstation_cluster;
```

```sql+sqlite
select
  name,
  display_name,
  create_time,
  update_time,
  uid,
  location,
  project
from
  gcp_workstations_workstation_cluster;
```

### List clusters with network configuration
Identify clusters along with their network and subnetwork configuration to understand how workstations are connected.

```sql+postgres
select
  name,
  display_name,
  network,
  subnetwork,
  location
from
  gcp_workstations_workstation_cluster;
```

```sql+sqlite
select
  name,
  display_name,
  network,
  subnetwork,
  location
from
  gcp_workstations_workstation_cluster;
```

### List clusters that are being updated
Analyze which clusters are currently reconciling (being updated). This is useful for tracking environments undergoing changes.

```sql+postgres
select
  name,
  display_name,
  reconciling,
  update_time,
  location
from
  gcp_workstations_workstation_cluster
where
  reconciling = true;
```

```sql+sqlite
select
  name,
  display_name,
  reconciling,
  update_time,
  location
from
  gcp_workstations_workstation_cluster
where
  reconciling = 1;
```

### List clusters by location
Discover which clusters exist in a specific location to understand the geographical distribution of your development environments.

```sql+postgres
select
  name,
  display_name,
  network,
  subnetwork,
  location
from
  gcp_workstations_workstation_cluster
where
  location = 'us-central1';
```

```sql+sqlite
select
  name,
  display_name,
  network,
  subnetwork,
  location
from
  gcp_workstations_workstation_cluster
where
  location = 'us-central1';
```

### List clusters with private cluster configuration
Explore clusters that have private cluster configuration enabled to understand which clusters are using private networking.

```sql+postgres
select
  name,
  display_name,
  private_cluster_config,
  network,
  subnetwork,
  location
from
  gcp_workstations_workstation_cluster
where
  private_cluster_config is not null;
```

```sql+sqlite
select
  name,
  display_name,
  private_cluster_config,
  network,
  subnetwork,
  location
from
  gcp_workstations_workstation_cluster
where
  private_cluster_config is not null;
```

### List clusters with conditions
Analyze clusters that have status conditions to understand any issues or warnings associated with the clusters.

```sql+postgres
select
  name,
  display_name,
  conditions,
  location
from
  gcp_workstations_workstation_cluster
where
  conditions is not null;
```

```sql+sqlite
select
  name,
  display_name,
  conditions,
  location
from
  gcp_workstations_workstation_cluster
where
  conditions is not null;
```

### List recently created clusters
Identify clusters that were created recently to track new development environment deployments.

```sql+postgres
select
  name,
  display_name,
  create_time,
  location
from
  gcp_workstations_workstation_cluster
order by
  create_time desc
limit
  10;
```

```sql+sqlite
select
  name,
  display_name,
  create_time,
  location
from
  gcp_workstations_workstation_cluster
order by
  create_time desc
limit
  10;
```

