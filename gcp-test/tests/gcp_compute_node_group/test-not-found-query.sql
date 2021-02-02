select name, id, size, kind
from gcp.gcp_compute_node_group
where name = 'dummy-{{ resourceName }}'