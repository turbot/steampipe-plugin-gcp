select name, bgp_asn, description, kind, bgp_advertise_mode, self_link
from gcp.gcp_compute_router
where name = '{{ resourceName }}'