select name, cluster_name, akas, location, project
from gcp.gcp_kubernetes_node_pool
where name = '{{ resourceName }}' and location = '{{ output.location.value }}' and cluster_name = '{{ output.cluster_name.value }}'