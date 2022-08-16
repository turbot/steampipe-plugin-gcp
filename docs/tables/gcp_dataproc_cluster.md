# Table: gcp_dataproc_cluster

Dataproc is a fully managed and highly scalable service for running Apache Spark, Apache Flink, Presto, and 30+ open source tools and frameworks. Use Dataproc for data lake modernization, ETL, and secure data science, at planet scale, fully integrated with Google Cloud, at a fraction of the cost.

## Examples

### Basic info

```sql
select
  cluster_name,
  cluster_uuid,
  config,
  state,
  tags
from
  gcp_dataproc_cluster;
```

### List the clusters which are in error state

```sql
select
  cluster_name,
  cluster_uuid,
  state
from
  gcp_dataproc_cluster
where
  state = 'ERROR';
```

### Get config details of a cluster

```sql
select
  cluster_name,
  config -> 'endpointConfig' as endpoint_config,
  config -> 'configBucket' as config_bucket,
  config -> 'shieldedInstanceConfig' as shielded_instance_config,
  config -> 'masterConfig' as master_config
from
  gcp_dataproc_cluster
where
  cluster_name = 'cluster-5824';
```
