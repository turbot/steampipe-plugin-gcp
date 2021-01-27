select name, description, kind
from gcp.gcp_compute_global_forwarding_rule
where title = '{{resourceName}}'