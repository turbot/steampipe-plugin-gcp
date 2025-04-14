---
title: "Steampipe Table: gcp_compute_backend_bucket - Query Google Cloud Compute Backend Buckets using SQL"
description: "Allows users to query Backend Buckets in Google Cloud Compute, providing insights into the configuration, status, and other metadata related to backend buckets."
folder: "Compute"
---

# Table: gcp_compute_backend_bucket - Query Google Cloud Compute Backend Buckets using SQL

A Google Cloud Compute Backend Bucket is a Google Cloud service that allows you to store and serve static web content directly from a Google Cloud Storage bucket. It is typically used in conjunction with HTTP(S) Load Balancing and URL Maps to serve static website content, images, videos, or downloads. It is a globally available resource that scales to meet your needs and provides several benefits including low latency, high scalability, and cost-effectiveness.

## Table Usage Guide

The `gcp_compute_backend_bucket` table provides insights into Backend Buckets within Google Cloud Compute. As a DevOps engineer or system administrator, explore backend bucket-specific details through this table, including its configuration, status, and other metadata. Utilize it to monitor and manage the storage and serving of your static web content, ensuring optimal performance and cost-effectiveness.

## Examples

### Basic info
Explore the fundamental details of your Google Cloud Platform's compute backend buckets. This is useful for gaining insights into each bucket's name, unique identifier, and description.

```sql+postgres
select
  name,
  id,
  description,
  bucket_name
from
  gcp_compute_backend_bucket;
```

```sql+sqlite
select
  name,
  id,
  description,
  bucket_name
from
  gcp_compute_backend_bucket;
```

### List of backend buckets where cloud CDN is not enabled
Explore which backend buckets in your Google Cloud Platform configuration do not have the Cloud CDN feature enabled. This is useful for identifying potential areas to improve content delivery and network performance.

```sql+postgres
select
  name,
  id,
  enable_cdn
from
  gcp_compute_backend_bucket
where
  not enable_cdn;
```

```sql+sqlite
select
  name,
  id,
  enable_cdn
from
  gcp_compute_backend_bucket
where
  enable_cdn = 0;
```

### Backend bucket count per storage bucket
Analyze the distribution of backend buckets across storage buckets in your Google Cloud Platform. This can help to balance resource allocation and optimize storage utilization.

```sql+postgres
select
  bucket_name,
  count(*) as backend_bucket_count
from
  gcp_compute_backend_bucket
group by
  bucket_name;
```

```sql+sqlite
select
  bucket_name,
  count(*) as backend_bucket_count
from
  gcp_compute_backend_bucket
group by
  bucket_name;
```