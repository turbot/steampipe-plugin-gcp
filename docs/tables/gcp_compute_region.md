---
title: "Steampipe Table: gcp_compute_region - Query Google Cloud Platform Compute Regions using SQL"
description: "Allows users to query Google Cloud Platform Compute Regions, providing insights into the regional resources available for deploying computing resources."
folder: "Compute"
---

# Table: gcp_compute_region - Query Google Cloud Platform Compute Regions using SQL

A Compute Region in Google Cloud Platform is a specific geographical location where you can deploy your resources. Each region is a separate geographic area that consists of zones. Zones are essentially the deployment areas within a region where users can deploy their resources to ensure high availability and disaster recovery.

## Table Usage Guide

The `gcp_compute_region` table provides insights into the Compute Regions within Google Cloud Platform (GCP). As a cloud engineer or system administrator, use this table to explore region-specific details, including available zones, quotas, and related metadata. This can be particularly useful for planning resource deployment, ensuring high availability, and managing disaster recovery strategies.

## Examples

### List of compute regions which are down
Identify the specific regions in your Google Cloud Platform's compute engine that are currently non-operational. This is useful for quickly pinpointing areas of your infrastructure that may be causing disruptions or outages.

```sql+postgres
select
  name,
  id,
  status
from
  gcp_compute_region
where
  status = 'DOWN';
```

```sql+sqlite
select
  name,
  id,
  status
from
  gcp_compute_region
where
  status = 'DOWN';
```

### Get the quota info for a region (us-west1)
Analyze the settings to understand the quota limits for a specific region. This is useful for managing resources and preventing overuse.

```sql+postgres
select
  name,
  q -> 'metric' as quota_metric,
  q -> 'limit' as quota_limit
from
  gcp_compute_region,
  jsonb_array_elements(quotas) as q
where
  name = 'us-west1'
order by
  quota_metric;
```

```sql+sqlite
select
  name,
  json_extract(q.value, '$.metric') as quota_metric,
  json_extract(q.value, '$.limit') as quota_limit
from
  gcp_compute_region,
  json_each(quotas) as q
where
  name = 'us-west1'
order by
  quota_metric;
```

### Get the available zone info of each region
Explore which zones are available in each region to optimize resource allocation and manage your resources efficiently across different geographical locations.

```sql+postgres
select
  name,
  zone_names
from
  gcp_compute_region;
```

```sql+sqlite
select
  name,
  zone_names
from
  gcp_compute_region;
```

### Count the available zone in each region
Identify the number of available zones within each region to better distribute resources and maintain system balance.

```sql+postgres
select
  name,
  jsonb_array_length(zone_names) as zone_count
from
  gcp_compute_region;
```

```sql+sqlite
select
  name,
  json_array_length(zone_names) as zone_count
from
  gcp_compute_region;
```