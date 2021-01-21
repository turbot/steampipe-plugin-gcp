# Table:  gcp_iam_policy

An Identity and Access Management (IAM) policy, which specifies access controls for Google Cloud resources. ... Members can be user accounts, service accounts, Google groups, and domains (such as G Suite).

## Examples

### List of project members and there roles

```sql
select
title,
  p -> 'members' as member,
  split_part(p ->> 'role', '/', 4) as role
from
  gcp_iam_policy,
  jsonb_array_elements(bindings) as p;
```