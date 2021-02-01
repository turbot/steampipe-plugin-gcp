# Table:  gcp_compute_router

Cloud Router is a fully distributed and managed Google Cloud service that programs custom dynamic routes and scales with network traffic.

## Examples

### Cloud router basic info

```sql
select
  name,
  asn,
  advertise_mode
from
  gcp_compute_router;
```


### List all routers with custom route advertisements

```sql
select
  name,
  asn,
  advertise_mode,
  advertised_ip_ranges
from
  gcp_compute_router
where advertise_mode = 'CUSTOM';
```
