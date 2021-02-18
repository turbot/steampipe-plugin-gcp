select name, description, allow_global_access, kind, self_link, all_ports, ip_protocol, target, project
from gcp.gcp_compute_global_forwarding_rule
where name = '{{resourceName}}'