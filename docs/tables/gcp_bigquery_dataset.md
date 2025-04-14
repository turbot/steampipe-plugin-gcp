---
title: "Steampipe Table: gcp_bigquery_dataset - Query Google Cloud Platform BigQuery Datasets using SQL"
description: "Allows users to query BigQuery Datasets in Google Cloud Platform, providing insights into dataset configurations and access controls."
folder: "BigQuery"
---

# Table: gcp_bigquery_dataset - Query Google Cloud Platform BigQuery Datasets using SQL

Google Cloud Platform's BigQuery is a serverless, highly scalable, and cost-effective multi-cloud data warehouse designed for business agility. BigQuery Datasets are top-level containers that are used to organize and control access to tables and views. They provide a mechanism for grouping tables and views, and setting permissions at a group level.

## Table Usage Guide

The `gcp_bigquery_dataset` table provides insights into BigQuery Datasets within Google Cloud Platform. As a data analyst or data engineer, explore dataset-specific details through this table, including access controls, locations, and associated metadata. Utilize it to uncover information about datasets, such as their default partition type, default table expiration settings, and the labels applied to them.

## Examples

### Basic info
Explore the creation times and locations of your Google Cloud Platform BigQuery datasets. This is useful for understanding when and where your data resources were established, which can aid in resource management and optimization.

```sql+postgres
select
  dataset_id,
  self_link,
  creation_time,
  location
from
  gcp_bigquery_dataset;
```

```sql+sqlite
select
  dataset_id,
  self_link,
  creation_time,
  location
from
  gcp_bigquery_dataset;
```

### List datasets which are not encrypted using CMK
Identify instances where your datasets are not encrypted using a Customer-Managed Key (CMK). This helps in enhancing the security of your data by ensuring that all datasets are encrypted with a key that you manage and control.

```sql+postgres
select
  dataset_id,
  location,
  kms_key_name
from
  gcp_bigquery_dataset
where
  kms_key_name is null;
```

```sql+sqlite
select
  dataset_id,
  location,
  kms_key_name
from
  gcp_bigquery_dataset
where
  kms_key_name is null;
```

### List publicly accessible datasets
Discover the segments that are publicly accessible for data analysis and manipulation. This is useful in identifying potential data privacy risks or for making data publicly available for collaboration.

```sql+postgres
select
  dataset_id,
  location,
  ls as access_policy
from
  gcp_bigquery_dataset,
  jsonb_array_elements(access) as ls
where
  ls ->> 'specialGroup' = 'allAuthenticatedUsers'
  or ls ->> 'iamMember' = 'allUsers';
```

```sql+sqlite
select
  dataset_id,
  location,
  ls.value as access_policy
from
  gcp_bigquery_dataset,
  json_each(access) as ls
where
  json_extract(ls.value, '$.specialGroup') = 'allAuthenticatedUsers'
  or json_extract(ls.value, '$.iamMember') = 'allUsers';
```

### List datasets which do not have owner tag key
Discover datasets that lack an assigned owner within the Google Cloud Platform's BigQuery service. This query can be useful to identify potential gaps in data ownership and accountability.

```sql+postgres
select
  dataset_id,
  location
from
  gcp_bigquery_dataset
where
  tags -> 'owner' is null;
```

```sql+sqlite
select
  dataset_id,
  location
from
  gcp_bigquery_dataset
where
  json_extract(tags, '$.owner') is null;
```