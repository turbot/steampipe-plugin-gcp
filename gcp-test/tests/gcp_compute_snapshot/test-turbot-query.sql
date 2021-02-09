select title, akas
from gcp.gcp_compute_snapshot
where name = '{{resourceName}}'