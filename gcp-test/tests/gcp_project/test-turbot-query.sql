select
  akas,
  name,
  title
from
  gcp_project
where
  name = '{{ output.current_project_name.value }}';