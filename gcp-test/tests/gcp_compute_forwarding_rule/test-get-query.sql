select name, description, allow_global_access, kind, self_link, all_ports, ip_protocol, region, backend_service, load_balancing_scheme, labels, location, project
from gcp.gcp_compute_forwarding_rule
where name = '{{ resourceName }}'