select name, description
from gcp.gcp_compute_backend_service
where title = '{{ resourceName }}'