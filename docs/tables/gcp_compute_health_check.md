---
title: "Steampipe Table: gcp_compute_health_check - Query GCP Compute Engine Health Checks using SQL"
description: "Allows users to query GCP Compute Engine Health Checks, both global and regional, including their protocol, probe intervals, thresholds and protocol-specific settings."
folder: "Compute"
---

# Table: gcp_compute_health_check - Query GCP Compute Engine Health Checks using SQL

A health check in GCP Compute Engine determines whether backends (instance groups or network endpoint groups) of a load balancer respond to traffic. Health checks can be global or regional and probe backends with TCP, SSL, HTTP, HTTPS, HTTP/2 or gRPC.

## Table Usage Guide

The `gcp_compute_health_check` table provides insights into global and regional health checks within Google Cloud Compute Engine. As a network engineer, use it to review probe protocols, intervals, timeouts and thresholds, and to find health checks with aggressive or lenient settings. Legacy HTTP and HTTPS health checks (`httpHealthChecks` and `httpsHealthChecks`) are not included.

## Examples

### Basic info
Explore the health checks in a project, their type and where they are defined.

```sql+postgres
select
  name,
  type,
  location,
  check_interval_sec,
  timeout_sec
from
  gcp_compute_health_check;
```

```sql+sqlite
select
  name,
  type,
  location,
  check_interval_sec,
  timeout_sec
from
  gcp_compute_health_check;
```

### List HTTP health checks with their request path and port
Identify the path and port used by HTTP health checks.

```sql+postgres
select
  name,
  location,
  http_health_check ->> 'requestPath' as request_path,
  http_health_check ->> 'port' as port
from
  gcp_compute_health_check
where
  type = 'HTTP';
```

```sql+sqlite
select
  name,
  location,
  json_extract(http_health_check, '$.requestPath') as request_path,
  json_extract(http_health_check, '$.port') as port
from
  gcp_compute_health_check
where
  type = 'HTTP';
```

### List regional health checks
Find health checks defined in a region, as used by regional and internal load balancers.

```sql+postgres
select
  name,
  type,
  location
from
  gcp_compute_health_check
where
  location_type = 'REGIONAL';
```

```sql+sqlite
select
  name,
  type,
  location
from
  gcp_compute_health_check
where
  location_type = 'REGIONAL';
```

### List health checks with logging disabled
Find health checks that do not log probe state changes.

```sql+postgres
select
  name,
  location
from
  gcp_compute_health_check
where
  not coalesce((log_config ->> 'enable')::boolean, false);
```

```sql+sqlite
select
  name,
  location
from
  gcp_compute_health_check
where
  log_config is null
  or json_extract(log_config, '$.enable') is not 1;
```
