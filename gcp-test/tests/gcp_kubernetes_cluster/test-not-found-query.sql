select name, id
from gcp.gcp_kubernetes_cluster
where name = 'dummy-{{ resourceName }}'