# Table: gcp_kubernetes_node_pool

A node pool is a group of nodes within a cluster that all have the same configuration. Node pools use a NodeConfig specification. Each node in the pool has a Kubernetes node label, cloud.google.com/gke-nodepool , which has the node pool's name as its value.

## Examples

### Basic info

```sql
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

```sql
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


### List maximum pods for each node

```sql
select
  name,
  cluster_name,
  max_pods_constraint ->> 'maxPodsPerNode' as max_mods_per_node
from
  gcp_kubernetes_node_pool;
```
