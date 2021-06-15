# Table: gcp_compute_address

Represents an address resource.

Each virtual machine instance has an ephemeral internal IP address and, optionally, an external IP address. To communicate between instances on the same network, you can use an instance's internal IP address. To communicate with the Internet and instances outside of the same network, you must specify the instance's external IP address.

## Examples

### Basic info

```sql
select
  address,
  id,
  address_type,
  creation_timestamp,
  ip_version,
  status,
  subnetwork,
  location
from
  gcp_compute_address;
```

### List of address which are not in use

```sql
select
  address,
  address_type,
  creation_timestamp,
  status
from
  gcp_compute_address where status != 'IN_USE' ;
```

### Address count by each network_tier

```sql
select
  network_tier,
  count(*)
from
  gcp_compute_address
group by
  network_tier
order by network_tier;
```

### Get details of users that are using an address

```sql
select
  name,
  address,
  id,
  jsonb_pretty(users)
from
  gcp_compute_address where name= 'test2';
```
