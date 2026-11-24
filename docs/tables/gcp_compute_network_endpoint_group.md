---
title: "Steampipe Table: gcp_compute_network_endpoint_group - Query GCP Compute Engine Network Endpoint Groups using SQL"
description: "Allows users to query GCP Compute Engine Network Endpoint Groups (zonal, regional and global), including their endpoint type, size, network and serverless or Private Service Connect targets."
folder: "Compute"
---

# Table: gcp_compute_network_endpoint_group - Query GCP Compute Engine Network Endpoint Groups using SQL

A network endpoint group (NEG) in GCP Compute Engine is a group of backend endpoints or services for a load balancer. Zonal NEGs hold VM or IP:port endpoints, regional NEGs point to serverless services or Private Service Connect targets, and global NEGs hold internet or hybrid endpoints.

## Table Usage Guide

The `gcp_compute_network_endpoint_group` table provides insights into zonal, regional and global network endpoint groups within Google Cloud Compute Engine. Use it to see which kinds of backends your load balancers use, how many endpoints each group has, and which network, Cloud Run service or Private Service Connect target a group points to. Individual endpoints are not listed, since listing them requires a POST call.

## Examples

### Basic info
Explore network endpoint groups, their type, size and location.

```sql+postgres
select
  name,
  network_endpoint_type,
  size,
  location_type,
  location
from
  gcp_compute_network_endpoint_group;
```

```sql+sqlite
select
  name,
  network_endpoint_type,
  size,
  location_type,
  location
from
  gcp_compute_network_endpoint_group;
```

### Count network endpoint groups by endpoint type
Understand which kinds of backends are used in the project.

```sql+postgres
select
  network_endpoint_type,
  count(*)
from
  gcp_compute_network_endpoint_group
group by
  network_endpoint_type;
```

```sql+sqlite
select
  network_endpoint_type,
  count(*)
from
  gcp_compute_network_endpoint_group
group by
  network_endpoint_type;
```

### List empty network endpoint groups
Find network endpoint groups that have no endpoints.

```sql+postgres
select
  name,
  network_endpoint_type,
  location
from
  gcp_compute_network_endpoint_group
where
  size = 0
  and network_endpoint_type in ('GCE_VM_IP', 'GCE_VM_IP_PORT', 'NON_GCP_PRIVATE_IP_PORT');
```

```sql+sqlite
select
  name,
  network_endpoint_type,
  location
from
  gcp_compute_network_endpoint_group
where
  size = 0
  and network_endpoint_type in ('GCE_VM_IP', 'GCE_VM_IP_PORT', 'NON_GCP_PRIVATE_IP_PORT');
```

### List serverless network endpoint groups with their Cloud Run service
See which Cloud Run services are exposed through a load balancer.

```sql+postgres
select
  name,
  location,
  cloud_run ->> 'service' as cloud_run_service
from
  gcp_compute_network_endpoint_group
where
  network_endpoint_type = 'SERVERLESS';
```

```sql+sqlite
select
  name,
  location,
  json_extract(cloud_run, '$.service') as cloud_run_service
from
  gcp_compute_network_endpoint_group
where
  network_endpoint_type = 'SERVERLESS';
```
