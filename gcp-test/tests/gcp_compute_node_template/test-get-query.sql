select name, description, kind, node_type, self_link, server_binding_type, region, location, project, node_affinity_labels
from gcp.gcp_compute_node_template
where name = '{{ resourceName }}'