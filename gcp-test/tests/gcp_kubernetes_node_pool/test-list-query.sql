select name, cluster_name, akas, location, project
from gcp.gcp_kubernetes_node_pool
where title = '{{ resourceName }}'