select name, description
from gcp.gcp_compute_node_template
where title = '{{ resourceName }}'