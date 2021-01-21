select name, service_account_id, key_type, key_algorithm, key_origin
from gcp.gcp_service_account_key
where name = '{{ output.resource_id.value }}'