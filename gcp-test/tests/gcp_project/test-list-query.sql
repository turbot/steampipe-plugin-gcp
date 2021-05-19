select name, project_id
from gcp_project
where project_id = '{{ output.project_id.value }}';