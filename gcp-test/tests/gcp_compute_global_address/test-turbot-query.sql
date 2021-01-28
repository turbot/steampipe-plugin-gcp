select title, akas
from gcp.gcp_compute_global_address
where name = '{{resourceName}}'