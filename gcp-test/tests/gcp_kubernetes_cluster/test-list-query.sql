select name, location, services_ipv4_cidr, akas
from gcp.gcp_kubernetes_cluster
where title = '{{ resourceName }}'