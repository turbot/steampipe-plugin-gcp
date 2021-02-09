# Table: gcp_compute_target_pool

The Target Pools resource defines a group of instances that should receive incoming traffic from forwarding rules. When a forwarding rule directs traffic to a target pool, Google Compute Engine picks an instance from these target pools based on a hash of the source IP and port and the destination IP and port.

## Examples

### Basic info

```sql
select
  name,
  id,
  location
from
  gcp_compute_target_pool;
```

### List of target pools and attached instances that receives incoming traffic

```sql
select
  name,
  id,
  split_part(i, '/', 11) as instance_name
from
  gcp_compute_target_pool,
  jsonb_array_elements_text(instances) as i;
```

### List of Health checks attached to each target pool

```sql
select
  name,
  id,
  split_part(h, '/', 10) as health_check
from
  gcp_compute_target_pool,
  jsonb_array_elements_text(health_checks) as h;
```
