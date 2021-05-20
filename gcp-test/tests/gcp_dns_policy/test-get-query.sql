select name, title, kind, enable_inbound_forwarding, enable_logging, project
from gcp.gcp_dns_policy
where name = '{{ resourceName }}';