---
title: "Steampipe Table: gcp_iam_role - Query Google Cloud IAM Roles using SQL"
description: "Allows users to query IAM Roles in Google Cloud, specifically the role name, title, description, permissions, and stage, providing insights into role-specific details and permissions."
folder: "IAM"
---

# Table: gcp_iam_role - Query Google Cloud IAM Roles using SQL

Google Cloud IAM (Identity and Access Management) is a service that allows you to manage access control by defining who (identity) has what access (role) for which resource. It provides unified view into security policy across your entire organization, with built-in auditing to ease compliance processes. IAM Roles are a collection of permissions that you can grant to the identities interacting with your Google Cloud resources.

## Table Usage Guide

The `gcp_iam_role` table provides insights into IAM Roles within Google Cloud IAM. As a DevOps engineer, explore role-specific details through this table, including permissions, role description, and associated metadata. Utilize it to uncover information about roles, such as those with wildcard permissions, the role titles, and the verification of role descriptions.

## Examples

### IAM role basic info
Explore the basic information about IAM roles in your GCP environment to understand their configuration and status. This can help in managing access control and ensuring security compliance.

```sql+postgres
select
  name,
  role_id,
  deleted,
  description,
  title
from
  gcp_iam_role;
```

```sql+sqlite
select
  name,
  role_id,
  deleted,
  description,
  title
from
  gcp_iam_role;
```

### List of IAM roles which are in BETA stage
Discover the segments that are still in the BETA stage within the IAM roles. This can be useful to assess the elements within your infrastructure that might need additional testing or development.

```sql+postgres
select
  name,
  description,
  stage
from
  gcp_iam_role
where
  stage = 'BETA';
```

```sql+sqlite
select
  name,
  description,
  stage
from
  gcp_iam_role
where
  stage = 'BETA';
```

### List of IAM customer managed roles
Discover the custom roles within your IAM configuration that are not managed by Google Cloud Platform. This can help in understanding the distribution of responsibilities and access controls within your organization.

```sql+postgres
select
  name,
  description,
  stage
from
  gcp_iam_role
where
  is_gcp_managed = false;
```

```sql+sqlite
select
  name,
  description,
  stage
from
  gcp_iam_role
where
  is_gcp_managed = 0;
```