select name, description, kind, ike_version, peer_external_gateway_interface, peer_ip, self_link, region, project, target_vpn_gateway, vpn_gateway_interface, location
from gcp.gcp_compute_vpn_tunnel
where name = '{{ resourceName }}'