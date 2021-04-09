select title, akas
from gcp.gcp_compute_target_https_proxy
where name = '{{ resourceName }}';