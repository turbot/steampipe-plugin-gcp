select name, asn, description, kind, advertise_mode, self_link
from gcp.gcp_compute_router
where name = '{{ resourceName }}'