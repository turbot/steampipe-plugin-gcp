---
title: "Steampipe Table: gcp_dataproc_cluster - Query Google Cloud Platform Dataproc Clusters using SQL"
description: "Allows users to query Google Cloud Platform Dataproc Clusters, providing insights into cluster configurations, status, and associated metadata."
folder: "Dataproc"
---

# Table: gcp_dataproc_cluster - Query Google Cloud Platform Dataproc Clusters using SQL

Google Cloud Dataproc is a fast, easy-to-use, fully managed cloud service for running Apache Spark and Apache Hadoop clusters in a simpler, more cost-efficient way. Operations that used to take hours or days take seconds or minutes instead, and you pay only for the resources you use. Dataproc also easily integrates with other Google Cloud services, giving you a powerful and complete data processing platform.

## Table Usage Guide

The `gcp_dataproc_cluster` table provides insights into Dataproc Clusters within Google Cloud Platform. As a data engineer, you can explore cluster-specific details through this table, including configurations, status, and associated metadata. Use it to uncover information about clusters, such as those with specific configurations, the operational status of clusters, and verification of associated metadata.

## Examples

### Basic info
Explore the configuration and status of your Google Cloud Platform's Dataproc clusters. This can help you assess the current state and settings of your clusters for better resource management and optimization.

```sql+postgres
select
  cluster_name,
  cluster_uuid,
  config,
  state,
  tags
from
  gcp_dataproc_cluster;
```

```sql+sqlite
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
Explore which clusters are experiencing errors to troubleshoot and resolve issues promptly, ensuring smooth operations. This is crucial in a real-world scenario where maintaining the health and functionality of clusters is vital for various applications and services.

```sql+postgres
select
  cluster_name,
  cluster_uuid,
  state
from
  gcp_dataproc_cluster
where
  state = 'ERROR';
```

```sql+sqlite
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
Explore the configuration details of a specific cluster to gain insights into various aspects like endpoint configuration, bucket configuration, shielded instance configuration, and master configuration. This can be particularly useful for understanding and managing the cluster's settings and configurations.

```sql+postgres
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

```sql+sqlite
select
  cluster_name,
  json_extract(config, '$.endpointConfig') as endpoint_config,
  json_extract(config, '$.configBucket') as config_bucket,
  json_extract(config, '$.shieldedInstanceConfig') as shielded_instance_config,
  json_extract(config, '$.masterConfig') as master_config
from
  gcp_dataproc_cluster
where
  cluster_name = 'cluster-5824';
```