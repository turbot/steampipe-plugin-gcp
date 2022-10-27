# Table:  gcp_compute_instance_group

An instance group is a collection of virtual machine (VM) instances that you can manage as a single entity.

## Examples

### Basic Info

```sql
select 
  name,
  description,
  self_link, 
  size,
  location, 
  akas, 
  project
from 
  gcp_compute_instance_group;
```

### Get number of instances per instance group
```sql
select
  name,
  size as no_of_instances
from
  gcp_compute_instance_group;
```

### Get instance details of each instance group
```sql
select
  g.name,
  ins.name as instance_name,
  ins.status as instance_status
from
  gcp_compute_instance_group as g,
  jsonb_array_elements(instances) as i,
  gcp_compute_instance as ins
where
  (i ->> 'instance') = ins.self_link;
```

### Get network and subnetwork info of each instance group

```sql
select
  g.name as instance_group_name,
  n.name as network_name,
  s.name as subnetwork_name,
  s.ip_cidr_range,
  s.gateway_address,
  n.location
from
  gcp_compute_instance_group as g,
  gcp_compute_network as n,
  gcp_compute_subnetwork as s
where
  g.network = n.self_link
  and g.subnetwork = s.self_link;
```