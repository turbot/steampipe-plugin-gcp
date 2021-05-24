# Table:  gcp_kubernetes_cluster

A cluster is the foundation of Google Kubernetes Engine (GKE): the Kubernetes objects that represent the containerized applications all run on top of a cluster. In GKE, a cluster consists of at least one control plane and multiple worker machines called nodes.

## Examples

### Basic info

```sql
select
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


### List of all zonal clusters

```sql
select
  name,
  location_type
from
  gcp_kubernetes_cluster
where
  location_type = 'Zonal';
```


### List clusters where node auto upgrade is enabled

```sql
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


### List clusters which uses default Service Account

```sql
select
  name,
  location_type,
  node_config ->> 'serviceAccount' service_account
from
  gcp_kubernetes_cluster
where
  node_config ->> 'serviceAccount' = 'default';
```


### List clusters where legacy authorization is enabled

```sql
select
  name,
  location_type,
  legacy_abac_enabled
from
  gcp_kubernetes_cluster
where
  legacy_abac_enabled;
```


### List clusters where shielded nodes features are disabled

```sql
select
  name,
  location_type,
  shielded_nodes_enabled
from
  gcp_kubernetes_cluster
where
  not shielded_nodes_enabled;
```


### List clusters where secrets in etcd are not encrypted

```sql
select
  name,
  database_encryption_state
from
  gcp_kubernetes_cluster
where
  database_encryption_state <> 'ENCRYPTED';
```


### Node configuration of clusters

```sql
select
  name,
  node_config ->> 'diskSizeGb' as disk_size_gb,
  node_config ->> 'diskType' as disk_type,
  node_config ->> 'imageType' as image_type,
  node_config ->> 'machineType' as machine_type,
  node_config ->> 'diskType' as disk_type,
  node_config -> 'metadata' ->> 'disable-legacy-endpoints' as disable_legacy_endpoints,
  node_config ->> 'serviceAccount' as service_account,
  node_config -> 'shieldedInstanceConfig' ->> 'enableIntegrityMonitoring' as enable_integrity_monitoring
from
  gcp_kubernetes_cluster;
```
