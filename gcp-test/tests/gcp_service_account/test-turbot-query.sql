select title, akas
from gcp.gcp_service_account
where name = '{{ output.resource_id.value }}'