select name, title, project
from gcp.gcp_secret_manager_secret
where name = '{{ output.resource_name.value }}'