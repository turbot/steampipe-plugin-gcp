select
  name,
  ip_cidr_range,
  gateway_address,
  broadcast(ip_cidr_range),
  netmask(ip_cidr_range),
  network(ip_cidr_range),
  pow(2, 32 - masklen(ip_cidr_range)) -1 as hosts_per_subnet
from
  gcp_compute_subnetwork
where name = '{{ resourceName }}';