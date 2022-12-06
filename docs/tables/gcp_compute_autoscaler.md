# Table: gcp_compute_autoscaler

Autoscalers automatically add or delete instances from a managed instance group according to your defined autoscaling policy.

## Examples

### Basic Info

```sql
select 
  name,
  description,
  self_link, 
  status,
  location, 
  akas,
  recommended_size
from 
  gcp_compute_autoscaler;
```

### Get auto scaling policy for autoscalers

```sql
select
  title,
  autoscaling_policy ->> 'mode' as mode,
  autoscaling_policy -> 'cpuUtilization' ->> 'predictiveMethod' as cpu_utilization_method,
  autoscaling_policy -> 'cpuUtilization' ->> 'utilizationTarget' as cpu_utilization_target,
  autoscaling_policy ->> 'maxNumReplicas' as max_replicas,
  autoscaling_policy ->> 'minNumReplicas' as min_replicas,
  autoscaling_policy ->> 'coolDownPeriodSec' as cool_down_period_sec
from
  gcp_compute_autoscaler;
```

### Get autoscalers with configuration errors

```sql
select 
  name,
  description,
  self_link, 
  status,
  location, 
  akas,
  recommended_size
from 
  gcp_compute_autoscaler
where 
  status = 'ERROR';
```

### Get instance groups having autoscaling enabled

```sql
select
  a.title as autoscaler_name,
  g.name as instance_group_name,
  g.description as instance_group_description,
  g.size as instance_group_size
from
  gcp_compute_instance_group g,
  gcp_compute_autoscaler a
where
  g.name = split_part(a.target, 'instanceGroupManagers/', 2);
```