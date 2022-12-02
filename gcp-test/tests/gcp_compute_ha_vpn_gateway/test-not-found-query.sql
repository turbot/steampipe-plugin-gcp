select
  name,
  id,
  description,
  location,
  self_link,
  kind
from
  gcp_compute_ha_vpn_gateway
where 
  name = '{{ resourceName }}-dummy';