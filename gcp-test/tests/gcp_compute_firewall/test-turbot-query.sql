select title, akas
from gcp.gcp_compute_firewall
where name = '{{ resourceName }}'