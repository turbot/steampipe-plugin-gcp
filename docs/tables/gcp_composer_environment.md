---
title: "Steampipe Table: gcp_composer_environment - Query GCP Composer Environments using SQL"
description: "Allows users to query GCP Composer environments, providing detailed information on environment configurations, resources, and associated services."
folder: "Cloud Composer"
---

# Table: gcp_composer_environment - Query GCP Composer Environments using SQL

Google Cloud Composer is a fully managed workflow orchestration service built on Apache Airflow. The `gcp_composer_environment` table in Steampipe allows you to query information about Cloud Composer environments in your GCP environment, including details about Airflow configurations, node settings, security configurations, and more.

## Table Usage Guide

The `gcp_composer_environment` table is useful for cloud administrators, DevOps engineers, and data engineers who need to gather detailed insights into their Cloud Composer environments. You can query various aspects of the environment, such as its Airflow URI, node configurations, software settings, and network access controls. This table is particularly helpful for monitoring environment states, managing Airflow settings, and ensuring that your workflows are running in a secure and optimized environment.

## Examples

### Basic info
Retrieve basic information about Composer environments, including their name, location, and state.

```sql+postgres
select
  name,
  location,
  state,
  airflow_uri,
  create_time
from
  gcp_composer_environment;
```

```sql+sqlite
select
  name,
  location,
  state,
  airflow_uri,
  create_time
from
  gcp_composer_environment;
```

### List environments by size
Identify Composer environments based on their size, such as "ENVIRONMENT_SIZE_SMALL" or "ENVIRONMENT_SIZE_LARGE."

```sql+postgres
select
  name,
  environment_size,
  location,
  project
from
  gcp_composer_environment
where
  environment_size = 'ENVIRONMENT_SIZE_SMALL';
```

```sql+sqlite
select
  name,
  environment_size,
  location,
  project
from
  gcp_composer_environment
where
  environment_size = 'ENVIRONMENT_SIZE_SMALL';
```

### List environments with specific node configurations
Retrieve environments with specific node configurations, such as node count and machine types.

```sql+postgres
select
  name,
  node_count,
  node_config ->> 'machineType' as machine_type,
  location,
  project
from
  gcp_composer_environment
where
  node_count > 0;
```

```sql+sqlite
select
  name,
  node_count,
  json_extract(node_config, '$.machineType') as machine_type,
  location,
  project
from
  gcp_composer_environment
where
  node_count > 0;
```

### List environments with specific network access controls
Identify environments that have specific network-level access controls for the Airflow web server.

```sql+postgres
select
  name,
  jsonb_path_query_array(web_server_network_access_control, '$.allowedIpRanges[*].value') as allowed_ip_ranges,
  location,
  project
from
  gcp_composer_environment
where
  web_server_network_access_control is not null;
```

```sql+sqlite
select
  name,
  json_extract(web_server_network_access_control, '$.allowedIpRanges[0].value') as allowed_ip_ranges,
  location,
  project
from
  gcp_composer_environment
where
  web_server_network_access_control is not null;
```

### List environments with scheduled snapshot configurations
Fetch environments that have scheduled snapshot configurations enabled for recovery.

```sql+postgres
select
  name,
  recovery_config ->> 'scheduledSnapshotsConfig' as scheduled_snapshots_config,
  location,
  project
from
  gcp_composer_environment
where
  recovery_config ->> 'scheduledSnapshotsConfig' is not null;
```

```sql+sqlite
select
  name,
  json_extract(recovery_config, '$.scheduledSnapshotsConfig') as scheduled_snapshots_config,
  location,
  project
from
  gcp_composer_environment
where
  json_extract(recovery_config, '$.scheduledSnapshotsConfig') is not null;
```