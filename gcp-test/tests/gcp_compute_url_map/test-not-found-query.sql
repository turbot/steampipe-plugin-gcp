select name, id, kind, description
from gcp.gcp_compute_url_map
where name = 'dummy-{{ resourceName }}'