# Table: gcp_compute_ha_vpn_gateway

Represents a VPN gateway running in GCP. This virtual device is managed by Google, but used only by you. This type of VPN Gateway allows for the creation of VPN solutions with higher availability than classic Target VPN Gateways.

### Basic info

```sql
select
  name,
  id,
  description,
  location,
  self_link,
  kind
from
  gcp_compute_ha_vpn_gateway;
```

### List VPN interfaces for all VPN gateways

```sql
select
  name as vpn_gateway_name,
  i ->> 'id' as vpn_interface_id,
  i ->> 'ipAddress' as vpn_interface_ip_address
from
  gcp_compute_ha_vpn_gateway g,
  jsonb_array_elements(vpn_interfaces) i;
```

### Get network detail per VPN gateway

```sql
select
  n.name,
  n.id,
  n.creation_timestamp,
  n.mtu,
  n.routing_mode,
  n.location,
  n.project
from
  gcp_compute_ha_vpn_gateway g,
  gcp_compute_network n
where
  g.network = n.self_link;
```