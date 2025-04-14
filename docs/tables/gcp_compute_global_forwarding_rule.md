---
title: "Steampipe Table: gcp_compute_global_forwarding_rule - Query GCP Compute Global Forwarding Rules using SQL"
description: "Allows users to query GCP Compute Global Forwarding Rules, providing insights into network traffic routing and load balancing configurations."
folder: "Compute"
---

# Table: gcp_compute_global_forwarding_rule - Query GCP Compute Global Forwarding Rules using SQL

A Global Forwarding Rule in Google Cloud Platform (GCP) is a component of the Cloud Load Balancing service. It is responsible for forwarding traffic from specified IP addresses to a target instance, target proxy, or target pool. These rules can be used to distribute incoming network traffic across multiple compute instances to ensure that no single instance is overwhelmed.

## Table Usage Guide

The `gcp_compute_global_forwarding_rule` table provides insights into the Global Forwarding Rules within Google Cloud Platform's Compute service. As a network engineer, you can use this table to explore details about each rule, including the IP addresses it handles, its target instances, and its associated load balancing configurations. This can be especially beneficial in optimizing your network traffic distribution and ensuring efficient load balancing across your compute instances.

## Examples

### Basic info
Gain insights into the details of global forwarding rules, such as their names, IDs, IP addresses, protocols, port ranges, and targets within the Google Cloud Platform. This can be useful in understanding the networking configuration and traffic routing in your cloud environment.

```sql+postgres
select
  name,
  id,
  ip_address,
  ip_protocol,
  port_range,
  target
from
  gcp_compute_global_forwarding_rule;
```

```sql+sqlite
select
  name,
  id,
  ip_address,
  ip_protocol,
  port_range,
  target
from
  gcp_compute_global_forwarding_rule;
```

### List global forwarding rules which are globally accessible
Determine the areas in which global forwarding rules are set to be globally accessible, allowing for a broadened network reach and enhanced connectivity. This can be particularly useful in understanding the scope of your network access and identifying potential security considerations.

```sql+postgres
select
  name,
  id,
  ip_address,
  allow_global_access
from
  gcp_compute_global_forwarding_rule
where
  allow_global_access;
```

```sql+sqlite
select
  name,
  id,
  ip_address,
  allow_global_access
from
  gcp_compute_global_forwarding_rule
where
  allow_global_access = '1';
```

### List global forwarding rules where mirroring collector is enabled (i.e load balancer can be used as a collector for packet mirroring)
Discover the segments that have enabled the packet mirroring feature, allowing the load balancer to collect data. This is useful in analyzing network traffic for security monitoring or troubleshooting.

```sql+postgres
select
  name,
  id,
  is_mirroring_collector
from
  gcp_compute_global_forwarding_rule
where
  is_mirroring_collector;
```

```sql+sqlite
select
  name,
  id,
  is_mirroring_collector
from
  gcp_compute_global_forwarding_rule
where
  is_mirroring_collector = '1';
```