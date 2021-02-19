# Table: gcp_service_account

A service account is a special type of Google account intended to represent a non-human user that needs to authenticate and be authorized to access data in Google APIs

## Examples

### List of email ids associated with the service account

```sql
select
  display_name,
  name as service_account,
  email
from
  gcp_service_account;
```


### Find service accounts with policies that grant public access

```sql
select
  name,
  split_part(s ->> 'role', '/', 2) as role,
  entity
from
  gcp_service_account,
  jsonb_array_elements(iam_policy -> 'bindings') as s,
  jsonb_array_elements_text(s -> 'members') as entity
where
  entity = 'allUsers'
  or entity = 'allAuthenticatedUsers';
```