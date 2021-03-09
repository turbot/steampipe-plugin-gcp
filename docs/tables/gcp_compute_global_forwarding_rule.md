# Table: gcp_compute_global_forwarding_rule

Compute global forwarding rule and its corresponding IP address represent the frontend configuration of a Google Cloud load balancer. A forwarding rule specifies a backend service, target proxy, or target pool. A forwarding rule and its IP address are internal or external.

### Basic info

```sql
select
  name,
  id,
  ip_address,
  ip_protocol,
  port_range,
  target
from
  gcp_compute_global_forwarding_rule;
```


### List global forwarding rules which are globally accessible

```sql
select
  name,
  id,
  ip_address,
  allow_global_access
from
  gcp_compute_global_forwarding_rule
where
  allow_global_access;
```


### List global forwarding rules where mirroring collector is enabled (i.e load balancer can be used as a collector for packet mirroring)

```sql
select
  name,
  id,
  is_mirroring_collector
from
  gcp_compute_global_forwarding_rule
where
  is_mirroring_collector;
```