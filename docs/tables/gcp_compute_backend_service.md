---
title: "Steampipe Table: gcp_compute_backend_service - Query Google Cloud Compute Engine Backend Services using SQL"
description: "Allows users to query Google Cloud Compute Engine Backend Services, specifically providing insights into the configuration and status of these services."
folder: "Compute"
---

# Table: gcp_compute_backend_service - Query Google Cloud Compute Engine Backend Services using SQL

Google Cloud Compute Engine Backend Services are a part of Google Cloud's Load Balancing feature, providing a scalable, reliable, and efficient way to distribute traffic across various instances. They manage HTTP(S) Load Balancing by directing traffic to available instances based on the capacity and load of each instance. This service helps in optimizing resource utilization and minimizing latency.

## Table Usage Guide

The `gcp_compute_backend_service` table provides insights into Google Cloud Compute Engine Backend Services. As a DevOps engineer, you can explore details about these services through this table, including their configurations, statuses, and associated instances. Utilize it to manage and monitor the distribution of traffic across your instances, ensuring optimal resource utilization and performance.

## Examples

### Backend info of backend service
Determine the areas in which your Google Cloud Compute backend service is balancing its workload, along with identifying the associated network endpoint groups. This can help you manage workload distribution and optimize network performance.

```sql+postgres
select
  name,
  id,
  b ->> 'balancingMode' as balancing_mode,
  split_part(b ->> 'group', '/', 10) as network_endpoint_groups
from
  gcp_compute_backend_service,
  jsonb_array_elements(backends) as b;
```

```sql+sqlite
Error: The corresponding SQLite query is unavailable.
```

### List of backend services where health check is not configured
Discover the segments that lack health check configurations within the Google Cloud Platform's backend services. This can help in identifying potential vulnerabilities and ensuring optimal performance of the services.

```sql+postgres
select
  name,
  id,
  self_link,
  health_checks
from
  gcp_compute_backend_service
where
  health_checks is null;
```

```sql+sqlite
select
  name,
  id,
  self_link,
  health_checks
from
  gcp_compute_backend_service
where
  health_checks is null;
```

### List of backend services where connection draining timeout is less than 300 sec
Determine the areas in which backend services may experience connection issues due to a draining timeout of less than 300 seconds. This can be useful for troubleshooting and optimizing network performance.

```sql+postgres
select
  name,
  id,
  connection_draining_timeout_sec
from
  gcp_compute_backend_service
where
  connection_draining_timeout_sec < 300;
```

```sql+sqlite
select
  name,
  id,
  connection_draining_timeout_sec
from
  gcp_compute_backend_service
where
  connection_draining_timeout_sec < 300;
```

### List of backend services where logging is not enabled
Discover the segments that have logging disabled in your backend services. This can help in identifying areas where crucial event tracking might be missing, aiding in improving system monitoring and error detection.

```sql+postgres
select
  name,
  id,
  log_config_enable
from
  gcp_compute_backend_service
where
   not log_config_enable;
```

```sql+sqlite
select
  name,
  id,
  log_config_enable
from
  gcp_compute_backend_service
where
  log_config_enable = 0;
```