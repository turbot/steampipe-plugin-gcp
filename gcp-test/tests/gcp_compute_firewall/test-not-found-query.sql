select name, id, direction, kind
from gcp.gcp_compute_firewall
where name = 'dummy-{{ resourceName }}'