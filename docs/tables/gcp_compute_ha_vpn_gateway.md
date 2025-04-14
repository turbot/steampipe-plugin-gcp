---
title: "Steampipe Table: gcp_compute_ha_vpn_gateway - Query GCP Compute HA VPN Gateways using SQL"
description: "Allows users to query HA VPN Gateways in GCP Compute, specifically providing details about the gateways such as their ID, name, creation timestamp, and associated network details."
folder: "Compute"
---

# Table: gcp_compute_ha_vpn_gateway - Query GCP Compute HA VPN Gateways using SQL

Highly Available (HA) VPN Gateway is a resource within Google Cloud's GCP Compute service. It provides an entry point for connecting your Virtual Private Cloud (VPC) network to your on-premises network or another VPC network through an IPsec VPN tunnel. The HA VPN gateways are designed to provide high reliability, high throughput, and low latency.

## Table Usage Guide

The `gcp_compute_ha_vpn_gateway` table provides insights into HA VPN Gateways within Google Cloud's GCP Compute service. As a network administrator, explore gateway-specific details through this table, including their ID, name, creation timestamp, and associated network details. Utilize it to uncover information about gateways, such as their status, interfaces, and the regions they are located in.

## Examples

### Basic info
Explore which High Availability (HA) VPN gateways are present in your Google Cloud Platform (GCP) setup. This can be beneficial in managing network connectivity and ensuring secure data transmission.

```sql+postgres
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

```sql+sqlite
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
Explore the VPN gateways within your Google Cloud Platform to identify each VPN interface's unique ID and IP address. This can be particularly useful for network management and troubleshooting connectivity issues.

```sql+postgres
select
  name as vpn_gateway_name,
  i ->> 'id' as vpn_interface_id,
  i ->> 'ipAddress' as vpn_interface_ip_address
from
  gcp_compute_ha_vpn_gateway g,
  jsonb_array_elements(vpn_interfaces) i;
```

```sql+sqlite
select
  name as vpn_gateway_name,
  json_extract(i.value, '$.id') as vpn_interface_id,
  json_extract(i.value, '$.ipAddress') as vpn_interface_ip_address
from
  gcp_compute_ha_vpn_gateway g,
  json_each(vpn_interfaces) as i;
```

### Get network detail per VPN gateway
Explore the specific details of each VPN gateway in your network, such as its creation timestamp, routing mode, location, and associated project. This can aid in understanding your network's configuration and how each VPN gateway contributes to the overall network structure.

```sql+postgres
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

```sql+sqlite
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