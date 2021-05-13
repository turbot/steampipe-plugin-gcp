select name, title, akas
from gcp.gcp_compute_project_metadata
where name = '{{ output.project_id.value }}';