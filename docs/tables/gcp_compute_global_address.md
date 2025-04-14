---
title: "Steampipe Table: gcp_compute_global_address - Query Google Cloud Compute Engine Global Addresses using SQL"
description: "Allows users to query Google Cloud Compute Engine Global Addresses, providing detailed information about the allocation method, network tier, and purpose of each address."
folder: "Compute"
---

# Table: gcp_compute_global_address - Query Google Cloud Compute Engine Global Addresses using SQL

Google Cloud Compute Engine Global Addresses are a resource in Google Cloud Platform that are used to reserve IP addresses for your project. These addresses can be used for various purposes such as HTTP(S), SSL proxy, and TCP proxy load balancing, or Cloud NAT. They can be either external or internal, with the allocation method either being automatic or manual.

## Table Usage Guide

The `gcp_compute_global_address` table provides insights into Global Addresses within Google Cloud Compute Engine. As a network engineer, explore address-specific details through this table, including the allocation method, network tier, and purpose of each address. Utilize it to manage and monitor the IP addresses reserved for your project, ensuring optimal utilization and network configuration.

## Examples

### List of internal address type global addresses
Discover the segments that consist of internal type global addresses in the Google Cloud Platform. This helps in understanding the distribution and usage of internal addresses within your cloud environment.

```sql+postgres
select
  name,
  id,
  address,
  address_type
from
  gcp_compute_global_address
where
  address_type = 'INTERNAL';
```

```sql+sqlite
select
  name,
  id,
  address,
  address_type
from
  gcp_compute_global_address
where
  address_type = 'INTERNAL';
```

### List of unused global addresses
Assess the elements within your Google Cloud Platform to identify unused global addresses. This can help optimize resource utilization and reduce costs by pinpointing areas of potential waste.

```sql+postgres
select
  name,
  address,
  status
from
  gcp_compute_global_address
where
  status <> 'IN_USE';
```

```sql+sqlite
select
  name,
  address,
  status
from
  gcp_compute_global_address
where
  status <> 'IN_USE';
```

### List of global addresses used for VPC peering
Explore which global addresses are used for the purpose of Virtual Private Cloud (VPC) peering. This is beneficial in understanding your network's connectivity and identifying potential areas for optimization or troubleshooting.

```sql+postgres
select
  name,
  address,
  purpose
from
  gcp_compute_global_address
where
  purpose = 'VPC_PEERING';
```

```sql+sqlite
select
  name,
  address,
  purpose
from
  gcp_compute_global_address
where
  purpose = 'VPC_PEERING';
```