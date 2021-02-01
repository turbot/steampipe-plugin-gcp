# Table:  gcp_compute_route

Google Cloud routes define the paths that network traffic takes from a virtual machine (VM) instance to other destinations. These destinations can be inside Google Cloud Virtual Private Cloud (VPC) network (for example, in another VM) or outside it.

## Examples

### Route basic info

```sql
select
  name,
  dest_range,
  priority,
  network
from
  gcp_compute_route;
```


### List all system-generated default routes

```sql
select
  name,
  dest_range,
  priority,
  next_hop_gateway
from
  gcp_compute_route
where priority = 1000 and dest_range = '0.0.0.0/0';
```
