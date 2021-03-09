select name, id, kind, description
from gcp.gcp_compute_node_template
where name = 'dummy-{{ resourceName }}'