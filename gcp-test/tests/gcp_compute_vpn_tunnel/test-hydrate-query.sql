select name, description, kind, ike_version, self_link
from gcp.gcp_compute_vpn_tunnel
where name = '{{ resourceName }}'