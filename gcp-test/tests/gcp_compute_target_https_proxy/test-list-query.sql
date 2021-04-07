select name, description
from gcp.gcp_compute_target_https_proxy
where akas::text = '["{{ output.resource_aka.value }}"]';