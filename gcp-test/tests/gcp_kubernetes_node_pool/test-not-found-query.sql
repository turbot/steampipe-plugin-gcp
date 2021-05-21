select title, akas, project, location
from gcp.gcp_kubernetes_node_pool
where name = '{{ resourceName }}dd'