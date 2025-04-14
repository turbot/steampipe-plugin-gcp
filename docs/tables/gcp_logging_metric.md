---
title: "Steampipe Table: gcp_logging_metric - Query Google Cloud Platform Logging Metrics using SQL"
description: "Allows users to query Logging Metrics in Google Cloud Platform, specifically retrieving information about user-defined metrics based on logs, which can be used for monitoring and automating responses to specific log entries."
folder: "Cloud Logging"
---

# Table: gcp_logging_metric - Query Google Cloud Platform Logging Metrics using SQL

Google Cloud Platform's Logging Metrics is a feature that allows users to create user-defined metrics based on logs. These metrics can be used for monitoring and automating responses to specific log entries. It provides a way to track the count of log entries that match specific filters.

## Table Usage Guide

The `gcp_logging_metric` table provides insights into user-defined metrics in Google Cloud Platform's Logging. As a system administrator or DevOps engineer, explore metric-specific details through this table, including metric descriptors, metric type, and associated metadata. Utilize it to monitor and automate responses to specific log entries, helping to ensure system stability and performance.

## Examples

### Filter info of each metric
Explore which logging metrics are currently set up in your Google Cloud Platform (GCP) environment. This can help you understand what kind of data is being collected and monitored, which is critical for effective system management and troubleshooting.

```sql+postgres
select
  name as metric_name,
  description,
  filter
from
  gcp_logging_metric;
```

```sql+sqlite
select
  name as metric_name,
  description,
  filter
from
  gcp_logging_metric;
```

### Bucket configuration details of the logging metrics
Explore the configuration of logging metrics in Google Cloud Platform, focusing on specific parameters such as growth factor, scale, offset, and width. This can be useful in understanding the structure and distribution of data within your logging metrics, aiding in efficient data management and analysis.

```sql+postgres
select
  name,
  exponential_buckets_options_growth_factor,
  exponential_buckets_options_num_finite_buckets,
  exponential_buckets_options_scale,
  linear_buckets_options_num_finite_buckets,
  linear_buckets_options_offset,
  linear_buckets_options_width,
  explicit_buckets_options_bounds
from
  gcp_logging_metric;
```

```sql+sqlite
select
  name,
  exponential_buckets_options_growth_factor,
  exponential_buckets_options_num_finite_buckets,
  exponential_buckets_options_scale,
  linear_buckets_options_num_finite_buckets,
  linear_buckets_options_offset,
  linear_buckets_options_width,
  explicit_buckets_options_bounds
from
  gcp_logging_metric;
```