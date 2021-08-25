select title, akas
from gcp.gcp_compute_resource_policy
where name = '{{ resourceName }}';