select title, tags, akas
from gcp.gcp_dns_managed_zone
where name = '{{ resourceName }}';