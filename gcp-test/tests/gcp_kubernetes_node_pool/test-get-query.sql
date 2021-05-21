select name, id, location, project, cluster_name
from gcp.gcp_kubernetes_node_pool
where name = '{{ resourceName }}' and location = '{{ output.location.value }}' and cluster_name = '{{ output.location.value }}'