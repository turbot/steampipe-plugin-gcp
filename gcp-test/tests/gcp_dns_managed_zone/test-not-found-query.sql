select name, id, kind, description
from gcp.gcp_dns_managed_zone
where name = 'dummy-{{ resourceName }}'