select name, service_account_id
from gcp.gcp_service_account_key
where name = 'dummy-{{resourceName}}'