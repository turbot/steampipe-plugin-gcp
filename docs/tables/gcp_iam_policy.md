---
title: "Steampipe Table: gcp_iam_policy - Query Google Cloud IAM Policies using SQL"
description: "Allows users to query IAM Policies in Google Cloud, specifically the policy bindings and members, providing insights into access control and permissions."
folder: "IAM"
---

# Table: gcp_iam_policy - Query Google Cloud IAM Policies using SQL

Google Cloud Identity and Access Management (IAM) provides the right tools to manage resource permissions with minimum fuss and high automation. It offers unified control across the entire suite of Google Cloud resources. IAM Policies are the primary resources in IAM that bind a set of members to a role, thus defining what actions the members can take on the resources.

## Table Usage Guide

The `gcp_iam_policy` table provides insights into IAM Policies within Google Cloud Identity and Access Management (IAM). As a Security Analyst, explore policy-specific details through this table, including bindings, roles, and associated members. Utilize it to uncover information about the policies, such as those with broad access, the binding of members to roles, and the verification of permissions.

## Examples

### List of project members with their roles
Explore which roles are assigned to different project members. This can help in managing access control and ensuring appropriate permissions are allocated.

```sql+postgres
select
  entity,
  p ->> 'role' as role
from
  gcp_iam_policy,
  jsonb_array_elements(bindings) as p,
  jsonb_array_elements_text(p -> 'members') as entity;
```

```sql+sqlite
select
  e.value as entity,
  json_extract(p.value, '$.role') as role
from
  gcp_iam_policy,
  json_each(bindings) as p,
  json_each(json_extract(p.value, '$.members')) as e;
```

### List of members with owner roles
Explore which members have been assigned the 'owner' role in your Google Cloud Platform IAM policy. This is useful for gaining insights into access control and ensuring appropriate permissions are in place.

```sql+postgres
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

```sql+sqlite
Error: SQLite does not support split functions.
```