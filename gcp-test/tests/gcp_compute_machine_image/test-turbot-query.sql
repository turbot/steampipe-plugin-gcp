select title, akas
from gcp.gcp_compute_machine_image
where name = '{{ resourceName }}';