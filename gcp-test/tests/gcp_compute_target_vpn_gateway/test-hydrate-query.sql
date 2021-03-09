select name, kind, description, self_link
from gcp.gcp_compute_target_vpn_gateway
where name = '{{ resourceName }}'