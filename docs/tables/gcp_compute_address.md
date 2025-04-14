---
title: "Steampipe Table: gcp_compute_address - Query Google Cloud Compute Engine Addresses using SQL"
description: "Allows users to query Addresses in Google Cloud Compute Engine, specifically details related to their configuration, status, and associated resources."
folder: "Compute"
---

# Table: gcp_compute_address - Query Google Cloud Compute Engine Addresses using SQL

An Address in Google Cloud Compute Engine is a static, external IP address that you reserve and assign to your instances. These reserved addresses can be used by cloud resources in the same region. They provide a reliable endpoint for your cloud services and are essential for load balancing and other network configurations.

## Table Usage Guide

The `gcp_compute_address` table provides insights into the static, external IP addresses reserved within Google Cloud Compute Engine. As a network administrator, explore address-specific details through this table, including their status, type, and associated instances. Utilize it to oversee the allocation and usage of addresses, ensuring optimal network configurations and resource allocation.

## Examples

### Basic info
This query can be used to gain insights into the various IP addresses associated with your Google Cloud Platform (GCP) Compute Engine. It helps in understanding the type, creation time, IP version, status, associated subnetwork, and location of each address, which can be beneficial for managing and optimizing your network infrastructure.

```sql+postgres
select
  address,
  id,
  address_type,
  creation_timestamp,
  ip_version,
  status,
  subnetwork,
  location
from
  gcp_compute_address;
```

```sql+sqlite
select
  address,
  id,
  address_type,
  creation_timestamp,
  ip_version,
  status,
  subnetwork,
  location
from
  gcp_compute_address;
```

### List of address which are not in use
Discover the segments that consist of unused addresses in your Google Cloud Platform Compute Engine. This can help manage resources effectively by identifying and potentially freeing up unused addresses.

```sql+postgres
select
  address,
  address_type,
  creation_timestamp,
  status
from
  gcp_compute_address where status != 'IN_USE' ;
```

```sql+sqlite
select
  address,
  address_type,
  creation_timestamp,
  status
from
  gcp_compute_address where status != 'IN_USE' ;
```

### Address count by each network_tier
Analyze the distribution of addresses across different network tiers in your Google Cloud Platform to better understand your resource allocation. This could be useful in optimizing network performance and managing costs effectively.

```sql+postgres
select
  network_tier,
  count(*)
from
  gcp_compute_address
group by
  network_tier
order by network_tier;
```

```sql+sqlite
select
  network_tier,
  count(*)
from
  gcp_compute_address
group by
  network_tier
order by network_tier;
```

### Get details of users that are using an address
Gain insights into the details of users who are utilizing a specific address. This can be particularly useful for understanding user distribution and managing resources efficiently.

```sql+postgres
select
  name,
  address,
  id,
  jsonb_pretty(users)
from
  gcp_compute_address where name= 'test2';
```

```sql+sqlite
select
  name,
  address,
  id,
  users
from
  gcp_compute_address where name= 'test2';
```