select name, size, description, kind, self_link, zone, project, named_ports, zone_name, location
from gcp.gcp_compute_instance_group
where name = '{{ resourceName }}'