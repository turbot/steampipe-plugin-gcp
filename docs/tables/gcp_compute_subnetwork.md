---
title: "Steampipe Table: gcp_compute_subnetwork - Query Google Cloud Compute Engine Subnetworks using SQL"
description: "Allows users to query Google Cloud Compute Engine Subnetworks, providing insights into the configuration and status of each subnetwork."
folder: "Compute"
---

# Table: gcp_compute_subnetwork - Query Google Cloud Compute Engine Subnetworks using SQL

Google Cloud Compute Engine Subnetworks are regional resources, each within a specific region, that contain IP address ranges. Subnetworks can be used to partition the IP space of a network into segments, improving network security and efficiency. They are associated with a network and region, and can have policies that control outbound internet access.

## Table Usage Guide

The `gcp_compute_subnetwork` table provides insights into the subnetworks within Google Cloud Compute Engine. As a network administrator, explore subnetwork-specific details through this table, including IP ranges, associated network and region, and outbound internet access policies. Utilize it to uncover information about subnetworks, such as their configuration, status, and the partitioning of the IP space of a network.

## Examples

### Subnetwork basic info
Explore which subnetworks in your Google Cloud Platform have private IP Google access enabled. This can help determine areas where you may want to tighten security or reconfigure access permissions.

```sql+postgres
select
  name,
  gateway_address,
  ip_cidr_range,
  ipv6_cidr_range,
  private_ip_google_access,
  id,
  network_name
from
  gcp_compute_subnetwork;
```

```sql+sqlite
select
  name,
  gateway_address,
  ip_cidr_range,
  ipv6_cidr_range,
  private_ip_google_access,
  id,
  network_name
from
  gcp_compute_subnetwork;
```

### List of subnetworks where users have compute admin access assigned in a resource policy
Explore which subnetworks have users with compute admin access assigned, allowing you to understand the distribution of administrative privileges within your network resources. This query is useful for identifying potential security risks and ensuring appropriate access management.

```sql+postgres
select
  name,
  id,
  jsonb_array_elements_text(p -> 'members') as members,
  p ->> 'role' as role
from
  gcp_compute_subnetwork,
  jsonb_array_elements(iam_policy -> 'bindings') as p
where
  p ->> 'role' = 'roles/compute.admin';
```

```sql+sqlite
select
  s.name,
  s.id,
  json_extract(p.value, '$.members') as members,
  json_extract(p.value, '$.role') as role
from
  gcp_compute_subnetwork as s,
  json_each(iam_policy, '$.bindings') as p
where
  json_extract(p.value, '$.role') = 'roles/compute.admin';
```

### Secondary IP info of each subnetwork
Identify the secondary IP ranges within each subnetwork in your Google Cloud Platform. This can help you understand the distribution and usage of IP addresses within your network infrastructure.

```sql+postgres
select
  name,
  id,
  p ->> 'rangeName' as range_name,
  p ->> 'ipCidrRange' as ip_cidr_range
from
  gcp_compute_subnetwork,
  jsonb_array_elements(secondary_ip_ranges) as p;
```

```sql+sqlite
select
  s.name,
  s.id,
  json_extract(p.value, '$.rangeName') as range_name,
  json_extract(p.value, '$.ipCidrRange') as ip_cidr_range
from
  gcp_compute_subnetwork as s,
  json_each(secondary_ip_ranges) as p;
```

### Subnet count per network
Analyze your network configuration to understand how many subnetworks exist within each network. This can be useful for assessing the complexity and segmentation of your network infrastructure.

```sql+postgres
select
  network,
  count(*) as subnet_count
from
  gcp_compute_subnetwork
group by
  network;
```

```sql+sqlite
select
  network,
  count(*) as subnet_count
from
  gcp_compute_subnetwork
group by
  network;
```

### List subnetworks having VPC flow logging set to false
Determine the areas in which VPC flow logging is not enabled within your Google Cloud Platform subnetworks. This can help identify potential security vulnerabilities and improve network monitoring and troubleshooting.

```sql+postgres
select
  name,
  id,
  enable_flow_logs
from
  gcp_compute_subnetwork
where
  not enable_flow_logs;
```

```sql+sqlite
select
  name,
  id,
  enable_flow_logs
from
  gcp_compute_subnetwork
where
  not enable_flow_logs;
```

### IP Info subnets
Explore which subnetworks in your Google Cloud Platform (GCP) compute environment have specific characteristics. This query can help pinpoint the specific locations where the number of hosts per subnet may need to be adjusted for optimal network performance.

```sql+postgres
select
  name,
  id,
  ip_cidr_range,
  gateway_address,
  broadcast(ip_cidr_range),
  netmask(ip_cidr_range),
  network(ip_cidr_range),
  pow(2, 32 - masklen(ip_cidr_range)) -1 as hosts_per_subnet
from
  gcp_compute_subnetwork;
```

```sql+sqlite
Error: SQLite does not support CIDR operations.
```