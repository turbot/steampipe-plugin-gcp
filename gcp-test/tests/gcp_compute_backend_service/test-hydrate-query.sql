select name, enable_cdn, description, kind, self_link
from gcp.gcp_compute_backend_service
where name = '{{ resourceName }}'