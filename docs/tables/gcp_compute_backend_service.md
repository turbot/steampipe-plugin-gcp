# Table: gcp_compute_backend_service

A backend service defines how Google Cloud load balancers distribute traffic. The backend service configuration contains a set of values, such as the protocol used to connect to back-ends, various distribution and session settings, health checks, and timeouts.

### Backend info of backend service

```sql
select
  name,
  id,
  b ->> 'balancingMode' as balancing_mode,
  split_part(b ->> 'group', '/', 10) as network_endpoint_groups
from
  gcp_compute_backend_service,
  jsonb_array_elements(backends) as b;
```


### List of backend services where health check is not configured

```sql
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

```sql
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

```sql
select
  name,
  id,
  log_config_enable
from
  gcp_compute_backend_service
where
   not log_config_enable;
```