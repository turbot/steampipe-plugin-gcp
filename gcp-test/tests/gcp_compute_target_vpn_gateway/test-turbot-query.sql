select title, akas
from gcp.gcp_compute_target_vpn_gateway
where name = '{{ resourceName }}'