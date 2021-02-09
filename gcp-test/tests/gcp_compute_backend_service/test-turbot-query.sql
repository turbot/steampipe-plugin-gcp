select title, akas
from gcp.gcp_compute_backend_service
where name = '{{ resourceName }}'