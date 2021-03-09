select name, description, kind
from gcp.gcp_compute_forwarding_rule
where title = '{{ resourceName }}'