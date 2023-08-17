select name, title, akas
from gcp.gcp_compute_machine_type
where akas::text = '["{{ output.resource_aka.value }}"]' and zone = 'us-east1-b';