select name, description
from gcp.gcp_compute_ssl_policy
where akas::text = '["{{ output.resource_aka.value }}"]';