select title, akas, tags
from gcp.gcp_compute_forwarding_rule
where name = '{{ resourceName }}'