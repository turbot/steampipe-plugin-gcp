---
title: "Steampipe Table: gcp_compute_route - Query Google Cloud Compute Routes using SQL"
description: "Allows users to query Google Cloud Compute Routes, providing details on the paths that network traffic takes from a virtual machine (VM) instance to other destinations."
folder: "Compute"
---

# Table: gcp_compute_route - Query Google Cloud Compute Routes using SQL

Google Cloud Compute Routes are a part of Google Cloud's networking infrastructure. They define the paths that network traffic takes from a virtual machine (VM) instance to other destinations. These destinations can be inside your virtual private cloud (VPC) network or outside of it.

## Table Usage Guide

The `gcp_compute_route` table provides insights into Google Cloud Compute Routes within Google Cloud Platform's networking infrastructure. As a network engineer, explore route-specific details through this table, including the network that the route applies to, the destination range of outgoing packets, and the next hop (the first stop on the way to the final destination). Utilize it to uncover information about routes, such as those with specific network tags, the priority of the routes, and the instances associated with the routes.


## Examples

### Route basic info
Explore the basics of your network routing configuration in Google Cloud Platform (GCP). This can help pinpoint specific areas that may require adjustments to optimize network traffic flow.

```sql+postgres
select
  name,
  dest_range,
  priority,
  network
from
  gcp_compute_route;
```

```sql+sqlite
select
  name,
  dest_range,
  priority,
  network
from
  gcp_compute_route;
```

### List of routes that are not applied to default network
Discover the segments that are not part of the default network. This is useful to identify any routes that may have been incorrectly assigned or overlooked during network configuration.

```sql+postgres
select
  name,
  id,
  network_name as network
from
  gcp_compute_route
where
  network_name <> 'default';
```

```sql+sqlite
select
  name,
  id,
  network_name as network
from
  gcp_compute_route
where
  network_name <> 'default';
```

### List of system-generated default routes
Explore instances where potential misconfigurations have been detected on routes, in order to proactively address any issues and maintain optimal network performance.

```sql+postgres
select
  name,
  dest_range,
  priority,
  next_hop_gateway
from
  gcp_compute_route
where
  priority = 1000
  and dest_range = '0.0.0.0/0';
```

```sql+sqlite
select
  name,
  dest_range,
  priority,
  next_hop_gateway
from
  gcp_compute_route
where
  priority = 1000
  and dest_range = '0.0.0.0/0';
```

# List of warning messages for potential misconfigurations detected on routes

```sql+postgres
select
  name,
  warnings
from
  gcp_compute_route
where
  warnings is not null;
```

```sql+sqlite
select
  name,
  warnings
from
  gcp_compute_route
where
  warnings is not null;
```