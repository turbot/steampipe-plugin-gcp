select name, asn, description, kind, advertise_mode, advertised_groups, advertised_ip_ranges, self_link, location, project, network
from gcp.gcp_compute_router
where name = '{{ resourceName }}'