---
title: "Steampipe Table: gcp_audit_log_config - Query GCP Audit Log Configurations using SQL"
description: "Allows users to query GCP Audit Log Configurations, providing insights into audit logging settings for Google Cloud services."
folder: "Audit"
---

# Table: gcp_audit_log_config - Query GCP Audit Log Configurations using SQL

Google Cloud Audit Logs maintains audit logs for each Google Cloud project, folder, and organization. These logs help you answer the question of "who did what, where, and when?" within your Google Cloud environment. The audit log configuration determines which services are enabled for audit logging and what types of operations are logged.

## Table Usage Guide

The `gcp_audit_log_config` table provides insights into audit log configurations within Google Cloud Platform. As a security analyst or cloud administrator, explore configuration-specific details through this table, including enabled services, log types, and exempted members. Utilize it to understand your audit logging coverage, verify logging settings, and ensure compliance with your organization's auditing requirements.

## Examples

### Basic info
Explore which services have audit logging enabled and what types of logs are being collected. This helps in understanding your audit logging coverage and identifying potential gaps in logging configuration.

```sql+postgres
select
  service,
  log_type,
  exempted_members
from
  gcp_audit_log_config;
```

```sql+sqlite
select
  service,
  log_type,
  exempted_members
from
  gcp_audit_log_config;
```

### List services with DATA_WRITE logging enabled
Identify which services have data write operations logging enabled. This is useful for ensuring critical data modifications are being tracked appropriately.

```sql+postgres
select
  service,
  log_type,
  exempted_members
from
  gcp_audit_log_config
where
  log_type = 'DATA_WRITE';
```

```sql+sqlite
select
  service,
  log_type,
  exempted_members
from
  gcp_audit_log_config
where
  log_type = 'DATA_WRITE';
```

### Find services with exempted members
Discover which services have specific members exempted from audit logging. This helps identify potential security gaps where actions might not be logged.

```sql+postgres
select
  service,
  log_type,
  exempted_members
from
  gcp_audit_log_config
where
  exempted_members is not null
  and jsonb_array_length(exempted_members) > 0;
```

```sql+sqlite
select
  service,
  log_type,
  exempted_members
from
  gcp_audit_log_config
where
  exempted_members is not null
  and json_array_length(exempted_members) > 0;
```

### Group services by log type
Analyze the distribution of audit log types across services to understand your overall audit logging strategy.

```sql+postgres
select
  log_type,
  count(*) as service_count,
  array_agg(service) as services
from
  gcp_audit_log_config
group by
  log_type;
```

```sql+sqlite
select
  log_type,
  count(*) as service_count,
  group_concat(service) as services
from
  gcp_audit_log_config
group by
  log_type;
``` 