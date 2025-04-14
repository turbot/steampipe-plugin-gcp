---
title: "Steampipe Table: gcp_compute_forwarding_rule - Query GCP Compute Forwarding Rules using SQL"
description: "Allows users to query GCP Compute Forwarding Rules, providing information about their configurations and operational status."
folder: "Compute"
---

# Table: gcp_compute_forwarding_rule - Query GCP Compute Forwarding Rules using SQL

A GCP Compute Forwarding Rule is a resource within Google Cloud Platform's Compute Engine service. It specifies which network traffic is directed to which specific load balancer components. Forwarding rules are associated with specific IP addresses.

## Table Usage Guide

The `gcp_compute_forwarding_rule` table provides insights into forwarding rules within Google Cloud Platform's Compute Engine service. As a network engineer or system administrator, you can explore specific details about each forwarding rule, including their associated IP addresses, target proxies, and port ranges. Use this table to understand your network traffic direction and management within your GCP environment.

## Examples

### Basic info
Explore the configuration of your Google Cloud Platform's compute forwarding rules to gain insights into the load balancing scheme and backend service. This can be useful in determining areas where network traffic is being directed and ensuring optimal distribution of workload.Explore which IP addresses are associated with your load balancing scheme in Google Cloud Platform. This can help you understand how your network traffic is being directed and managed.

```sql+postgres
select
  name,
  id,
  self_link,
  backend_service,
  ip_address,
  load_balancing_scheme
from
  gcp_compute_forwarding_rule;
```

```sql+sqlite
select
  name,
  id,
  self_link,
  backend_service,
  ip_address,
  load_balancing_scheme
from
  gcp_compute_forwarding_rule;
```

### List of forwarding rules which are not globally accessible
Identify the forwarding rules that are not accessible globally. This query is useful for ensuring network security by pinpointing potential vulnerabilities.Explore which forwarding rules in your Google Cloud Platform (GCP) compute environment are not globally accessible. This can help ensure your network configuration aligns with your security and accessibility requirements.

```sql+postgres
select
  name,
  id,
  allow_global_access
from
  gcp_compute_forwarding_rule
where
  not allow_global_access;
```

```sql+sqlite
select
  name,
  id,
  allow_global_access
from
  gcp_compute_forwarding_rule
where
  allow_global_access = 0;
```

### List of EXTERNAL forwarding rules
Explore which forwarding rules are set to 'EXTERNAL' in the Google Cloud Platform's Compute Engine. This can help assess network traffic routing configurations for security or optimization purposes.Discover the segments that utilize external load balancing schemes within your Google Cloud Platform's forwarding rules. This can help manage traffic flow and optimize resource allocation within your network infrastructure.

```sql+postgres
select
  name,
  id,
  load_balancing_scheme
from
  gcp_compute_forwarding_rule
where
  load_balancing_scheme = 'EXTERNAL';
```

```sql+sqlite
select
  name,
  id,
  load_balancing_scheme
from
  gcp_compute_forwarding_rule
where
  load_balancing_scheme = 'EXTERNAL';
```