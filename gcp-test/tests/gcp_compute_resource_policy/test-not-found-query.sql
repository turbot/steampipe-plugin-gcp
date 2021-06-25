select name, description
from gcp.gcp_compute_resource_policy
where name = 'dummy-{{ resourceName }}';