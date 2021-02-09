select title, akas
from gcp.gcp_compute_target_pool
where name = '{{ resourceName }}'