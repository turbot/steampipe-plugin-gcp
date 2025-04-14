---
title: "Steampipe Table: gcp_compute_vpn_tunnel - Query GCP Compute VPN Tunnels using SQL"
description: "Allows users to query GCP Compute VPN Tunnels, providing details about the VPN Tunnels in Google Cloud Platform's Compute service."
folder: "Compute"
---

# Table: gcp_compute_vpn_tunnel - Query GCP Compute VPN Tunnels using SQL

Google Cloud Platform's VPN Tunnels are part of the Compute service, providing secure encrypted connections between your network and your cloud network. They are used to securely extend your private network into your Google Cloud Virtual Private Cloud network through an IPsec VPN connection. They provide an essential layer of security for your data.

## Table Usage Guide

The `gcp_compute_vpn_tunnel` table provides insights into VPN tunnels within Google Cloud Platform's Compute service. As a network engineer, explore tunnel-specific details through this table, including the associated network, target VPN gateway, and routing configuration. Utilize it to uncover information about tunnels, such as their statuses, the IKE versions used, and the shared secrets for the tunnels.

## Examples

### VPN tunnel basic info
Explore the basic information about VPN tunnels to understand their configuration and security settings. This can be useful for network administrators to assess and manage the VPN infrastructure within their organization.

```sql+postgres
select
  name,
  id,
  peer_ip,
  shared_secret_hash
from
  gcp_compute_vpn_tunnel;
```

```sql+sqlite
select
  name,
  id,
  peer_ip,
  shared_secret_hash
from
  gcp_compute_vpn_tunnel;
```

### Get VPN gateway peer details
Determine the areas in which you can gain insights into the details of your VPN gateway peers. This can be beneficial in understanding the configuration and connectivity of your virtual private network.

```sql+postgres
select
  name,
  peer_ip,
  vpn_gateway_name
from
  gcp_compute_vpn_tunnel;
```

```sql+sqlite
select
  name,
  peer_ip,
  vpn_gateway_name
from
  gcp_compute_vpn_tunnel;
```

### Traffic selector info of each tunnel
This example helps you identify the traffic selectors for each VPN tunnel in your Google Cloud Platform. It's particularly useful for network administrators seeking to understand how traffic is being directed and managed within their VPN infrastructure.

```sql+postgres
select
  name,
  jsonb_array_elements_text(local_traffic_selector) as local_traffic_selector,
  jsonb_array_elements_text(remote_traffic_selector) as remote_traffic_selector
from
  gcp_compute_vpn_tunnel;
```

```sql+sqlite
select
  name,
  json_each.value as local_traffic_selector,
  json_each.value as remote_traffic_selector
from
  gcp_compute_vpn_tunnel,
  json_each(gcp_compute_vpn_tunnel.local_traffic_selector),
  json_each(gcp_compute_vpn_tunnel.remote_traffic_selector);
```