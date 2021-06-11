select title, akas, project, location
from gcp.gcp_kubernetes_node_pool
where name = '{{ resourceName }}dd' and location = '{{ output.location.value }}' and cluster_name = '{{ output.cluster_name.value }}'