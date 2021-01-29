# Table: gcp_compute_instance

A Virtual Private Cloud (VPC) network is a virtual version of a physical network, implemented inside of Google's production network.

## Examples

### List networks having auto create subnetworks feature disabled

```sql
select
  name,
  id,
  auto_create_subnetworks
from
  gcp_compute_network
where
  not auto_create_subnetworks;
```

### List networks having routing_mode set to REGIONAL

```sql
select
  name,
  id,
  routing_mode
from
  gcp_compute_network
where
  routing_mode = 'REGIONAL';
```

### List the subnets stats to the network

```sql
select
  name,
  count(d) as num_subnets
from
  gcp_compute_network as i,
  jsonb_array_elements(subnetworks) as d
group by
    name;
```
