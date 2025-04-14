---
title: "Steampipe Table: gcp_bigquery_table - Query BigQuery Tables using SQL"
description: "Allows users to query BigQuery Tables in Google Cloud Platform (GCP), specifically the table details including creation time, expiration time, labels, and more."
folder: "BigQuery"
---

# Table: gcp_bigquery_table - Query BigQuery Tables using SQL

BigQuery is a fully-managed, serverless data warehouse that enables super-fast SQL queries using the processing power of Google's infrastructure. It allows you to analyze large datasets in real-time with a SQL-like syntax, without the need to manage the underlying infrastructure. BigQuery is designed to be not only fast, but also easy to use, and cost-effective.

## Table Usage Guide

The `gcp_bigquery_table` table provides insights into BigQuery Tables within Google Cloud Platform (GCP). As a data analyst or data engineer, explore table-specific details through this table, including creation time, expiration time, labels, and more. Utilize it to uncover information about tables, such as those with specific labels, the status of tables, and the verification of table metadata.

## Examples

### Basic info
Discover the segments that are part of your Google Cloud BigQuery data, including their creation time and geographical location. This can aid in understanding the distribution and age of your data across different regions.

```sql+postgres
select
  table_id,
  dataset_id,
  self_link,
  creation_time,
  location
from
  gcp_bigquery_table;
```

```sql+sqlite
select
  table_id,
  dataset_id,
  self_link,
  creation_time,
  location
from
  gcp_bigquery_table;
```

### List tables that are not encrypted using CMK
Discover the segments that are potentially vulnerable due to the absence of encryption using a Customer-Managed Key (CMK). This aids in identifying areas that require enhanced security measures.

```sql+postgres
select
  table_id,
  dataset_id,
  location,
  kms_key_name
from
  gcp_bigquery_table
where
  kms_key_name is null;
```

```sql+sqlite
select
  table_id,
  dataset_id,
  location,
  kms_key_name
from
  gcp_bigquery_table
where
  kms_key_name is null;
```

### List tables which do not have owner tag key
Discover the segments that lack an owner tag key within your Google Cloud Platform's BigQuery datasets. This can be beneficial for identifying potential areas of unclaimed resources or orphaned datasets that need attention.

```sql+postgres
select
  dataset_id,
  location
from
  gcp_bigquery_table
where
  tags -> 'owner' is null;
```

```sql+sqlite
select
  dataset_id,
  location
from
  gcp_bigquery_table
where
  json_extract(tags, '$.owner') is null;
```