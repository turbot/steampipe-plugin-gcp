---
title: "Steampipe Table: gcp_compute_security_policy - Query Google Cloud Armor Security Policies using SQL"
description: "Allows users to query Google Cloud Armor Security Policies, providing insights into policy details, rules, and configuration."
folder: "Compute"
---

# Table: gcp_compute_security_policy

Google Cloud Armor Security Policies protect your applications from DDoS and application attacks. This table lets you query policy details, rules, and configuration.

## Table Usage Guide

The `gcp_compute_security_policy` table provides insights into Cloud Armor security policies for your GCP projects.

**Important Notes:**
- This table supports optional quals. Queries with optional quals are optimised to use security policy filters. Optional quals are supported for the following columns:
  - `id`
  - `type`
  - `description`
  - `self_link`
  - `filter`: For additional details regarding the filter string, please refer to the documentation at https://cloud.google.com/compute/docs/reference/rest/v1/securityPolicies/list?filter#query-parameters.

## Examples

### List all security policies
List all Google Cloud Armor security policies in your GCP projects, including their names, IDs, descriptions, and self links.

```sql+postgres
select
  name,
  id,
  description,
  self_link
from
  gcp_compute_security_policy;
```

```sql+sqlite
select
  name,
  id,
  description,
  self_link
from
  gcp_compute_security_policy;
```

### Get a security policy by name
Retrieve details for a specific security policy by its name, including its rules, labels, and associated project.

```sql+postgres
select
  name,
  id,
  description,
  rules,
  labels,
  project
from
  gcp_compute_security_policy
where
  name = 'my-security-policy';
```

```sql+sqlite
select
  name,
  id,
  description,
  rules,
  labels,
  project
from
  gcp_compute_security_policy
where
  name = 'my-security-policy';
```

### List all rules for each security policy
Show the rules defined for each security policy, helping you review policy configurations across your environment.

```sql+postgres
select
  name,
  rules
from
  gcp_compute_security_policy;
```

```sql+sqlite
select
  name,
  rules
from
  gcp_compute_security_policy;
```

### Show all policies with adaptive protection enabled
Identify security policies that have adaptive protection enabled, which helps protect against advanced DDoS attacks.

```sql+postgres
select
  name,
  adaptive_protection_config
from
  gcp_compute_security_policy
where
  adaptive_protection_config -> 'layer7DdosDefenseConfig' ->> 'enable' = 'true';
```

```sql+sqlite
select
  name,
  adaptive_protection_config
from
  gcp_compute_security_policy
where
  json_extract(adaptive_protection_config, '$.layer7DdosDefenseConfig.enable') = 'true';
```

### Filter security policies by ID and description 
Use the filter parameter to query security policies with specific criteria. This example shows how to filter by ID and description.

```sql+postgres
select
  name,
  id,
  creation_timestamp,
  description,
  self_link
from
  gcp_compute_security_policy
where
  filter = 'id = 4811866613213140474 AND description = "Default security policy for: tet5s"';
```

```sql+sqlite
select
  name,
  id,
  creation_timestamp,
  description,
  self_link
from
  gcp_compute_security_policy
where
  filter = 'id = 4811866613213140474 AND description = "Default security policy for: tet5s"';
```
