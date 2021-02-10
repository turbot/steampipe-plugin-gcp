select name, kind, description, self_link, project
from gcp.gcp_compute_target_vpn_gateway
where name = '{{ resourceName }}'