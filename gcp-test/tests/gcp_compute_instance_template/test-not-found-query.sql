select name, id, description
from gcp.gcp_compute_instance_template
where name = 'dummy-{{ resourceName }}'