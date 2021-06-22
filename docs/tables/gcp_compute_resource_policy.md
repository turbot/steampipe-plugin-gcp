# Table: gcp_compute_resource_policy

A policy that can be attached to a resource to specify or schedule actions on that resource.

## Examples

### Basic info

```sql
select
  name,
  status,
  self_link
from
  gcp_compute_resource_policy;
```


### List of policy used to schedule an instance

```sql
select
  p.name as policy_name,
  i.name,
  p.instance_schedule_policy
from
  gcp_compute_resource_policy as p
  join gcp_compute_instance as i on i.resource_policies ?| array[p.self_link]
where
  p.instance_schedule_policy is not null;
```


### List invalid policy

```sql
select
  name,
  self_link,
  status
from
  gcp_compute_resource_policy
where
  status = 'INVALID';
```
