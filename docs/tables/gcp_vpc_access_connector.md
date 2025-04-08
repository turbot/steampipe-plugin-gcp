---
title: "Steampipe Table: gcp_vpc_access_connector - Query GCP VPC Access Connectors using SQL"
description: "Allows users to query GCP VPC Access Connectors, providing detailed information on connector configurations, associated projects, and network settings."
folder: "VPC Access"
---

# Table: gcp_vpc_access_connector - Query GCP VPC Access Connectors using SQL

Google Cloud VPC Access Connector provides a way to enable serverless applications to connect securely to your Virtual Private Cloud (VPC) network. The `gcp_vpc_access_connector` table in Steampipe allows you to query information about VPC Access Connectors in your GCP environment, including their IP ranges, network settings, and associated projects.

## Table Usage Guide

The `gcp_vpc_access_connector` table is useful for cloud administrators and network engineers who need to gather detailed insights into their VPC Access Connectors. You can query various aspects of the connectors, such as their machine types, throughput configurations, state, and associated projects. This table is particularly useful for managing and monitoring network configurations, ensuring secure connectivity, and optimizing resource usage.

## Examples

### Basic info
Retrieve basic information about VPC Access Connectors, including their name, location, and state.

```sql+postgres
select
  name,
  location,
  state,
  network,
  machine_type
from
  gcp_vpc_access_connector;
```

```sql+sqlite
select
  name,
  location,
  state,
  network,
  machine_type
from
  gcp_vpc_access_connector;
```

### List connectors with specific IP CIDR ranges
Identify connectors that are using specific IP CIDR ranges, which can help in managing IP address allocation and avoiding conflicts.

```sql+postgres
select
  name,
  ip_cidr_range,
  network,
  location
from
  gcp_vpc_access_connector
where
  ip_cidr_range = '10.8.0.0/28';
```

```sql+sqlite
select
  name,
  ip_cidr_range,
  network,
  location
from
  gcp_vpc_access_connector
where
  ip_cidr_range = '10.8.0.0/28';
```

### List connectors by network and throughput
Retrieve connectors that are part of a specific VPC network and have a specific throughput configuration, which can be useful for optimizing network performance.

```sql+postgres
select
  name,
  network,
  min_throughput,
  max_throughput
from
  gcp_vpc_access_connector
where
  network = 'default'
  and max_throughput >= 1000;
```

```sql+sqlite
select
  name,
  network,
  min_throughput,
  max_throughput
from
  gcp_vpc_access_connector
where
  network = 'default'
  and max_throughput >= 1000;
```

### List the projects associated with the connectors
Identify VPC Access Connectors that are being used by specific projects, which can help in understanding project dependencies and managing access.

```sql+postgres
select
  name,
  jsonb_array_elements_text(connected_projects) as project_name,
  network,
  location
from
  gcp_vpc_access_connector
where
  connected_projects is not null;
```

```sql+sqlite
select
  name,
  json_extract(connected_projects, '$[0]') as project_name,
  network,
  location
from
  gcp_vpc_access_connector
where
  connected_projects is not null;
```

### List connectors by state
Retrieve a list of connectors filtered by their state and project, which can help in monitoring the status of connectors in specific environments.

```sql+postgres
select
  name,
  state,
  project,
  location
from
  gcp_vpc_access_connector
where
  state = 'READY';
```

```sql+sqlite
select
  name,
  state,
  project,
  location
from
  gcp_vpc_access_connector
where
  state = 'READY';
```

### Connectors with their associated subnets
Retrieve information about VPC Access Connectors and their associated subnets.

```sql+postgres
select
  c.name as connector_name,
  c.location,
  c.network,
  s ->> 'name' as subnet_name,
  s ->> 'ipCidrRange' as subnet_ip_range
from
  gcp_vpc_access_connector c,
  jsonb_array_elements(c.subnet) as s;
```

```sql+sqlite
select
  c.name as connector_name,
  c.location,
  c.network,
  json_extract(s.value, '$.name') as subnet_name,
  json_extract(s.value, '$.ipCidrRange') as subnet_ip_range
from
  gcp_vpc_access_connector c,
  json_each(c.subnet) as s;
```
