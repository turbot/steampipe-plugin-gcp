---
title: "Steampipe Table: gcp_billing_account - Query GCP Billing Accounts using SQL"
description: "Allows users to query GCP Billing Accounts, providing insights into the billing information, payment status, and open status of each account."
folder: "Billing"
---

# Table: gcp_billing_account - Query GCP Billing Accounts using SQL

A GCP Billing Account is a resource in Google Cloud Platform that is linked to a Google payments profile. This resource is used to define who pays for a given set of Google Cloud resources and Google Maps Platform APIs. Access control to the billing account is established by IAM roles.

## Table Usage Guide

The `gcp_billing_account` table provides detailed insights into billing accounts within Google Cloud Platform. As a financial analyst or cloud administrator, explore billing account-specific details through this table, including billing information, payment status, and open status. Utilize it to uncover information about billing accounts, such as those with outstanding payments or to verify the open status of accounts.

**Important Notes**
- This table requires the `billing.viewer` permission to retrieve billing account details.

## Examples

### Basic info
Explore which Google Cloud Platform billing accounts are open and their associated projects and locations. This is useful for auditing purposes and to identify any potential cost-related issues.

```sql+postgres
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

```sql+sqlite
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
Explore the relationship between billing account members and their assigned roles within a Google Cloud Platform (GCP) environment. This can be useful for auditing purposes or to ensure appropriate access levels are maintained.

```sql+postgres
select
  name,
  display_name,
  jsonb_array_elements_text(p -> 'members') as member,
  p ->> 'role' as role
from
  gcp_billing_account,
  jsonb_array_elements(iam_policy -> 'bindings') as p;
```

```sql+sqlite
select
  name,
  display_name,
  json_extract(p.value, '$.members') as member,
  json_extract(p.value, '$.role') as role
from
  gcp_billing_account,
  json_each(iam_policy, '$.bindings') as p;
```

### List accounts whose members have billing admin access
Explore which accounts have members with billing admin access. This is useful to identify potential areas of financial risk and ensure appropriate access control.

```sql+postgres
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

```sql+sqlite
select
  name,
  display_name,
  json_extract(i.value, '$.members') as members,
  json_extract(i.value, '$.role') as role
from
  gcp_billing_account,
  json_each(iam_policy, '$.bindings') as i
where
  json_extract(i.value, '$.role') like '%billing.admin';
```

### List billing accounts that are open
Explore open billing accounts to gain insights into their associated projects and locations, helping you to effectively manage and monitor your financial resources in the GCP environment.

```sql+postgres
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

```sql+sqlite
select
  name,
  display_name,
  project,
  location
from
  gcp_billing_account
where
  open = 1;
```