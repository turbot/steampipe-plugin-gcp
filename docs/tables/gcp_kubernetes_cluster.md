---
title: "Steampipe Table: gcp_kubernetes_cluster - Query Google Cloud Platform Kubernetes Clusters using SQL"
description: "Allows users to query Kubernetes Clusters in Google Cloud Platform, specifically providing details about the cluster configurations, node pools, network settings, and more."
folder: "GKE"
---

# Table: gcp_kubernetes_cluster - Query Google Cloud Platform Kubernetes Clusters using SQL

A Kubernetes Cluster in Google Cloud Platform is a managed environment for deploying, managing, and scaling your containerized applications using Google infrastructure. The cluster consists of at least one cluster control plane and multiple worker machines called nodes. These nodes are Google Compute Engine virtual machines that run the Kubernetes processes necessary to make them part of the cluster.

## Table Usage Guide

The `gcp_kubernetes_cluster` table provides insights into Kubernetes Clusters within Google Cloud Platform. As a DevOps engineer, explore cluster-specific details through this table, including configurations, node pools, network settings, and more. Utilize it to uncover information about clusters, such as their status, zones, and the associated services and workloads.

## Examples

### Basic info
Explore which Kubernetes clusters in your Google Cloud Platform (GCP) are active and where they are located. This can help you manage resources and understand your network's geographical distribution.

```sql+postgres
select
  id,
  name,
  location_type,
  status,
  cluster_ipv4_cidr,
  max_pods_per_node,
  current_node_count,
  endpoint,
  location
from
  gcp_kubernetes_cluster;
```

```sql+sqlite
select
  id,
  name,
  location_type,
  status,
  cluster_ipv4_cidr,
  max_pods_per_node,
  current_node_count,
  endpoint,
  location
from
  gcp_kubernetes_cluster;
```

### List zonal clusters
Explore which Kubernetes clusters in your Google Cloud Platform are configured as zonal. This is useful to understand the geographical distribution of your resources and plan for redundancy or disaster recovery.

```sql+postgres
select
  name,
  location_type
from
  gcp_kubernetes_cluster
where
  location_type = 'ZONAL';
```

```sql+sqlite
select
  name,
  location_type
from
  gcp_kubernetes_cluster
where
  location_type = 'ZONAL';
```

### List clusters with node auto-upgrade enabled
Explore clusters that have the node auto-upgrade feature enabled. This is useful for ensuring your systems are always up-to-date with the latest features and security updates.

```sql+postgres
select
  name,
  location_type,
  n -> 'management' ->> 'autoUpgrade' node_auto_upgrade
from
  gcp_kubernetes_cluster,
  jsonb_array_elements(node_pools) as n
where
  n -> 'management' ->> 'autoUpgrade' = 'true';
```

```sql+sqlite
select
  name,
  location_type,
  json_extract(n.value, '$.management.autoUpgrade') as node_auto_upgrade
from
  gcp_kubernetes_cluster,
  json_each(node_pools) as n
where
  json_extract(n.value, '$.management.autoUpgrade') = 'true';
```

### List clusters with default service account
Identify instances where clusters are using the default service account in Google Cloud Platform's Kubernetes service. This can help in improving security by ensuring each cluster uses a unique service account.

```sql+postgres
select
  name,
  location_type,
  node_config ->> 'ServiceAccount' service_account
from
  gcp_kubernetes_cluster
where
  node_config ->> 'ServiceAccount' = 'default';
```

```sql+sqlite
select
  name,
  location_type,
  json_extract(node_config, '$.ServiceAccount') service_account
from
  gcp_kubernetes_cluster
where
  json_extract(node_config, '$.ServiceAccount') = 'default';
```

### List clusters with legacy authorization enabled
Determine the areas in which legacy authorization is still enabled on clusters. This is useful to identify potential security risks and areas for improvement in your Google Cloud Platform Kubernetes setup.

```sql+postgres
select
  name,
  location_type,
  legacy_abac_enabled
from
  gcp_kubernetes_cluster
where
  legacy_abac_enabled;
```

```sql+sqlite
select
  name,
  location_type,
  legacy_abac_enabled
from
  gcp_kubernetes_cluster
where
  legacy_abac_enabled = 1;
```

### List clusters with shielded nodes features disabled
Discover the segments that have the shielded nodes feature disabled in your Kubernetes clusters. This can help you identify potential security risks and enhance the protection of your clusters.

```sql+postgres
select
  name,
  location_type,
  shielded_nodes_enabled
from
  gcp_kubernetes_cluster
where
  not shielded_nodes_enabled;
```

```sql+sqlite
select
  name,
  location_type,
  shielded_nodes_enabled
from
  gcp_kubernetes_cluster
where
  shielded_nodes_enabled = 0;
```

### List clusters where secrets in etcd are not encrypted
Determine the areas in which sensitive information in your clusters is not secured. This is useful for identifying potential security vulnerabilities and ensuring data protection standards are met.

```sql+postgres
select
  name,
  database_encryption_state
from
  gcp_kubernetes_cluster
where
  database_encryption_state <> 'ENCRYPTED';
```

```sql+sqlite
select
  name,
  database_encryption_state
from
  gcp_kubernetes_cluster
where
  database_encryption_state <> 'ENCRYPTED';
```

### Node configuration of clusters
Explore the configuration settings of your clusters to understand their disk size, machine type, and other important parameters. This can be useful for optimizing your resources, ensuring security measures are in place, and maintaining efficient operation of your clusters.

```sql+postgres
select
  name,
  node_config ->> 'Disksizegb' as disk_size_gb,
  node_config ->> 'Disktype' as disk_type,
  node_config ->> 'Imagetype' as image_type,
  node_config ->> 'Machinetype' as machine_type,
  node_config ->> 'Disktype' as disk_type,
  node_config -> 'Metadata' ->> 'disable-legacy-endpoints' as disable_legacy_endpoints,
  node_config ->> 'Serviceaccount' as service_account,
  node_config -> 'Shieldedinstanceconfig' ->> 'EnableIntegrityMonitoring' as enable_integrity_monitoring
from
  gcp_kubernetes_cluster;
```

```sql+sqlite
select
  name,
  json_extract(node_config, '$.Disksizegb') as disk_size_gb,
  json_extract(node_config, '$.Disktype') as disk_type,
  json_extract(node_config, '$.Imagetype') as image_type,
  json_extract(node_config, '$.Machinetype') as machine_type,
  json_extract(node_config, '$.Disktype') as disk_type,
  json_extract(json_extract(node_config, '$.Metadata'), '$.disable-legacy-endpoints') as disable_legacy_endpoints,
  json_extract(node_config, '$.ServiceAccount') as service_account,
  json_extract(json_extract(node_config, '$.ShieldedInstanceConfig'), '$.EnableIntegrityMonitoring') as enable_integrity_monitoring
from
  gcp_kubernetes_cluster;
```