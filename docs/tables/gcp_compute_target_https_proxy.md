---
title: "Steampipe Table: gcp_compute_target_https_proxy - Query GCP Compute Engine Target HTTPS Proxies using SQL"
description: "Allows users to query GCP Compute Engine Target HTTPS Proxies, specifically their configurations and associated SSL certificates, providing insights into the HTTPS traffic routing and SSL/TLS settings."
folder: "Compute"
---

# Table: gcp_compute_target_https_proxy - Query GCP Compute Engine Target HTTPS Proxies using SQL

A Target HTTPS Proxy is a component of GCP Compute Engine used for directing incoming HTTPS traffic to a URL map. It is associated with one or more SSL certificates for secure connections and a URL map that routes the traffic. Target HTTPS Proxies are regional resources and are required by external HTTP(S) load balancers.

## Table Usage Guide

The `gcp_compute_target_https_proxy` table provides insights into Target HTTPS Proxies within Google Cloud Compute Engine. As a network engineer, explore proxy-specific details through this table, including associated SSL certificates, URL maps, and regional settings. Utilize it to uncover information about proxies, such as their configurations, the SSL/TLS settings, and the routing of HTTPS traffic.

## Examples

### Basic info
Explore the configuration of your HTTPS proxy settings in Google Cloud Platform (GCP) to identify potential issues or areas for improvement. This query will assist you in gaining insights into the binding status of your proxies, enhancing your network's performance and security.

```sql+postgres
select
  name,
  id,
  self_link,
  proxy_bind
from
  gcp_compute_target_https_proxy;
```

```sql+sqlite
select
  name,
  id,
  self_link,
  proxy_bind
from
  gcp_compute_target_https_proxy;
```

### Get SSL policy details for each target HTTPS proxy
Analyze the settings to understand the SSL policy associated with each target HTTPS proxy in your Google Cloud Compute environment. This could be beneficial for maintaining security standards and protocols across your digital assets.

```sql+postgres
select
  name,
  id,
  ssl_policy
from
  gcp_compute_target_https_proxy;
```

```sql+sqlite
select
  name,
  id,
  ssl_policy
from
  gcp_compute_target_https_proxy;
```