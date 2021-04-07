# Table:  gcp_monitoring_alert_policy

An alerting policy describes a set of conditions that you want to monitor. These conditions might relate to the state of an unhealthy system or to resource consumption.

## Examples

### Basic info

```sql
select
  display_name,
  name,
  enabled,
  documentation ->> 'content' as doc_content,
  tags
from
  gcp_monitoring_alert_policy;
```


### Get the creation record for each alert policy

```sql
select
  display_name,
  name,
  creation_record ->> 'mutateTime' as mutation_time,
  creation_record ->> 'mutatedBy' as mutated_by,
from
  gcp_monitoring_alert_policy;
```


### Get the condition details for each alert policy

```sql
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