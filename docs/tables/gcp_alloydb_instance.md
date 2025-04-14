---
title: "Steampipe Table: gcp_alloydb_instance - Query Google Cloud Platform AlloyDB Instances using SQL"
description: "Allows users to query Google Cloud Platform AlloyDB Instances, offering insights into instance configurations, status, and associated metadata."
folder: "AlloyDB"
---

# Table: gcp_alloydb_instance - Query Google Cloud Platform AlloyDB Instances using SQL

Google Cloud AlloyDB is a fully managed PostgreSQL-compatible database service, optimized for high performance with the reliability of Google's infrastructure. AlloyDB Instances within a cluster offer flexible options for scaling and redundancy, meeting the needs of demanding database applications.

## Table Usage Guide

The `gcp_alloydb_instance` table can be queried to retrieve detailed information about individual AlloyDB instances within a cluster. It is invaluable for database administrators and system architects looking to monitor instance-specific performance, understand configuration details, and manage resource allocation effectively.

## Examples

### Basic Information Query

Retrieve essential details about all AlloyDB instances to quickly get an overview of your instances' configurations and statuses.

```sql+postgres
select
  name,
  state,
  availability_type,
  display_name
from
  gcp_alloydb_instance;
```

```sql+sqlite
select
  name,
  state,
  availability_type,
  display_name
from
  gcp_alloydb_instance;
```

### List Instances by Availability Type

This query is useful for identifying instances according to their availability setup, which helps in understanding your database's fault tolerance and geographic distribution.

```sql+postgres
select
  name,
  availability_type,
  state
from
  gcp_alloydb_instance
where
  availability_type = 'REGIONAL';
```

```sql+sqlite
select
  name,
  availability_type,
  state
from
  gcp_alloydb_instance
where
  availability_type = 'REGIONAL';
```

### Find Instances with Specific Labels

Labels are key-value pairs assigned to instances and can be used for organization and access control. This query helps filter instances based on these labels.

```sql+postgres
select
  name,
  labels
from
  gcp_alloydb_instance
where
  labels -> 'environment' = 'production';
```

```sql+sqlite
select
  name,
  json_extract(labels, '$.environment') as environment
from
  gcp_alloydb_instance
where
  environment = 'production';
```

### Instances Currently in Maintenance

This query helps in identifying which instances are currently undergoing maintenance, allowing for better planning and reduced downtime impact.

```sql+postgres
select
  name,
  state
from
  gcp_alloydb_instance
where
  state = 'MAINTENANCE';
```

```sql+sqlite
select
  name,
  state
from
  gcp_alloydb_instance
where
  state = 'MAINTENANCE';
```

### Detailed View of Instance Configurations

Gain a deeper understanding of the configurations for a specific instance, useful for troubleshooting or planning upgrades.

```sql+postgres
select
  name,
  machine_config,
  client_connection_config,
  ip_address
from
  gcp_alloydb_instance
where
  name = 'instance-12345';
```

```sql+sqlite
select
  name,
  json_extract(machine_config, '$') as machine_config,
  json_extract(client_connection_config, '$') as client_connection_config,
  ip_address
from
  gcp_alloydb_instance
where
  name = 'instance-12345';
```

### Get node details of the instances

Retrieve information about the nodes.

```sql+postgres
select
  name,
  node -> 'Id' as node_id,
  node -> 'Ip' as node_ip,
  node -> 'State' as node_state,
  node -> 'ZoneId' as node_zone_id
from
  gcp_alloydb_instance,
  jsonb_array_elements(nodes) as node;
```

```sql+sqlite
select
  name,
  json_extract(node, '$.Id') as node_id,
  json_extract(node, '$.Ip') as node_ip,
  json_extract(node, '$.State') as node_state,
  json_extract(node, '$.ZoneId') as node_zone_id
from
  gcp_alloydb_instance,
  json_each(json(nodes)) as node;
```