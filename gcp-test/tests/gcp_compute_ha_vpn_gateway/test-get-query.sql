select
  name,
  description,
  location,
  self_link,
  kind
from
  gcp_compute_ha_vpn_gateway
where 
  name = '{{ output.resource_name.value }}';