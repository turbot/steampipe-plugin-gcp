select tags, title, akas
from gcp.gcp_compute_instance
where name = '{{ resourceName }}';