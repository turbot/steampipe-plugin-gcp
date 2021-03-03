select title, akas
from gcp.gcp_compute_instance_group
where name = '{{ resourceName }}';