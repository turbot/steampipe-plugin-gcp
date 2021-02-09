select title, akas
from gcp.gcp_compute_vpn_tunnel
where name = '{{ resourceName }}'