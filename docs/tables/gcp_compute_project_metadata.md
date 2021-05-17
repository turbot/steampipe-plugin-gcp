# Table: gcp_compute_project_metadata

Compute project metadata authoritatively manages metadata common to all instances for a project in GCE.

## Examples

### Basic info

```sql
select
  name,
  id,
  default_service_account,
  creation_timestamp
from
  gcp_compute_project_metadata;
```


### Check if OS Login is enabled for Linux instances in the project

```sql
select
  name,
  id
from
  gcp_compute_project_metadata,
  jsonb_array_elements(common_instance_metadata -> 'items') as q
where
  common_instance_metadata -> 'items' @> '[{"key": "enable-oslogin"}]'
  and q ->> 'key' ilike 'enable-oslogin'
  and q ->> 'value' not ilike 'TRUE';
```
