# Table: gcp_compute_vpn_tunnel

A VPN tunnel connects two VPN gateways and serves as a virtual medium through which encrypted traffic is passed.

### VPN tunnel basic info

```sql
select
  name,
  id,
  peer_ip,
  shared_secret_hash
from
  gcp_compute_vpn_tunnel;
```


### Get VPN gateway peer details

```sql
select
  name,
  peer_ip,
  target_vpn_gateway
from
  gcp_compute_vpn_tunnel;
```
