select name, akas
from gcp.gcp_kubernetes_cluster
where name = 'dummy-{{ resourceName }}'