select name, id, kind
from gcp.gcp_compute_route
where name = 'dummy-{{ resourceName }}'