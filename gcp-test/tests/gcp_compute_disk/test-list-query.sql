select name, self_link, title, akas, location, project, tags, labels
from gcp.gcp_compute_disk
where akas::text = '["{{ output.resource_aka.value }}"]'