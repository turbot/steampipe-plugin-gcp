select
  name,
  title
from
  gcp_dns_policy
where akas :: text = '["{{ output.resource_aka.value }}"]';