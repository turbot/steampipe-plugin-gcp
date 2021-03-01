select name, description
from gcp.gcp_compute_instance_group
where title = '{{ resourceName }}';