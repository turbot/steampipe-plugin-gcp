select name, description, kind, auto_create_subnetworks, self_link, project, title, akas
from gcp.gcp_compute_network
where akas::text = '["{{ output.resource_aka.value }}"]'