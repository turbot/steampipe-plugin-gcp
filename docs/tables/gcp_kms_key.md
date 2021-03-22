# Table:  table_gcp_key

A Cloud KMS key is a named object containing one or more key versions, along with metadata for the key. A key exists on exactly one key ring tied to a specific location.

## Examples

### Basic info

```sql
select
  name,
  create_time,
  rotation_period
from
  gcp_kms_key;
```


### List of keys older than 30 days

```sql
select
  name,
  create_time
from
  gcp_kms_key
where
  create_time <= (create_time - interval '30' day)
order by
  name;
```