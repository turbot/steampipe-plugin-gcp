# Table:  table_gcp_key_ring

A key ring organizes keys in a specific Google Cloud location and allows you to manage access control on groups of keys. A key ring's name does not need to be unique across a Google Cloud project, but must be unique within a given location. After creation, a key ring cannot be deleted. Key rings do not incur storage costs.

## Examples

### Basic info

```sql
select
  name,
  create_time
from
  gcp_kms_key_ring;
```


### List of key rings older than 30 days

```sql
select
  name,
  create_time
from
  gcp_kms_key_ring
where
  create_time <= (current_date - interval '30' day)
order by
  create_time;
```