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

### Get SSL policy details for each target SSL proxy

```sql
select
  name,
  id,
  ssl_policy
from
  gcp_compute_target_ssl_proxy;
```

### Get SSL certificates used to authenticate connections to Backends

```sql
select
  name,
  id,
  ssl_certificates
from
  gcp_compute_target_ssl_proxy;
```
