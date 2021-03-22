# Table: gcp_compute_target_ssl_proxy

SSL Proxy Load Balancing terminates SSL connections from the client and creates new connections to the backends.

## Examples

### Basic info

```sql
select
  name,
  id,
  self_link
from
  gcp_compute_target_ssl_proxy;
```

### Get SSL policy details

```sql
select
  name,
  id,
  ssl_policy
from
  gcp_compute_target_ssl_proxy;
```
