select name, self_link, health_checks, instances
from gcp.gcp_compute_target_pool
where name = '{{ resourceName }}'