select name, id, kind, description
from gcp.gcp_compute_target_https_proxy
where name = 'dummy-{{ resourceName }}';