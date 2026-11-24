---
title: "Steampipe Table: gcp_compute_target_http_proxy - Query GCP Compute Engine Target HTTP Proxies using SQL"
description: "Allows users to query GCP Compute Engine Target HTTP Proxies, both global and regional, including the URL map they route to."
folder: "Compute"
---

# Table: gcp_compute_target_http_proxy - Query GCP Compute Engine Target HTTP Proxies using SQL

A Target HTTP Proxy is a component of GCP Compute Engine that terminates HTTP connections from clients and routes requests to a URL map. Target HTTP proxies are used by global and regional application load balancers and can be global or regional resources.

## Table Usage Guide

The `gcp_compute_target_http_proxy` table provides insights into global and regional Target HTTP Proxies within Google Cloud Compute Engine. Use it to find load balancers that accept plain HTTP traffic and to follow the path from a forwarding rule to its URL map.

## Examples

### Basic info
Explore target HTTP proxies and the URL map each one routes to.

```sql+postgres
select
  name,
  id,
  url_map,
  location
from
  gcp_compute_target_http_proxy;
```

```sql+sqlite
select
  name,
  id,
  url_map,
  location
from
  gcp_compute_target_http_proxy;
```

### List regional target HTTP proxies
Find target HTTP proxies defined in a region.

```sql+postgres
select
  name,
  location,
  url_map
from
  gcp_compute_target_http_proxy
where
  location_type = 'REGIONAL';
```

```sql+sqlite
select
  name,
  location,
  url_map
from
  gcp_compute_target_http_proxy
where
  location_type = 'REGIONAL';
```

### List target HTTP proxies with proxy bind enabled
Identify proxies used by Traffic Director with proxy bind enabled.

```sql+postgres
select
  name,
  proxy_bind
from
  gcp_compute_target_http_proxy
where
  proxy_bind;
```

```sql+sqlite
select
  name,
  proxy_bind
from
  gcp_compute_target_http_proxy
where
  proxy_bind = 1;
```
