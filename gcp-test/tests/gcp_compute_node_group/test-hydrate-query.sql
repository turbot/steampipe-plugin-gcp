select name, description, kind, location, self_link
from gcp.gcp_compute_node_group
where name = '{{ resourceName }}'