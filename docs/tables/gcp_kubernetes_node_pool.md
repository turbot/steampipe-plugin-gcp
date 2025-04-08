---
title: "Steampipe Table: gcp_kubernetes_node_pool - Query GCP Kubernetes Node Pools using SQL"
description: "Allows users to query GCP Kubernetes Node Pools, providing details on configuration and status of each node pool within the cluster."
folder: "GKE"
---

# Table: gcp_kubernetes_node_pool - Query GCP Kubernetes Node Pools using SQL

A Kubernetes Node Pool is a group of nodes within a Google Cloud Kubernetes Engine (GKE) cluster that all have the same configuration. Node pools use a NodeConfig specification. They create a group of Google Compute Engine VM instances that serve as worker nodes for the applications running on your GKE clusters.

## Table Usage Guide

The `gcp_kubernetes_node_pool` table provides insights into the node pools within GCP Kubernetes Engine. As a DevOps or cloud engineer, explore node pool-specific details through this table, including node configurations, statuses, and associated metadata. Utilize it to uncover information about node pools, such as their configurations, the statuses of the nodes, and the details of the instances running the nodes.

## Examples

### Basic info
Explore the status and details of your Google Cloud Platform's Kubernetes node pools, such as the initial node count, version, and location. This can help you manage your resources effectively and understand the configuration of your clusters better.

```sql+postgres
select
  name,
  cluster_name,
  initial_node_count,
  version,
  status,
  location
from
  gcp_kubernetes_node_pool;
```

```sql+sqlite
select
  name,
  cluster_name,
  initial_node_count,
  version,
  status,
  location
from
  gcp_kubernetes_node_pool;
```

### List configuration info of each node
Explore the configuration details of each node within a Kubernetes cluster. This is useful to assess and manage resource allocation, such as disk size, machine type, and image type, and to review specific configurations like legacy endpoint usage and integrity monitoring settings.

```sql+postgres
select
  name,
  cluster_name,
  config ->> 'diskSizeGb' as disk_size_gb,
  config ->> 'diskType' as disk_type,
  config ->> 'imageType' as image_type,
  config ->> 'machineType' as machine_type,
  config -> 'metadata' ->> 'disable-legacy-endpoints' as disable_legacy_endpoints,
  config ->> 'serviceAccount' as machine_type,
  config -> 'shieldedInstanceConfig' ->> 'enableIntegrityMonitoring' as enable_integrity_monitoring
from
  gcp_kubernetes_node_pool;
```

```sql+sqlite
select
  name,
  cluster_name,
  json_extract(config, '$.diskSizeGb') as disk_size_gb,
  json_extract(config, '$.diskType') as disk_type,
  json_extract(config, '$.imageType') as image_type,
  json_extract(config, '$.machineType') as machine_type,
  json_extract(config, '$.metadata.disable-legacy-endpoints') as disable_legacy_endpoints,
  json_extract(config, '$.serviceAccount') as machine_type,
  json_extract(config, '$.shieldedInstanceConfig.enableIntegrityMonitoring') as enable_integrity_monitoring
from
  gcp_kubernetes_node_pool;
```

### List maximum pods for each node
Determine the capacity of each node in your Kubernetes cluster by identifying the maximum number of pods each node can run. This helps in efficient resource allocation and load balancing within the cluster.

```sql+postgres
select
  name,
  cluster_name,
  max_pods_constraint ->> 'maxPodsPerNode' as max_mods_per_node
from
  gcp_kubernetes_node_pool;
```

```sql+sqlite
select
  name,
  cluster_name,
  json_extract(max_pods_constraint, '$.maxPodsPerNode') as max_mods_per_node
from
  gcp_kubernetes_node_pool;
```

### List of all zonal node pools
Explore which node pools in your Google Cloud Platform Kubernetes service are zonal. This can help you manage and optimize your resources, as zonal node pools can offer different benefits and limitations compared to regional ones.

```sql+postgres
select
  name,
  location_type
from
  gcp_kubernetes_node_pool
where
  location_type = 'ZONAL';
```

```sql+sqlite
select
  name,
  location_type
from
  gcp_kubernetes_node_pool
where
  location_type = 'ZONAL';
```