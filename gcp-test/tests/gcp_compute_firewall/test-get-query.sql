select name, direction, description, kind, disabled, self_link, action, project, network, log_config_enable, allowed, source_tags
from gcp.gcp_compute_firewall
where name = '{{ resourceName }}'