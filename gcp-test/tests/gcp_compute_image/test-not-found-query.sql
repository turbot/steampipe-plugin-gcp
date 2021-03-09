select name, id, kind, description
from gcp.gcp_compute_image
where name = 'dummy-{{ resourceName }}' and project = '{{ output.project_id.value }}';