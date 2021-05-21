# Table:  gcp_kubernetes_cluster

Cloud Router is a fully distributed and managed Google Cloud service that programs custom dynamic routes and scales with network traffic.

## Examples

### Cloud router basic info

```sql
select
  name,
  bgp_asn,
  bgp_advertise_mode
from
  gcp_compute_router;
```


### NAT gateway info attached to router

```sql
select
  name,
  nat ->> 'name' as nat_name,
  nat ->> 'enableEndpointIndependentMapping' as enable_endpoint_independent_mapping,
  nat ->> 'natIpAllocateOption' as nat_ip_allocate_option,
  nat ->> 'sourceSubnetworkIpRangesToNat' as source_subnetwork_ip_ranges_to_nat
from
  gcp_compute_router,
  jsonb_array_elements(nats) as nat;
```


### List all routers with custom route advertisements

```sql
select
  name,
  bgp_asn,
  bgp_advertise_mode,
  bgp_advertised_ip_ranges
from
  gcp_compute_router
where bgp_advertise_mode = 'CUSTOM';
```
