select title, akas
from gcp.gcp_compute_backend_bucket
where name = '{{ resourceName }}'