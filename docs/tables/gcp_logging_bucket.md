---
title: "Steampipe Table: gcp_logging_bucket - Query Google Cloud Platform Logging Buckets using SQL"
description: "Allows users to query Logging Buckets in Google Cloud Platform, specifically the configuration and metadata, providing insights into logging data retention and management."
folder: "Cloud Logging"
---

# Table: gcp_logging_bucket - Query Google Cloud Platform Logging Buckets using SQL

A Logging Bucket in Google Cloud Platform is a container for logs. They provide a way to organize and control access to logs. These buckets are used to store logs based on retention and storage policies, ensuring that logs are kept and managed efficiently.

## Table Usage Guide

The `gcp_logging_bucket` table provides insights into Logging Buckets within Google Cloud Platform. As a system administrator, explore bucket-specific details through this table, including location, retention period, and associated metadata. Utilize it to manage and optimize your log data storage, understand your data retention policies, and ensure appropriate access controls are in place.

## Examples

### Basic info
Explore the lifecycle and retention details of your Google Cloud Platform logging buckets. This can help you understand how long logs are being retained and in what state, assisting in the management and planning of your logging strategy.

```sql+postgres
select
  name,
  lifecycle_state,
  description,
  retention_days
from
  gcp_logging_bucket;
```

```sql+sqlite
select
  name,
  lifecycle_state,
  description,
  retention_days
from
  gcp_logging_bucket;
```

### List locked buckets
Explore which Google Cloud Platform logging buckets are locked to prevent accidental deletion or alteration of crucial log data. This could be particularly useful for maintaining data security and integrity in your cloud environment.

```sql+postgres
select
  name,
  locked
from
  gcp_logging_bucket
where
  locked;
```

```sql+sqlite
select
  name,
  locked
from
  gcp_logging_bucket
where
  locked = 1;
```