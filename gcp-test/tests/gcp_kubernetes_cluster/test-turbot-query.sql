select title, akas
from gcp.gcp_kubernetes_cluster
where name = '{{ resourceName }}'