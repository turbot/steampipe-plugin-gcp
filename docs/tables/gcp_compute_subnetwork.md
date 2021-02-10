# Table: gcp_compute_subnetwork

A subnetwork (also known as a subnet) is a logical partition of a Virtual Private Cloud network with one primary IP range and zero or more secondary IP ranges.

## Examples

### Subnetwork basic info

```sql
select
  name,
  gateway_address,
  ip_cidr_range,
  ipv6_cidr_range,
  private_ip_google_access,
  id,
  network_name
from
  gcp_compute_subnetwork;
```

### List of subnetworks where users have compute admin access assigned in a resource policy

```sql
select
  name,
  id,
  jsonb_array_elements_text(p -> 'members') as members,
  p ->> 'role' as role
from
  gcp_compute_subnetwork,
  jsonb_array_elements(iam_policy -> 'bindings') as p
where
  p ->> 'role' = 'roles/compute.admin';
```

### Secondary IP info of each subnetwork

```sql
select
  name,
  id,
  p ->> 'rangeName' as range_name,
  p ->> 'ipCidrRange' as ip_cidr_range
from
  gcp_compute_subnetwork,
  jsonb_array_elements(secondary_ip_ranges) as p;
```

### Subnet count per network

```sql
select
  network,
  count(*) as subnet_count
from
  gcp_compute_subnetwork
group by
  network;
```

### List subnetworks having VPC flow logging set to false

```sql
select
  name,
  id,
  enable_flow_logs
from
  gcp_compute_subnetwork
where
  not enable_flow_logs;
```

### IP Info subnets
```sql
select
  name,
  id,
  ip_cidr_range
  gateway_address,
  broadcast(ip_cidr_range),
  netmask(ip_cidr_range),
  network(ip_cidr_range),
  pow(2, 32 - masklen(ip_cidr_range)) -1 as hosts_per_subnet
from
  gcp_compute_subnetwork;
```