select name, id, kind, description
from gcp.gcp_compute_vpn_tunnel
where name = 'dummy-{{ resourceName }}'