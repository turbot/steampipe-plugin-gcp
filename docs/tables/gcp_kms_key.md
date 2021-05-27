# Table: gcp_kms_key

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

### List keys older than 30 days

```sql
select
  name,
  create_time,
  rotation_period
from
  gcp_kms_key
where
  create_time <= (current_date - interval '30' day)
order by
  create_time;
```

### List keys where rotation period is 100000s

```sql
select
  name,
  create_time,
  rotation_period
from
  gcp_kms_key
where
  rotation_period = '100000s'
```
