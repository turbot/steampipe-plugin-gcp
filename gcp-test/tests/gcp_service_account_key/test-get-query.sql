select name, service_account_name, key_type, key_algorithm, key_origin, project, location
from gcp.gcp_service_account_key
where name = '{{ output.name.value }}' and service_account_name = '{{ output.service_account_name.value }}'