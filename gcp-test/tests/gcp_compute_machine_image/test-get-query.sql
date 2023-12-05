select name, title, kind, self_link
from gcp.gcp_compute_machine_image
where name = '{{ resourceName }}';