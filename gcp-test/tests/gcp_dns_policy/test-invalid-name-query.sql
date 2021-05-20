select name, title, id, kind, description, enable_inbound_forwarding, enable_logging, target_name_servers, project
from gcp.gcp_dns_policy
where name = '';