select name, description, kind, default_service, self_link
from gcp.gcp_compute_url_map
where name = '{{ resourceName }}'