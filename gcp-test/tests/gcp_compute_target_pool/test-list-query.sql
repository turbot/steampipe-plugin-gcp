select name, description, self_link
from gcp.gcp_compute_target_pool
where title = '{{resourceName}}'