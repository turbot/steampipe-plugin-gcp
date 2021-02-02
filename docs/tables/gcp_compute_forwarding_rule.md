# Table: gcp_compute_forwarding_rule

Compute forwarding rule and its corresponding IP address represent the frontend configuration of a Google Cloud load balancer. A forwarding rule specifies a backend service, target proxy, or target pool. A forwarding rule and its IP address are internal or external.

### Basic info

```sql
select
  name,
  id,
  self_link,
  backend_service,
  ip_address,
  load_balancing_scheme
from
  gcp_compute_forwarding_rule;
```


### List of forwarding rules which are globally accessible

```sql
select
  name,
  id,
  allow_global_access
from
  gcp_compute_forwarding_rule
where
  not allow_global_access;
```


### List of EXTERNAL forwarding rules

```sql
select
  name,
  id,
  load_balancing_scheme
from
  gcp_compute_forwarding_rule
where
  load_balancing_scheme = 'EXTERNAL';
```