# Table: gcp_compute_forwarding_rule

Compute forwarding rule and its corresponding IP address represent the frontend configuration of a Google Cloud load balancer. A forwarding rule specifies a backend service, target proxy, or target pool. A forwarding rule and its IP address are internal or external.

### Basic info

```sql
select
  name,
  id,
  self_link,
  backend_service,
  load_balancing_scheme,
from
  gcp_compute_forwarding_rule;
```


### Get the labels attached with forwarding rule

```sql
select
  name,
  id,
  labels
from
  gcp_compute_forwarding_rule;
```