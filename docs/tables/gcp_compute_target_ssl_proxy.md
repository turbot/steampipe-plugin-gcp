---
title: "Steampipe Table: gcp_compute_target_ssl_proxy - Query GCP Compute Engine Target SSL Proxies using SQL"
description: "Allows users to query Target SSL Proxies in GCP Compute Engine, providing insights into the configuration and status of these proxies."
folder: "Compute"
---

# Table: gcp_compute_target_ssl_proxy - Query GCP Compute Engine Target SSL Proxies using SQL

A Target SSL Proxy is a component of Google Cloud Platform's Compute Engine that is used to forward SSL requests to a backend service. It is associated with SSL certificates and provides SSL termination for HTTPS load balancers. This resource is crucial for managing and controlling the traffic to your backend services.

## Table Usage Guide

The `gcp_compute_target_ssl_proxy` table provides valuable insights into the Target SSL Proxies within Google Cloud Platform's Compute Engine. As a network engineer, you can use this table to explore the details of each proxy, including its associated SSL certificates, the backend service it is linked to, and its current status. Utilize it to monitor and manage the traffic flow to your backend services, ensuring optimal performance and security.

## Examples

### Basic info
Explore the basic details of your SSL proxies in Google Cloud Platform's Compute Engine service. This can help you manage and organize your proxies for efficient network traffic control.

```sql+postgres
select
  name,
  id,
  self_link
from
  gcp_compute_target_ssl_proxy;
```

```sql+sqlite
select
  name,
  id,
  self_link
from
  gcp_compute_target_ssl_proxy;
```

### Get SSL policy details for each target SSL proxy
Explore the security configurations of your network by identifying the SSL policies applied to each target SSL proxy within your Google Cloud Platform compute environment. This can help in maintaining the security standards and compliance of your network infrastructure.

```sql+postgres
select
  name,
  id,
  ssl_policy
from
  gcp_compute_target_ssl_proxy;
```

```sql+sqlite
select
  name,
  id,
  ssl_policy
from
  gcp_compute_target_ssl_proxy;
```