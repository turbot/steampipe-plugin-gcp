select name, id, kind, description
from gcp.google_compute_vpn_gateway
where name = 'dummy-{{ resourceName }}'