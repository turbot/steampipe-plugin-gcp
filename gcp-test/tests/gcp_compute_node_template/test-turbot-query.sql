select title, akas
from gcp.gcp_compute_node_template
where name = '{{ resourceName }}'