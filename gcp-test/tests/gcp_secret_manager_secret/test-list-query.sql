select name, title
from gcp.gcp_secret_manager_secret
where akas::text = '["{{ output.resource_aka.value }}"]'