# Table: gcp_billing_account

Cloud Billing accounts pay for usage costs in Google Cloud projects and Google Maps Platform projects.To use Google Cloud resources in a project, billing must be enabled on the project. Billing is enabled when the project is linked to an active Cloud Billing account.

**_Please note_**: This table requires the `billing.viewer` permission to retrieve billing account details.

## Examples

### Basic info

```sql
select
  name,
  display_name,
  master_billing_account,
  open,
  project,
  location
from
  gcp_billing_account;
```

### Get the billing account members and their associated IAM roles

```sql
select
  name,
  display_name,
  jsonb_array_elements_text(p -> 'members') as member,
  p ->> 'role' as role
from
  gcp_billing_account,
  jsonb_array_elements(iam_policy -> 'bindings') as p;
```

### List accounts whose members have billing admin access

```sql
select
  name,
  display_name,
  jsonb_array_elements_text(i -> 'members') as members,
  i ->> 'role' as role
from
  gcp_billing_account,
  jsonb_array_elements(iam_policy -> 'bindings') as i
where
  i ->> 'role' like '%billing.admin';
```

### List billing accounts that are open

```sql
select
  name,
  display_name,
  project,
  location
from
  gcp_billing_account
where
  open;
```
