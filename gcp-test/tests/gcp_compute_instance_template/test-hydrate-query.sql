select name, instance_description, description, kind, instance_machine_type, self_link
from gcp.gcp_compute_instance_template
where name = '{{ resourceName }}'