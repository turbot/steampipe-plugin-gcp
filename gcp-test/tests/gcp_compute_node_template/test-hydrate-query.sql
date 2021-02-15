select name, description, kind, node_type, self_link, server_binding_type
from gcp.gcp_compute_node_template
where name = '{{ resourceName }}'