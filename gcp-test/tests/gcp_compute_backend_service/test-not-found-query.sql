select name, id, enable_cdn, load_balancing_scheme
from gcp.gcp_compute_backend_service
where name = 'dummy-{{ resourceName }}'