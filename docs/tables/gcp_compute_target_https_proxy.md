# Table: gcp_compute_target_https_proxy

A VPN tunnel connects two VPN gateways and serves as a virtual medium through which encrypted traffic is passed.

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

### Get SSL policy details associated with the Target Https Proxy

```sql
select
  name,
  id,
  ssl_policy
from
  gcp_compute_target_https_proxy;
```
