select name, title, akas
from gcp.gcp_compute_project
where name = '{{ output.project_id.value }}';