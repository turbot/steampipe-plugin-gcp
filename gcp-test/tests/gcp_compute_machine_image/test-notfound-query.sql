select name, id
from gcp.gcp_compute_machine_image
where name = 'dummy{{ resourceName }}';