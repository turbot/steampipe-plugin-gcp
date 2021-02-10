select name, description
from gcp.gcp_compute_instance_template
where title = '{{ resourceName }}'