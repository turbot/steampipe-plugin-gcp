select name, id, kind, description
from gcp.gcp_compute_ssl_policy
where name = 'dummy-{{ resourceName }}';