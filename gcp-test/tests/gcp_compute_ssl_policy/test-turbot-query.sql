select title, akas
from gcp.gcp_compute_ssl_policy
where name = '{{ resourceName }}';