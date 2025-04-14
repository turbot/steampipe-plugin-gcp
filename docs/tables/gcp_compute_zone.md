---
title: "Steampipe Table: gcp_compute_zone - Query Google Cloud Compute Engine Zones using SQL"
description: "Allows users to query Google Cloud Compute Engine Zones, providing information on each zone's status, region, and available CPU platforms."
folder: "Compute"
---

# Table: gcp_compute_zone - Query Google Cloud Compute Engine Zones using SQL

Google Cloud Compute Engine Zones are geographical locations where Google Cloud resources are deployed and managed. Each zone is a deployment area within a region and is designed to be isolated from failures in other zones. Zones are ideal for deploying high availability applications and for distributing resources to provide disaster recovery.

## Table Usage Guide

The `gcp_compute_zone` table provides insights into the zones within Google Cloud Compute Engine. As a Cloud Engineer, you can explore zone-specific details through this table, including the zone's status, region, and available CPU platforms. Utilize it to manage resource allocation and distribution, and to plan for disaster recovery and high availability applications.

## Examples

### Compute zone basic info
Explore which zones within your Google Cloud Platform's compute engine are active and where they are located. This is useful for understanding the distribution and status of your compute resources.

```sql+postgres
select
  name,
  id,
  region_name,
  status
from
  gcp_compute_zone;
```

```sql+sqlite
select
  name,
  id,
  region_name,
  status
from
  gcp_compute_zone;
```

### Get the available cpu platforms in each zone
Determine the areas in which different CPU platforms are available in each zone to optimize resource allocation and performance.

```sql+postgres
select
  name,
  available_cpu_platforms
from
  gcp_compute_zone;
```

```sql+sqlite
select
  name,
  available_cpu_platforms
from
  gcp_compute_zone;
```

### Get the zones which are down
Explore which zones in your GCP Compute environment are currently down. This is useful for quickly identifying areas of your infrastructure that may be experiencing issues.

```sql+postgres
select
  name,
  id,
  region_name,
  status
from
  gcp_compute_zone
where
  status = 'DOWN';
```

```sql+sqlite
select
  name,
  id,
  region_name,
  status
from
  gcp_compute_zone
where
  status = 'DOWN';
```