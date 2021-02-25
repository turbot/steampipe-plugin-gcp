select name, description
from gcp.gcp_compute_instance
where title = '{{ resourceName }}';