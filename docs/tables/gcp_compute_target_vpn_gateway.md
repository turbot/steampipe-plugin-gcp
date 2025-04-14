---
title: "Steampipe Table: gcp_compute_target_vpn_gateway - Query GCP Compute Engine Target VPN Gateways using SQL"
description: "Allows users to query Target VPN Gateways in GCP Compute Engine, enabling insights into the VPN gateways that funnel traffic between your virtual network and your on-premises network."
folder: "Compute"
---

# Table: gcp_compute_target_vpn_gateway - Query GCP Compute Engine Target VPN Gateways using SQL

A Target VPN Gateway in GCP Compute Engine is a virtual router that manages VPN tunnels, providing a way to securely connect networks. It handles traffic between your virtual network and your on-premises network, acting as a focal point for multiple VPN tunnels. This resource is critical for creating secure connections between GCP and your on-premises network.

## Table Usage Guide

The `gcp_compute_target_vpn_gateway` table provides insights into Target VPN Gateways within Google Cloud Platform's Compute Engine. As a network administrator, you can explore gateway-specific details through this table, such as the associated network, the number of tunnels, and the creation timestamp. Use this table to maintain an overview of your VPN connections, monitor the status of each gateway, and ensure secure network communications.

## Examples

### Target VPN gateway basic info
Determine the basic information about your target VPN gateway in Google Cloud Platform. This can be useful for understanding the structure and configuration of your virtual private network.

```sql+postgres
select
  name,
  id,
  self_link,
  kind
from
  gcp_compute_target_vpn_gateway;
```

```sql+sqlite
select
  name,
  id,
  self_link,
  kind
from
  gcp_compute_target_vpn_gateway;
```

### List of all tunnels connected with the gateway
This example helps you identify all the tunnels that are connected to a specific VPN gateway. It is useful when managing network connections and ensuring secure data transmission between different network segments.

```sql+postgres
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

```sql+sqlite
select
  gateway.name as vpn_gateway_name,
  tunnel.peer_ip,
  tunnel.name as tunnel_name
from
  gcp_compute_target_vpn_gateway as gateway,
  json_each(gateway.tunnels) as t
  join gcp_compute_vpn_tunnel as tunnel
  on t.value = tunnel.self_link;
```