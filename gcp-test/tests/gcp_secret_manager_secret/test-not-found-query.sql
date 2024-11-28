select name
from gcp.gcp_secret_manager_secret
where name = 'dummy-{{ output.resource_name.value }}'