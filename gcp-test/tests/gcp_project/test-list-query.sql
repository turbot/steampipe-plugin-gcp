select
  name,
  project_id,
  project_number
from
  gcp_project
where
  project_id = '{{ output.project_id.value }}';