select name, dest_range, description, kind, network, self_link, next_hop_ip, project, priority
from gcp.gcp_compute_route
where name = '{{ resourceName }}'