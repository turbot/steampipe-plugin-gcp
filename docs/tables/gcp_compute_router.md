---
title: "Steampipe Table: gcp_compute_router - Query Google Cloud Platform Compute Routers using SQL"
description: "Allows users to query Compute Routers in Google Cloud Platform, providing insights into the configuration and status of cloud routers within a project."
folder: "Compute"
---

# Table: gcp_compute_router - Query Google Cloud Platform Compute Routers using SQL

A Compute Router in Google Cloud Platform is a resource that helps connect different networks within the cloud. It is a distributed, software-defined router that offers a reliable way to route traffic between virtual machines (VMs) and networks, regardless of their location. Compute Routers provide a scalable and flexible solution for managing network traffic in the cloud.

## Table Usage Guide

The `gcp_compute_router` table provides insights into Compute Routers within Google Cloud Platform. As a network administrator or cloud engineer, explore router-specific details through this table, including the network it is associated with, its region, and its operational status. Utilize it to manage and monitor your network's traffic routing configurations and to ensure optimal network performance.

## Examples

### Cloud router basic info
Explore the basic information related to cloud routers, such as their name and configuration settings. This could be useful for understanding their network behaviour and managing network routing protocols effectively.

```sql+postgres
select
  name,
  bgp_asn,
  bgp_advertise_mode
from
  gcp_compute_router;
```

```sql+sqlite
select
  name,
  bgp_asn,
  bgp_advertise_mode
from
  gcp_compute_router;
```

### NAT gateway info attached to router
Discover the settings of your NAT gateway that's linked to a router to understand its configuration and operational parameters. This can aid in network management by providing insights into features like endpoint independent mapping and IP allocation options.

```sql+postgres
select
  name,
  nat ->> 'name' as nat_name,
  nat ->> 'enableEndpointIndependentMapping' as enable_endpoint_independent_mapping,
  nat ->> 'natIpAllocateOption' as nat_ip_allocate_option,
  nat ->> 'sourceSubnetworkIpRangesToNat' as source_subnetwork_ip_ranges_to_nat
from
  gcp_compute_router,
  jsonb_array_elements(nats) as nat;
```

```sql+sqlite
select
  name,
  json_extract(nat.value, '$.name') as nat_name,
  json_extract(nat.value, '$.enableEndpointIndependentMapping') as enable_endpoint_independent_mapping,
  json_extract(nat.value, '$.natIpAllocateOption') as nat_ip_allocate_option,
  json_extract(nat.value, '$.sourceSubnetworkIpRangesToNat') as source_subnetwork_ip_ranges_to_nat
from
  gcp_compute_router,
  json_each(nats) as nat;
```

### List all routers with custom route advertisements
Explore which routers have custom route advertisements to better manage network traffic and understand your network's routing protocols. This is particularly useful when you want to assess the elements within your network that are using custom configurations for route advertisements.

```sql+postgres
select
  name,
  bgp_asn,
  bgp_advertise_mode,
  bgp_advertised_ip_ranges
from
  gcp_compute_router
where bgp_advertise_mode = 'CUSTOM';
```

```sql+sqlite
select
  name,
  bgp_asn,
  bgp_advertise_mode,
  bgp_advertised_ip_ranges
from
  gcp_compute_router
where bgp_advertise_mode = 'CUSTOM';
```