select name, title, id, kind, description 
from gcp.gcp_dns_policy
where name = 'dummy-{{ resourceName }}';