---
title: "Steampipe Table: gcp_compute_network - Query Google Cloud Platform Compute Networks using SQL"
description: "Allows users to query Compute Networks in Google Cloud Platform, specifically providing insights into network configurations, including subnetworks, firewall rules, and routing information."
folder: "Compute"
---

# Table: gcp_compute_network - Query Google Cloud Platform Compute Networks using SQL

A Google Cloud Platform Compute Network is a virtual version of the traditional physical networks that exist within and between physical data centers. A network provides the communication path between your Compute Engine virtual machine (VM) instances. They are global resources, spanning all regions, and are used to define the network topology, such as subnetworks and network peering connections.

## Table Usage Guide

The `gcp_compute_network` table provides insights into Compute Networks within Google Cloud Platform. As a network engineer or cloud architect, you can use this table to explore network-specific details, including its subnetworks, firewall rules, and routing configurations. This allows you to manage and optimize your network infrastructure effectively, ensuring secure and efficient communication paths within your Google Cloud environment.

## Examples

### List networks having auto create subnetworks feature disabled
Identify the networks that have the auto-create subnetworks feature turned off. This can be useful for assessing network configurations where manual subnet creation is preferred for more control over network segmentation.

```sql+postgres
select
  name,
  id,
  auto_create_subnetworks
from
  gcp_compute_network
where
  not auto_create_subnetworks;
```

```sql+sqlite
select
  name,
  id,
  auto_create_subnetworks
from
  gcp_compute_network
where
  auto_create_subnetworks = 0;
```

### List networks having routing_mode set to REGIONAL
Discover the segments that have their routing mode set to 'REGIONAL' within your network settings. This can be useful in understanding and managing network traffic flow within specific regions.

```sql+postgres
select
  name,
  id,
  routing_mode
from
  gcp_compute_network
where
  routing_mode = 'REGIONAL';
```

```sql+sqlite
select
  name,
  id,
  routing_mode
from
  gcp_compute_network
where
  routing_mode = 'REGIONAL';
```

### Subnets counts for each network
Explore which networks have the most subnets, allowing you to understand the distribution of subnets across your networks for better resource management and allocation.

```sql+postgres
select
  name,
  count(d) as num_subnets
from
  gcp_compute_network as i,
  jsonb_array_elements(subnetworks) as d
group by
  name;
```

```sql+sqlite
select
  g.name,
  count(d.value) as num_subnets
from
  gcp_compute_network g,
  json_each(g.subnetworks) as d
group by
  g.name;
```