# Table: gcp_bigquery_table

A BigQuery table contains individual records organized in rows. Each record is composed of columns (also called fields).

Every table is defined by a schema that describes the column names, data types, and other information. One can specify the schema of a table when it is created, or can create a table without a schema and declare the schema in the query job or load job that first populates it with data.

## Examples

### Basic info

```sql
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

```sql
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

```sql
select
  dataset_id,
  location
from
  gcp_bigquery_table
where
  tags -> 'owner' is null;
```
