---
title: "Steampipe Table: gcp_audit_policy - Query Google Cloud Platform Audit Policies using SQL"
description: "Allows users to query Audit Policies in Google Cloud Platform, specifically the settings and configurations of all audit policies, providing insights into compliance and security posture."
folder: "Audit Policy"
---

# Table: gcp_audit_policy - Query Google Cloud Platform Audit Policies using SQL

Google Cloud Audit Logs is a feature that maintains three audit logs for each Google Cloud project, folder, and organization: Admin Activity, Data Access, and System Event. These logs can be used to help you answer the question of "who did what, where, and when?" within your Google Cloud environment. Audit logs are critical for incident response, forensics, and establishing regulatory and compliance controls.

## Table Usage Guide

The `gcp_audit_policy` table provides insights into audit policies within Google Cloud Platform. As a security analyst, explore policy-specific details through this table, including policy settings, service conditions, and associated metadata. Utilize it to uncover information about policies, such as those with specific service conditions, the identity of the creator and the verification of policy settings.

## Examples

### Basic info
Determine the areas in which different types of logs are created by analyzing the audit policies within the Google Cloud Platform. This is useful for managing and understanding the audit trails in your environment.

```sql+postgres
select
  service,
  jsonb_array_elements(audit_log_configs) ->> 'logType' as log_type
from
  gcp_audit_policy;
```

```sql+sqlite
select
  service,
  json_extract(audit_log_configs, '$.logType') as log_type
from
  gcp_audit_policy,
  json_each(audit_log_configs);
```


### List of services which has data write access
Determine the areas in which certain services have data write access. This is useful for understanding potential security risks and ensuring only appropriate services have this level of access.

```sql+postgres
select
  service,
  log_type ->> 'logType' as log_type
from
  gcp_audit_policy,
  jsonb_array_elements(audit_log_configs) as log_type
where
  log_type ->> 'logType' = 'DATA_WRITE';
```

```sql+sqlite
select
  service,
  json_extract(log_type.value, '$.logType') as log_type
from
  gcp_audit_policy,
  json_each(audit_log_configs) as log_type
where
  json_extract(log_type.value, '$.logType') = 'DATA_WRITE';
```