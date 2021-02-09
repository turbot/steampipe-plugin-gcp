select title, akas
from gcp.gcp_compute_node_group
where name = '{{ resourceName }}'