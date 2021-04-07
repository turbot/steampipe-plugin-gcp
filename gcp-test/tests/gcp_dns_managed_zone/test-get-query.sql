select name, dns_name, description, kind, visibility, private_visibility_config_networks, labels, name_servers, location, project
from gcp.gcp_dns_managed_zone
where name = '{{ resourceName }}';