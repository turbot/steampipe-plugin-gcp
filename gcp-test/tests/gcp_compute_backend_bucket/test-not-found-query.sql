select name, id, description, kind
from gcp.gcp_compute_backend_bucket
where name = 'dummy-{{ resourceName }}'