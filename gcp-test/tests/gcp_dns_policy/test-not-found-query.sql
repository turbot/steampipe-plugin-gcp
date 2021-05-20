select
  name,
  id
from
  gcp.gcp_dns_policy
where
  name = 'dummy-{{ resourceName }}';