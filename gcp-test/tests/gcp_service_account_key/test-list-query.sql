select name, key_type
from gcp.gcp_service_account_key
where name = '{{ output.name.value }}'