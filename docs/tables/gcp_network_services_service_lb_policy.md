---
title: "Steampipe Table: gcp_network_services_service_lb_policy - Query GCP Network Services Service LB Policies using SQL"
description: "Allows users to query GCP Network Services Service LB Policies, including the load balancing algorithm, failover and isolation settings applied to backend services."
folder: "Network Services"
---

# Table: gcp_network_services_service_lb_policy - Query GCP Network Services Service LB Policies using SQL

A Service LB Policy is a Network Services resource that tunes how a global backend service spreads traffic across regions and zones. It sets the load balancing algorithm, automatic capacity drain, health based failover and isolation for the backend services that reference it.

## Table Usage Guide

The `gcp_network_services_service_lb_policy` table provides insights into Service LB Policies within Google Cloud Network Services. Use it to review which load balancing algorithm is applied and whether automatic capacity drain or failover is configured. Querying a project where the Network Services API is not enabled returns a 403 error, which can be ignored with the `ignore_error_codes` connection config argument.

## Examples

### Basic info
Explore service LB policies and their load balancing algorithm.

```sql+postgres
select
  title,
  load_balancing_algorithm,
  location,
  create_time
from
  gcp_network_services_service_lb_policy;
```

```sql+sqlite
select
  title,
  load_balancing_algorithm,
  location,
  create_time
from
  gcp_network_services_service_lb_policy;
```

### List policies with automatic capacity drain enabled
Find policies that move traffic away from unhealthy backends automatically.

```sql+postgres
select
  title,
  location
from
  gcp_network_services_service_lb_policy
where
  (auto_capacity_drain ->> 'enable')::boolean;
```

```sql+sqlite
select
  title,
  location
from
  gcp_network_services_service_lb_policy
where
  json_extract(auto_capacity_drain, '$.enable') = 1;
```

### Get failover health threshold for each policy
Review the health threshold that triggers failover.

```sql+postgres
select
  title,
  failover_config ->> 'failoverHealthThreshold' as failover_health_threshold
from
  gcp_network_services_service_lb_policy;
```

```sql+sqlite
select
  title,
  json_extract(failover_config, '$.failoverHealthThreshold') as failover_health_threshold
from
  gcp_network_services_service_lb_policy;
```
