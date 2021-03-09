select name, id, kind, description
from gcp.gcp_compute_router
where name = 'dummy-{{ resourceName }}'