---
title: "Steampipe Table: gcp_bigquery_job - Query GCP BigQuery Jobs using SQL"
description: "Allows users to query BigQuery Jobs in Google Cloud Platform (GCP), providing insights into job configuration, statistics, and status."
folder: "BigQuery"
---

# Table: gcp_bigquery_job - Query GCP BigQuery Jobs using SQL

BigQuery Jobs in Google Cloud Platform (GCP) represent actions such as SQL queries, load, export, copy, or other types of jobs. These jobs are used to manage asynchronous tasks such as SQL statements and data import/export. The jobs are primarily used for query and load operations and can be used to monitor the status of various operations.

## Table Usage Guide

The `gcp_bigquery_job` table provides insights into BigQuery Jobs within Google Cloud Platform (GCP). As a data analyst or data engineer, explore job-specific details through this table, including job configuration, statistics, and status. Utilize it to monitor the progress of data operations, understand the configuration of specific jobs, and analyze the overall performance of BigQuery operations.

## Examples

### Basic info
Explore the creation times and locations of jobs within Google Cloud's BigQuery service. This can help understand when and where specific tasks were initiated, providing valuable insight into resource usage and operational trends.

```sql+postgres
select
  job_id,
  self_link,
  creation_time,
  location
from
  gcp_bigquery_job;
```

```sql+sqlite
select
  job_id,
  self_link,
  creation_time,
  location
from
  gcp_bigquery_job;
```

### List running jobs
Explore which jobs are currently active in your Google Cloud BigQuery environment. This can help you manage resources effectively and monitor ongoing processes.

```sql+postgres
select
  job_id,
  self_link,
  creation_time,
  location
from
  gcp_bigquery_job
where
  state = 'RUNNING';
```

```sql+sqlite
select
  job_id,
  self_link,
  creation_time,
  location
from
  gcp_bigquery_job
where
  state = 'RUNNING';
```