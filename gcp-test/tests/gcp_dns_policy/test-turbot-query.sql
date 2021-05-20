select title, akas
from gcp_dns_policy
where name = '{{ resourceName }}';