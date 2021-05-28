select title, akas
from gcp.gcp_compute_target_ssl_proxy
where name = '{{ resourceName }}';