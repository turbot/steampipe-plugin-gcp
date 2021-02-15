select tags, title, akas
from gcp.gcp_compute_instance_template
where name = '{{ resourceName }}'