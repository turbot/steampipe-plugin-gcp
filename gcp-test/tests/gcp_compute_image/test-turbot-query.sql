select title, akas, tags
from gcp.gcp_compute_image
where name = '{{ resourceName }}' and project = '{{ output.project_id.value }}';