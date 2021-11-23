select name, id, kind
from gcp_compute_route
where name = 'dummy-{{ resourceName }}'