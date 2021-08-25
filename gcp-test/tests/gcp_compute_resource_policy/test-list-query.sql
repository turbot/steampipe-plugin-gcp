select name, description
from gcp.gcp_compute_resource_policy
where title = '{{ resourceName }}';