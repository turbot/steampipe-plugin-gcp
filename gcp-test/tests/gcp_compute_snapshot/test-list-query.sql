select name, description, source_disk, disk_size_gb, self_link, project, title, akas
from gcp.gcp_compute_snapshot
where akas::text = '["{{ output.resource_aka.value }}"]'