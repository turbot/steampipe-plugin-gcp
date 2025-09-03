---
title: "Steampipe Table: gcp_organization_audit_policy - Query Google Cloud Platform Organization Audit Policies using SQL"
description: "Allows users to query Organization Audit Policies in Google Cloud Platform, specifically the settings and configurations of all audit policies at the organization level, providing insights into compliance and security posture."
folder: "Audit Policy"
---

# Table: gcp_organization_audit_policy - Query Google Cloud Platform Organization Audit Policies using SQL

Google Cloud Audit Logs is a feature that maintains three audit logs for each Google Cloud project, folder, and organization: Admin Activity, Data Access, and System Event. These logs can be used to help you answer the question of "who did what, where, and when?" within your Google Cloud environment. Audit logs are critical for incident response, forensics, and establishing regulatory and compliance controls.

## Table Usage Guide

The `gcp_organization_audit_policy` table provides insights into audit policies within Google Cloud Platform organizations. As a security analyst, explore policy-specific details through this table, including policy settings, service conditions, and associated metadata. Utilize it to uncover information about policies, such as those with specific service conditions, the identity of the creator and the verification of policy settings.

## Examples

### Basic info
Determine the areas in which different types of logs are created by analyzing the audit policies within the Google Cloud Platform organizations. This is useful for managing and understanding the audit trails in your environment.

```sql+postgres
select
  organization_id,
  service,
  jsonb_array_elements(audit_log_configs) ->> 'logType' as log_type
from
  gcp_organization_audit_policy;
```

```sql+sqlite
select
  organization_id,
  service,
  json_extract(audit_log_configs, '$.logType') as log_type
from
  gcp_organization_audit_policy,
  json_each(audit_log_configs);
```

### List of services that have data write access
Determine the areas in which certain services have data write access. This is useful for understanding potential security risks and ensuring only appropriate services have this level of access.

```sql+postgres
select
  organization_id,
  service,
  log_type ->> 'logType' as log_type
from
  gcp_organization_audit_policy,
  jsonb_array_elements(audit_log_configs) as log_type
where
  log_type ->> 'logType' = 'DATA_WRITE';
```

```sql+sqlite
select
  organization_id,
  service,
  json_extract(log_type.value, '$.logType') as log_type
from
  gcp_organization_audit_policy,
  json_each(audit_log_configs) as log_type
where
  json_extract(log_type.value, '$.logType') = 'DATA_WRITE';
```

### Get audit policy for a specific organization
Retrieve audit policy details for a specific organization by its ID.

```sql+postgres
select
  organization_id,
  service,
  audit_log_configs
from
  gcp_organization_audit_policy
where
  organization_id = '123456789';
```

```sql+sqlite
select
  organization_id,
  service,
  audit_log_configs
from
  gcp_organization_audit_policy
where
  organization_id = '123456789';
```

### List organizations with admin activity logging enabled
Find organizations that have admin activity logging enabled for audit purposes.

```sql+postgres
select
  organization_id,
  service,
  log_type ->> 'logType' as log_type
from
  gcp_organization_audit_policy,
  jsonb_array_elements(audit_log_configs) as log_type
where
  log_type ->> 'logType' = 'ADMIN_READ';
```

```sql+sqlite
select
  organization_id,
  service,
  json_extract(log_type.value, '$.logType') as log_type
from
  gcp_organization_audit_policy,
  json_each(audit_log_configs) as log_type
where
  json_extract(log_type.value, '$.logType') = 'ADMIN_READ';
```
