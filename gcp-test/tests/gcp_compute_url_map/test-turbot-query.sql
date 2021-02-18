select title, akas
from gcp.gcp_compute_url_map
where name = '{{ resourceName }}'