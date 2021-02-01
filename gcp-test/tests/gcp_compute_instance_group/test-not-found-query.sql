select name, id, kind, description
from gcp.gcp_compute_instance_group
where name = 'dummy-{{ resourceName }}'