# Table: gcp_compute_project_metadata

A project is used to organize resources in a Google Cloud Platform environment.

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


### Check whether oslogin is not enabled for the project

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