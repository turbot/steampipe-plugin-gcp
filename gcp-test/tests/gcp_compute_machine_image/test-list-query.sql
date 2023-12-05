select name, title, akas
from gcp.gcp_compute_machine_image
where akas::text = '["{{ output.resource_aka.value }}"]';