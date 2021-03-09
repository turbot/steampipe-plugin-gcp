# Table: gcp_compute_instance_group

An instance group is a collection of virtual machine (VM) instances that you can manage as a single entity.

## Examples

### Basic info

```sql
select
  name,
  kind,
  size,
  self_link
from
  gcp_compute_instance_group;
```

### List of managed instance group

```sql
select
  name,
  kind,
  location,
  self_link
from
  gcp_compute_instance_group
where
  is_managed;
```

### List of instance group having no instances in it

```sql
select
  name,
  kind,
  location,
  self_link,
  size
from
  gcp_compute_instance_group
where
  size = 0;
```

### List of instances resides in this group

```sql
select
  grp.name as instance_group_name,
  instance.name as instance_name,
  instance.location as instance_location,
  instance.status as instance_state
from
  gcp_compute_instance_group as grp,
  gcp_compute_instance as instance,
  jsonb_array_elements(grp.instances) as ins
where
  ins ->> 'instance' = instance.self_link;
```
