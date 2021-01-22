# Table:  gcp_iam_policy

An Identity and Access Management (IAM) policy, which specifies access controls for Google Cloud resources. ... Members can be user accounts, service accounts, Google groups, and domains (such as G Suite).

## Examples

### List of project members with their roles

```sql
select
  entity,
  p ->> 'role' as role
from
  gcp_iam_policy,
  jsonb_array_elements(bindings) as p,
  jsonb_array_elements_text(p -> 'members') as entity;
```


### List of members with owner roles

```sql
select
  entity,
  p ->> 'role' as role
from
  gcp_iam_policy,
  jsonb_array_elements(bindings) as p,
  jsonb_array_elements_text(p -> 'members') as entity
where
  split_part(p ->> 'role', '/', 2) = 'owner';
```