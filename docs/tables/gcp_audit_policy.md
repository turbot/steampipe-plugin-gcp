# Table: gcp_audit_policy

Audit policy defines rules for which events are recorded as log entries.

### Basic info

```sql
select
  service,
  jsonb_array_elements(audit_log_configs) ->> 'logType' as log_type
from
  gcp_audit_policy;
```


### List of services which has data write access

```sql
select
  service,
  log_type ->> 'logType' as log_type
from
  gcp_audit_policy,
  jsonb_array_elements(audit_log_configs) as log_type
where
  log_type ->> 'logType' = 'DATA_WRITE';
```