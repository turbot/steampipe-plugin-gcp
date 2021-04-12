# Table: gcp_compute_target_https_proxy

Target proxies are referenced by one or more forwarding rules. In the case of external HTTP(S) load balancers and internal HTTP(S) load balancers, proxies route incoming requests to a URL map. In the case of SSL proxy load balancers and TCP proxy load balancers, target proxies route incoming requests directly to backend services.

## Examples

### Basic info

```sql
select
  name,
  id,
  self_link,
  proxy_bind
from
  gcp_compute_target_https_proxy;
```

### Get SSL policy details for each target HTTPS proxy

```sql
select
  name,
  id,
  ssl_policy
from
  gcp_compute_target_https_proxy;
```
