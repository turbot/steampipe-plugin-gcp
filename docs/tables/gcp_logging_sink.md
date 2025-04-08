---
title: "Steampipe Table: gcp_logging_sink - Query GCP Logging Sinks using SQL"
description: "Allows users to query GCP Logging Sinks, providing a comprehensive view of the logging sinks configured in the Google Cloud environment."
folder: "Cloud Logging"
---

# Table: gcp_logging_sink - Query GCP Logging Sinks using SQL

Google Cloud Logging Sinks are mechanisms in Google Cloud Platform (GCP) that allow you to route log entries from logging services to a variety of supported destinations. These destinations can be Cloud Storage buckets, Pub/Sub topics, or BigQuery datasets. Logging sinks give you the flexibility to manage, analyze, and act on your log data as you see fit.

## Table Usage Guide

The `gcp_logging_sink` table provides insights into Logging Sinks within Google Cloud Platform (GCP). As a cloud engineer, you can explore sink-specific details through this table, including the destination, filter, and exclusion details. Utilize it to uncover information about sinks, such as their configured destinations, the filters applied, and to verify if any exclusions are set.

## Examples

### List writer identity that writes the export logs of logging sink
Identify the unique identities responsible for writing export logs in your logging sink. This can help monitor and manage who is contributing to your logs, enhancing security and accountability.

```sql+postgres
select
  name,
  unique_writer_identity
from
  gcp_logging_sink;
```

```sql+sqlite
select
  name,
  unique_writer_identity
from
  gcp_logging_sink;
```


### List the destination path for each sink
Explore which logging sinks are directing data to specific destinations in your Google Cloud Platform. This can help you understand where your log data is being sent and ensure it's reaching the intended targets.

```sql+postgres
select
  name,
  destination
from
  gcp_logging_sink;
```

```sql+sqlite
select
  name,
  destination
from
  gcp_logging_sink;
```