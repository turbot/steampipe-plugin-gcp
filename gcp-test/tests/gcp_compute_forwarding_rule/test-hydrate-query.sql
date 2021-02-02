select name, description, allow_global_access, kind, region, self_link, all_ports, ip_protocol
from gcp.gcp_compute_forwarding_rule
where name = '{{ resourceName }}'