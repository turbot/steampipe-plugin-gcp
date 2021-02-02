select name, bgp_asn, description, kind, bgp_advertise_mode, bgp_advertised_groups, bgp_advertised_ip_ranges, self_link, location, project, network
from gcp.gcp_compute_router
where name = '{{ resourceName }}'