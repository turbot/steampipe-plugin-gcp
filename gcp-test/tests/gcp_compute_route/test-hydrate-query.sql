select name, dest_range, description, kind, network, self_link
from gcp.gcp_compute_route
where name = '{{ resourceName }}'