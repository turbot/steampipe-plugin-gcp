select name, direction, description, kind, action, self_link
from gcp.gcp_compute_firewall
where name = '{{ resourceName }}'