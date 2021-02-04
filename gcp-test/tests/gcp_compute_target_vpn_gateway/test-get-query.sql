select name, id, kind, description, region, self_link, project
from gcp.gcp_compute_target_vpn_gateway
where name = '{{ resourceName }}'