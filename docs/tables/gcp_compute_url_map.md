---
title: "Steampipe Table: gcp_compute_url_map - Query Google Cloud Compute Engine URL Maps using SQL"
description: "Allows users to query URL Maps in Google Cloud Compute Engine, providing details about the routing rules and host rules."
folder: "Compute"
---

# Table: gcp_compute_url_map - Query Google Cloud Compute Engine URL Maps using SQL

A URL Map in Google Cloud Compute Engine is a key component of the HTTP(S) load balancing service. It defines the rules that route HTTP(S) traffic from an HTTP(S) load balancer to backend services, backend buckets, or other URL map targets. URL Maps contain both host rules and path matchers to control the flow of traffic.

## Table Usage Guide

The `gcp_compute_url_map` table provides insights into URL Maps within Google Cloud Compute Engine. As a network engineer, explore URL Map-specific details through this table, including routing rules, host rules, and associated metadata. Utilize it to uncover information about URL Maps, such as the traffic flow, the backend services or buckets the traffic is directed to, and the verification of host rules.

## Examples

### Get the default backend service of each url-map
Explore which backend service is set as the default for each URL map in your Google Cloud Platform compute instance. This can help in understanding how your web traffic is being directed and managed.

```sql+postgres
select
  name,
  id,
  default_service_name
from
  gcp_compute_url_map;
```

```sql+sqlite
select
  name,
  id,
  default_service_name
from
  gcp_compute_url_map;
```

### Path matcher info of each url-map
Explore the relationship between URL maps and their associated path matchers in your Google Cloud Platform setup. This can help identify areas for optimization or troubleshooting within your web service routing configuration.

```sql+postgres
select
  name,
  id,
  p ->> 'name' as name,
  r ->> 'paths' as paths,
  split_part(r ->> 'service', '/', 10) as service
from
  gcp_compute_url_map,
  jsonb_array_elements(path_matchers) as p,
  jsonb_array_elements(p -> 'pathRules') as r;
```

```sql+sqlite
Error: SQLite does not support split_part function.
```

### Host rule info of each url-map
Explore which URL maps have specific host rules in your Google Cloud Platform. This can help in identifying potential misconfigurations or anomalies in the distribution of network traffic.

```sql+postgres
select
  name,
  id,
  p ->> 'hosts' as hosts,
  p ->> 'pathMatcher' as path_matcher
from
  gcp_compute_url_map,
  jsonb_array_elements(host_rules) as p;
```

```sql+sqlite
select
  u.name,
  u.id,
  json_extract(p.value, '$.hosts') as hosts,
  json_extract(p.value, '$.pathMatcher') as path_matcher
from
  gcp_compute_url_map as u,
  json_each(host_rules) as p;
```

### List of all global type url-maps
Explore the global URL maps within your Google Cloud Platform's compute service. This can aid in managing and routing network traffic effectively.

```sql+postgres
select
  name,
  id,
  location_type
from
  gcp_compute_url_map
where
  location_type = 'GLOBAL';
```

```sql+sqlite
select
  name,
  id,
  location_type
from
  gcp_compute_url_map
where
  location_type = 'GLOBAL';
```