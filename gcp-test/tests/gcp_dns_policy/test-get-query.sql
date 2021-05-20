select
  name,
  description,
  kind,
  enable_inbound_forwarding,
  enable_logging,
  project,
  networks,
  target_name_servers
from gcp.gcp_dns_policy
where name = '{{ resourceName }}';