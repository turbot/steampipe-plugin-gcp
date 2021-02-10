select name, description, kind, self_link, health_checks
from gcp.gcp_compute_target_pool
where name = '{{ resourceName }}'