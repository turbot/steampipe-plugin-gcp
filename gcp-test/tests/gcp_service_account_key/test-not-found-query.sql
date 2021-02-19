select name, service_account_name
from gcp.gcp_service_account_key
where name = 'dummy-{{ output.name.value }}' and service_account_name = '{{ output.service_account_name.value }}'