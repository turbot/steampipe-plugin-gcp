select name, size, description, kind, self_link, zone
from gcp.gcp_compute_instance_group
where name = '{{ resourceName }}';