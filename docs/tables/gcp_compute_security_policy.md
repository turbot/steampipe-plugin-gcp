---
title: "Steampipe Table: gcp_compute_security_policy - Query Google Cloud Armor Security Policies using SQL"
description: "Allows users to query Google Cloud Armor Security Policies, providing insights into policy details, rules, and configuration."
folder: "Compute"
---

# Table: gcp_compute_security_policy

Google Cloud Armor Security Policies protect your applications from DDoS and application attacks. This table lets you query policy details, rules, and configuration.

## Table Usage Guide

The `gcp_compute_security_policy` table provides insights into Cloud Armor security policies for your GCP projects.

### Examples

#### List all security policies
```sql
select
  name,
  id,
  description,
  self_link
from
  gcp_compute_security_policy;
```

#### Get a security policy by name
```sql
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

#### List all rules for each security policy
```sql
select
  name,
  rules
from
  gcp_compute_security_policy;
```

#### Find security policies with a specific label
```sql
select
  name,
  labels
from
  gcp_compute_security_policy
where
  labels ->> 'env' = 'prod';
```

#### Show all policies with adaptive protection enabled
```sql
select
  name,
  adaptive_protection_config
from
  gcp_compute_security_policy
where
  adaptive_protection_config -> 'layer7DdosDefenseConfig' ->> 'enable' = 'true';
```
