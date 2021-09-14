select name, location, services_ipv4_cidr, akas, project
from gcp.gcp_kubernetes_cluster
where name = '{{ resourceName }}' and location = '{{ output.location.value }}'