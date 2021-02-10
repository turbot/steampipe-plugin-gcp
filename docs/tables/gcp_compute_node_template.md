# Table: gcp_compute_node_template

Node templates specify properties for creating sole-tenant nodes, such as node type, vCPU and memory requirements, node affinity labels, and region.

## Examples

### List of n2-node-80-640 type node templates

```sql
select
  name,
  id,
  location,
  node_type
from
  gcp_compute_node_template
where
  node_type = 'n2-node-80-640';
```


### List of node templates where cpu overcommit is enabled

```sql
select
  name,
  id,
  node_type
from
  gcp_compute_node_template
where
  cpu_overcommit_type = 'ENABLED';
```


### Count of node templates per location

```sql
select
  location,
  count(*)
from
  gcp_compute_node_template
group by
  location;
```


### Find unused node templates

```sql
select
  t.name,
  t.id
from
  gcp_compute_node_template as t
left join
  gcp_compute_node_group as g on g.node_template = t.self_link
where
  g is null;
```
