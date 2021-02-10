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
  vpn_gateway_name
from
  gcp_compute_vpn_tunnel;
```

### Traffic selector info of each tunnel

```sql
select
  name,
  jsonb_array_elements_text(local_traffic_selector) as local_traffic_selector,
  jsonb_array_elements_text(remote_traffic_selector) as remote_traffic_selector
from
  gcp_compute_vpn_tunnel;
```
