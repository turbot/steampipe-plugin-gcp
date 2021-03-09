select name, id, description
from gcp.gcp_compute_image
where name = '' and project = '{{ output.project_id.value }}';