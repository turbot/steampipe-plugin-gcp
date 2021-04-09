# Table: gcp_bigquery_dataset

Datasets are top-level containers that are used to organize and control access to your tables and views.

## Examples

### Basic info

```sql
select
  dataset_id,
  self_link,
  creation_time,
  location
from
  gcp_bigquery_dataset;
```

### List datasets which are not encrypted using CMK

```sql
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

```sql
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

### List datasets which do not have owner tag key

```sql
select
  dataset_id,
  location
from
  gcp_bigquery_dataset
where
  tags -> 'owner' is null;
```
