select title, akas, tags
from gcp.google_compute_vpn_gateway
where name = '{{ resourceName }}'