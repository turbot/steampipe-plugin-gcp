# Table: gcp_compute_node_group

Represent a sole-tenant Node Group resource.

A sole-tenant node is a physical server that is dedicated to hosting VM instances only for specific project. Use sole-tenant nodes to keep the instances physically separated from instances in other projects, or to group the instances together on the same host hardware.

## Examples

### Node group basic info

```sql
select
  name,
  status,
  size,
  self_link
from
  gcp_compute_node_group;
```

### List of node groups where the autoscaler is not enabled

```sql
select
  name,
  id,
  status,
  autoscaling_policy_mode
from
  gcp_compute_node_group
where
  autoscaling_policy_mode <> 'ON';
```

### List of node groups with default maintenance settings

```sql
select
  name,
  id,
  status,
  autoscaling_policy_mode
from
  gcp_compute_node_group
where
  maintenance_policy = 'DEFAULT';
```


### List node types for node groups

```sql
select
  g.name,
  g.id,
  g.location,
  t.node_type
from
  gcp_compute_node_group as g,
  gcp_compute_node_template as t
where
  g.node_template = t.self_link;
```
