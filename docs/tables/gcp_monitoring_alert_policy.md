---
title: "Steampipe Table: gcp_monitoring_alert_policy - Query Google Cloud Monitoring Alert Policies using SQL"
description: "Allows users to query Google Cloud Monitoring Alert Policies, specifically providing insights into policy details, conditions, and notification channels."
folder: "Cloud Monitoring"
---

# Table: gcp_monitoring_alert_policy - Query Google Cloud Monitoring Alert Policies using SQL

Google Cloud Monitoring Alert Policies are a part of Google Cloud's operations suite that enables users to define conditions under which alerts are triggered. These alerts notify about specific events occurring in your Google Cloud environment. Alert policies are used to monitor resources, services, and applications running on Google Cloud.

## Table Usage Guide

The `gcp_monitoring_alert_policy` table provides insights into alert policies within Google Cloud Monitoring. As a DevOps engineer or a system administrator, you can explore policy-specific details through this table, including conditions, enabled status, and associated notification channels. Utilize it to uncover information about alert policies, such as those with specific conditions, the notification channels associated with each policy, and the verification of policy status.

## Examples

### Basic info
Explore which monitoring alert policies are enabled in your Google Cloud Platform. This can help you assess the current alerting configuration and discover any potential gaps in your monitoring strategy.

```sql+postgres
select
  display_name,
  name,
  enabled,
  documentation ->> 'content' as doc_content,
  tags
from
  gcp_monitoring_alert_policy;
```

```sql+sqlite
select
  display_name,
  name,
  enabled,
  json_extract(documentation, '$.content') as doc_content,
  tags
from
  gcp_monitoring_alert_policy;
```

### Get the creation record for each alert policy
Discover the segments that show when and by whom each alert policy was last modified. This can be particularly useful for auditing purposes or to track changes in alert policies over time.

```sql+postgres
select
  display_name,
  name,
  creation_record ->> 'mutateTime' as mutation_time,
  creation_record ->> 'mutatedBy' as mutated_by
from
  gcp_monitoring_alert_policy;
```

```sql+sqlite
select
  display_name,
  name,
  json_extract(creation_record, '$.mutateTime') as mutation_time,
  json_extract(creation_record, '$.mutatedBy') as mutated_by
from
  gcp_monitoring_alert_policy;
```

### Get the condition details for each alert policy
Discover the specifics of each alert policy, including the filter details and threshold values. This allows for a comprehensive understanding of the alert triggers, aiding in effective monitoring and management.

```sql+postgres
select
  display_name,
  con ->> 'displayName' as filter_display_name,
  con -> 'conditionThreshold' ->> 'filter' as filter,
  con -> 'conditionThreshold' ->> 'thresholdValue' as threshold_value,
  con -> 'conditionThreshold' ->> 'trigger' as trigger
from
  gcp_monitoring_alert_policy,
  jsonb_array_elements(conditions) as con;
```

```sql+sqlite
select
  display_name,
  json_extract(con.value, '$.displayName') as filter_display_name,
  json_extract(con.value, '$.conditionThreshold.filter') as filter,
  json_extract(con.value, '$.conditionThreshold.thresholdValue') as threshold_value,
  json_extract(con.value, '$.conditionThreshold.trigger') as trigger
from
  gcp_monitoring_alert_policy,
  json_each(conditions) as con;
```