select title, akas, preoject, location
from gcp.gcp_kubernetes_node_pool
where name = '{{ resourceName }}'