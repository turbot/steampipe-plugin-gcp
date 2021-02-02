select name, id, description, ip_address
from gcp.gcp_compute_forwarding_rule
where name = 'dummy-{{ resourceName }}'