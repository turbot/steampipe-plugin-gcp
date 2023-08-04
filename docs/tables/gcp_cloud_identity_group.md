# Table: gcp_cloud_identity_group

A Membership defines a relationship between a Group and an entity belonging to that Group, referred to as a "member".

**You must specify the parent resource** under which to list all the `Group` resources. Must be of the form `identitysources/{identity_source}` for external-identity-mapped groups or `customers/{customer}` for Google Groups in the `where` clause (`where parent='C046psxkn'`) to list the identity group memberships.

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

### Get members in each group

```sql
select
  m.group_name,
  m.name as member_name
from
  gcp_cloud_identity_group_membership m,
  gcp_cloud_identity_group g
where
  m.group_name = g.name
  and g.parent = 'C046psxkn';
```

### Get total number of members in each group

```sql
select
  m.group_name,
  count(m.*) as total_members
from
  gcp_cloud_identity_group_membership m,
  gcp_cloud_identity_group g
where
  m.group_name = g.name
  and g.parent = 'C046psxkn'
group by
  m.group_name;
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