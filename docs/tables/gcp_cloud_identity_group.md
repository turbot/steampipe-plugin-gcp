# Table: gcp_cloud_identity_group

A Membership defines a relationship between a Group and an entity belonging to that Group, referred to as a "member".

**You must specify the parent resource** in the `where` clause (`where parent='C046psxkn'`) to list the identity groups.

## Examples

### Basic info

```sql
select
  name,
  display_name,
  description,
  create_time,
  location,
  project
from
  gcp_cloud_identity_group
where
  parent = 'C046psxkn';
```

### Get details for a specific group

```sql
select
  name,
  display_name,
  description,
  create_time,
  location,
  project
from
  gcp_cloud_identity_group
where
  name = 'group_name';
```

### Get dynamic group settings

```sql
select
  name,
  display_name,
  dynamic_group_metadata ->> 'Status' as dynamic_group_status,
  queries ->> 'Query' as dynamic_group_query,
  queries ->> 'ResourceType' as dynamic_group_query_resource_type,
  project
from
  gcp_cloud_identity_group,
  jsonb_array_elements(dynamic_group_metadata -> 'Queries') as queries
where
  parent = 'C046psxkn';
```

### List groups created in the last 7 days

```sql
select
  name,
  display_name,
  description,
  create_time,
  location,
  project
from
  gcp_cloud_identity_group
where
  parent = 'C046psxkn'
  and create_time > now() - interval '7' day;
```
