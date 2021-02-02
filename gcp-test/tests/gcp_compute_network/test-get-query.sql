select name, description, gateway_ipv4, mtu, kind, auto_create_subnetworks, routing_mode, peerings, subnetworks, self_link, project, title, akas
from gcp.gcp_compute_network
where name = '{{ resourceName }}'