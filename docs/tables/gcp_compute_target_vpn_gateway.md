# Table: gcp_compute_target_vpn_gateway

A virtual private network lets you securely connect your Google Compute Engine resources to your own private network.

### Target VPN gateway basic info

```sql
select
  name,
  id,
  self_link,
  kind
from
  gcp_compute_target_vpn_gateway;
```


### List of all tunnels connected with the gateway

```sql
select
  gateway.name as vpn_gateway_name,
  tunnel.peer_ip,
  tunnel.name as tunnel_name
from 
  gcp_compute_target_vpn_gateway as gateway,
  jsonb_array_elements_text(tunnels) as t
  join gcp_compute_vpn_tunnel as tunnel
  on t = tunnel.self_link;
```
