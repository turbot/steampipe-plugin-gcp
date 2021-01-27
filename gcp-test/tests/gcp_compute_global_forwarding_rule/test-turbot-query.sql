select title, akas
from gcp.gcp_compute_global_forwarding_rule
where name = '{{resourceName}}'