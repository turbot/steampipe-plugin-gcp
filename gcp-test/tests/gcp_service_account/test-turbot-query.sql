select title, akas
from gcp.gcp_service_account
where name = '{{ output.name.value }}'