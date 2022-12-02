select
  title,
  akas
from
  gcp_compute_ha_vpn_gateway
where 
  name = '{{ resourceName }}';