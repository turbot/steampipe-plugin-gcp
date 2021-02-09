select name, id, description, kind, self_link
from gcp.gcp_compute_target_pool
where name = 'dummy-{{resourceName}}'