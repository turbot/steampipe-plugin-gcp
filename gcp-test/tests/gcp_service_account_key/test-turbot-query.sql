select title, akas
from gcp.gcp_service_account_key
where name = '{{ output.resource_id.value }}'