---
title: "Steampipe Table: gcp_secret_manager_secret - Query Google Cloud Platform Secret Manager Secrets using SQL"
description: "Allows users to query Secret Manager secrets in Google Cloud Platform, providing details about the secrets stored within Secret Manager."
folder: "Secret Manager"
---

# Table: gcp_secret_manager_secret - Query Google Cloud Platform Secret Manager Secrets using SQL

A Secret Manager Secret in Google Cloud Platform is a secure place to store and manage sensitive information, such as API keys, passwords, certificates, and other sensitive data. Secret Manager provides a central place and single source of truth to manage, access, and audit secrets across Google Cloud.

## Table Usage Guide

The `gcp_secret_manager_secret` table provides insights into secrets stored within the Google Cloud Secret Manager. As a security engineer, you can explore secret-specific details through this table, including the associated project, creation time, expiration time, and other metadata. Utilize it to understand the distribution and lifecycle of secrets for better management and security.

## Examples

### List all secrets in a specific project
Identify all the secrets stored within a specific Google Cloud project. This is useful for auditing and managing secrets within your project.

```sql+postgres
select
  name,
  project,
  create_time,
  expire_time
from
  gcp_secret_manager_secret
where
  project = 'my-gcp-project';
```

```sql+sqlite
select
  name,
  project,
  create_time,
  expire_time
from
  gcp_secret_manager_secret
where
  project = 'my-gcp-project';
```

### Find secrets that are about to expire
Identify secrets that are nearing their expiration date. This is useful for proactively managing and rotating secrets to maintain security.

```sql+postgres
select
  name,
  project,
  expire_time
from
  gcp_secret_manager_secret
where
  expire_time < now() + interval '30 days';
```

```sql+sqlite
select
  name,
  project,
  expire_time
from
  gcp_secret_manager_secret
where
  expire_time < datetime('now', '+30 days');
```

### Get details of a specific secret
Retrieve detailed information about a specific secret, including its labels, annotations, and replication policy.

```sql+postgres
select
  name,
  labels,
  annotations,
  replication,
  ttl
from
  gcp_secret_manager_secret
where
  name = 'my-secret';
```

```sql+sqlite
select
  name,
  labels,
  annotations,
  replication,
  ttl
from
  gcp_secret_manager_secret
where
  name = 'my-secret';
```

### Get user managed replication details of secrets
Retrieve replication details about the secrets.

```sql+postgres
select
  name,
  create_time,
  replication -> 'userManaged' -> 'replicas' as user_managed_replicas
from
  gcp_secret_manager_secret;
```

```sql+sqlite
select
  name,
  create_time,
  json_extract(replication, '$.userManaged.replicas') as user_managed_replicas
from
  gcp_secret_manager_secret;
```