select name, id, kind, description
from gcp.gcp_compute_target_vpn_gateway
where name = 'dummy-{{ resourceName }}'