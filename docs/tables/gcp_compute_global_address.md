# Table: gcp_compute_global_address

Global Address represents a Global Address resource. Global addresses are used for HTTP(S) load balancing.

### List of internal address type global addresses

```sql
select
  name,
  id,
  address,
  address_type
from
  gcp_compute_global_address
where
  address_type = 'INTERNAL';
```


### List of unused global addresses

```sql
select
  name,
  address,
  status
from
  gcp_compute_global_address
where
  status <> 'IN_USE';
```


### List of global addresses used for VPC peering

```sql
select
  name,
  address,
  purpose
from
  gcp_compute_global_address
where
  purpose = 'VPC_PEERING';
```