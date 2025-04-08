---
title: "Steampipe Table: gcp_service_account - Query Google Cloud Platform Service Accounts using SQL"
description: "Allows users to query Service Accounts in Google Cloud Platform, specifically the key details and permissions, providing insights into service account usage and security configurations."
folder: "Service Account"
---

# Table: gcp_service_account - Query Google Cloud Platform Service Accounts using SQL

A Service Account in Google Cloud Platform is a special type of account used by an application or a virtual machine (VM) instance, not a person. Applications use service accounts to make authorized API calls, authorized as either the service account itself, or as G Suite or Cloud Identity users through domain-wide delegation. These accounts can be created and managed by users, and they are tied to the lifecycle of the project in which they are created.

## Table Usage Guide

The `gcp_service_account` table provides insights into Service Accounts within Google Cloud Platform. As a security engineer, explore service account-specific details through this table, including permissions, roles, and associated metadata. Utilize it to uncover information about service accounts, such as those with excessive permissions, the roles assigned to each service account, and the verification of security configurations.

## Examples

### List of email ids associated with the service account
Explore which email IDs are linked to your service account to maintain a clear record of associated users. This can be particularly useful for managing access permissions and auditing user activities.

```sql+postgres
select
  display_name,
  name as service_account,
  email
from
  gcp_service_account;
```

```sql+sqlite
select
  display_name,
  name as service_account,
  email
from
  gcp_service_account;
```

### Find service accounts with policies that grant public access
Determine the areas in which service accounts have policies allowing public access. This is crucial for analyzing potential security risks and ensuring that sensitive data is not exposed to unauthorized users.

```sql+postgres
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

```sql+sqlite
select
  g.name,
  substr(
    json_extract(s.value, '$.role'),
    instr(json_extract(s.value, '$.role'), '/') + 1
  ) as role,
  e.value as entity
from
  gcp_service_account g,
  json_each(json_extract(g.iam_policy, '$.bindings')) as s,
  json_each(json_extract(s.value, '$.members')) as e
where
  e.value = 'allUsers'
  or e.value = 'allAuthenticatedUsers';
```